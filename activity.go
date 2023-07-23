package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/contrib/opentelemetry"
	"go.temporal.io/sdk/interceptor"
)

func GreatSuccess(ctx context.Context, name string) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("Hello", "team"))
	span.SetAttributes(attribute.Bool("winning", true))
	greeting := fmt.Sprintf("Always success %s", name)
	return greeting, nil
}

func FiftySuccess(ctx context.Context, name string) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("Chance", "50"))

	traceInterceptor, err := opentelemetry.NewTracingInterceptor(
		opentelemetry.TracerOptions{
			Tracer:         otel.GetTracerProvider().Tracer("Starter"),
			SpanContextKey: CustomContextKey,
		})

	c, err := client.Dial(client.Options{
		Interceptors: []interceptor.ClientInterceptor{traceInterceptor},
	})

	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()
	options := client.StartWorkflowOptions{
		ID:        "baking-workflow",
		TaskQueue: TaskQueueName,
	}

	span.AddEvent("Starting new trace with link")
	link := trace.LinkFromContext(ctx)
	newTraceCtx, newSpan := otel.Tracer("Driver").Start(context.Background(), "Split", trace.WithLinks(link))

	_, err = c.ExecuteWorkflow(newTraceCtx, options, Driver, name)
	if err != nil {
		log.Fatalln("unable to start worfklow", err)
	}
	newSpan.End()
	span.AddEvent("Ended new trace with link")

	roll := rand.Intn(10)
	if roll < 8 {
		greeting := fmt.Sprintf("Fifty success %s", name)
		return greeting, nil
	} else {
		return "", errors.New("Fifty Fail")
	}

}

func TwentySuccess(ctx context.Context, name string) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("Chance", "20"))
	roll := rand.Intn(10)
	if roll < 8 {
		greeting := fmt.Sprintf("Twenty success %s", name)
		return greeting, nil
	} else {
		return "", errors.New("Twenty Fail")
	}

}

func FlipSuccess(ctx context.Context, name string) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("Chance", "50"))
	roll := rand.Intn(10)
	if roll < 8 {
		greeting := fmt.Sprintf("Unlikely success %s", name)
		return greeting, nil
	} else {
		return "", errors.New("Unlikely Fail")
	}

}

func UnlikelySuccess(ctx context.Context, name string) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("Chance", "30"))
	roll := rand.Intn(10)
	if roll < 8 {
		greeting := fmt.Sprintf("Unlikely success %s", name)
		return greeting, nil
	} else {
		return "", errors.New("Unlikely Fail")
	}

}
