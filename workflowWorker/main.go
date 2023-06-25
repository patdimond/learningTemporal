package main

import (
	"context"
	"time"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/worker"
	//	"go.temporal.io/sdk/workflow"
	"github.com/uber-go/tally/v4"
	"github.com/uber-go/tally/v4/prometheus"
	sdktally "go.temporal.io/sdk/contrib/tally"
	prom "github.com/prometheus/client_golang/prometheus"

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
	c, err := client.Dial(client.Options{
		MetricsHandler: sdktally.NewMetricsHandler(newPrometheusScope(prometheus.Configuration{
			ListenAddress: "0.0.0.0:8077",
			TimerType:     "histogram",
		})),
	})
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
	w.RegisterWorkflow(app.DodgeyWorkflow)
	w.RegisterActivity(app.GreatSuccess)
	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func newPrometheusScope(c prometheus.Configuration) tally.Scope {
	reporter, err := c.NewReporter(
		prometheus.ConfigurationOptions{
			Registry: prom.NewRegistry(),
			OnError: func(err error) {
				log.Println("error in prometheus reporter", err)
			},
		},
	)
	if err != nil {
		log.Fatalln("error creating prometheus reporter", err)
	}
	scopeOpts := tally.ScopeOptions{
		CachedReporter:  reporter,
		Separator:       prometheus.DefaultSeparator,
		SanitizeOptions: &sdktally.PrometheusSanitizeOptions,
		Prefix:          "temporal_samples",
	}
	scope, _ := tally.NewRootScope(scopeOpts, time.Second)
	scope = sdktally.NewPrometheusNamingScope(scope)

	log.Println("prometheus metrics scope created")
	return scope
}
