package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.temporal.io/sdk/client"

	//"go.temporal.io/sdk/interceptor"
	//"go.temporal.io/sdk/workflow"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflow/v1"
	"go.temporal.io/api/workflowservice/v1"

	"understandingTemporal/app"
)

func main() {
	ctx := context.Background()
	exp, err := app.OtlpExporter("localhost:4317")
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(app.NewResource()),
	)
	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("workflowStarter")

	ctx, span := tracer.Start(ctx, "Start")
	defer span.End()

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: app.TaskQueueName,
	}
	name := "World"
	we, err := c.ExecuteWorkflow(ctx, options, app.DodgeyWorkflow, name)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Poll until it completes
	id := we.GetID()
	done := false

	for {
		req := &workflowservice.ListWorkflowExecutionsRequest{
			Query: fmt.Sprintf("WorkflowId = '%s'", id),
		}
		resp, listErr := c.ListWorkflow(ctx, req)

		if listErr != nil {
			fmt.Printf("Failed to poll workflow status")
			fmt.Printf(listErr.Error())
		}
		for _, v := range resp.Executions {
			if v.Execution.WorkflowId == id {
				printInfo(v)
				if v.Status != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
					printClose(v)
					done = true
				}
			}
		}
		if done {
			break
		}
		time.Sleep(1 * time.Second)
	}

	var greeting string
	err = we.Get(ctx, &greeting)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}
	printResults(greeting, we.GetID(), we.GetRunID())
}

func printInfo(execution *workflow.WorkflowExecutionInfo) {
	now := time.Now()
	fmt.Printf("Id: %s\nStatus: %s\nStart: %s\nStartDiff: %s\nExec: %s\nExecDiff: %s\n",
		execution.Execution.WorkflowId,
		execution.Status,
		execution.StartTime,
		now.Sub(*execution.StartTime),
		execution.ExecutionTime,
		now.Sub(*execution.ExecutionTime),
	)
}

func printClose(execution *workflow.WorkflowExecutionInfo) {
	end := execution.CloseTime
	fmt.Printf("StartToEnd: %s\nExecToEnd: %s\n",
		end.Sub(*execution.StartTime),
		end.Sub(*execution.ExecutionTime),
	)
}

func printResults(greeting string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", greeting)
}
