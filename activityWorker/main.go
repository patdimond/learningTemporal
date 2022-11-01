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

	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	traceInterceptor, err := opentelemetry.NewTracingInterceptor(opentelemetry.TracerOptions{})
	// This worker hosts Activity functions
	w := worker.New(c, app.TaskQueueName, worker.Options{Interceptors: []interceptor.WorkerInterceptor{app.NewWorkerInterceptor(), traceInterceptor}})
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
