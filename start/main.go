package main

import (
	"context"
	"fmt"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"

	//"go.temporal.io/sdk/interceptor"
	//"go.temporal.io/sdk/workflow"

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

	traceInterceptor, err := opentelemetry.NewTracingInterceptor(
		opentelemetry.TracerOptions{
			Tracer:         otel.GetTracerProvider().Tracer("Starter"),
			SpanContextKey: app.CustomContextKey,
		})

	c, err := client.Dial(client.Options{
		Interceptors: []interceptor.ClientInterceptor{traceInterceptor},
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		ID:        "greeting-workflow",
		TaskQueue: app.TaskQueueName,
	}
	name := "World"
	we, err := c.ExecuteWorkflow(ctx, options, app.Baker, name)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}
	var greeting string
	err = we.Get(ctx, &greeting)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}
	printResults(greeting, we.GetID(), we.GetRunID())
}

func printResults(greeting string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", greeting)
}
