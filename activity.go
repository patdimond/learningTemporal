package app

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func GreatSuccess(ctx context.Context, name string) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("Hello", "team"))
	greeting := fmt.Sprintf("Always success %s", name)
	return greeting, nil
}

func FiftySuccess(ctx context.Context, name string) (string, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("Chance", "50"))
	roll := rand.Intn(10)
	if roll < 5 {
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
	if roll < 2 {
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
	if roll < 5 {
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
	if roll < 3 {
		greeting := fmt.Sprintf("Unlikely success %s", name)
		return greeting, nil
	} else {
		return "", errors.New("Unlikely Fail")
	}

}
