package main

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
	//	"go.temporal.io/sdk/workflow"

	"understandingTemporal/app"
)

func main() {
	exp, err := app.OtlpExporter("localhost:4317")
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(app.NewResource()),
	)
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()
	otel.SetTracerProvider(tp)
	// Create the client object just once per process
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	traceInterceptor, err := opentelemetry.NewTracingInterceptor(
		opentelemetry.TracerOptions{
			Tracer:         otel.GetTracerProvider().Tracer("WorkflowWorker"),
			SpanContextKey: app.CustomContextKey,
		})
	// This worker hosts Workflow
	w := worker.New(c, app.TaskQueueName, worker.Options{
		Interceptors: []interceptor.WorkerInterceptor{
			app.NewWorkerInterceptor(),
			traceInterceptor,
		},
	})
	w.RegisterWorkflow(app.Baker)
	w.RegisterWorkflow(app.Driver)
	w.RegisterActivity(app.GreatSuccess)
	w.RegisterActivity(app.FiftySuccess)
	w.RegisterActivity(app.TwentySuccess)
	w.RegisterActivity(app.FlipSuccess)
	w.RegisterActivity(app.UnlikelySuccess)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
