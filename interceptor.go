package app

import (
	"context"
	"fmt"
	"time"

	//"go.temporal.io/sdk/activity"
	"go.opentelemetry.io/otel/trace"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/workflow"
)

const (
	activityDurationKey = "activity_duration"
	workflowDurationKey = "workflow_duration"
	testKey             = "from_context_duration"
)

type workerInterceptor struct {
	interceptor.WorkerInterceptorBase
}

type workflowInboundInterceptor struct {
	interceptor.WorkflowInboundInterceptorBase
}

type activityInboundInterceptor struct {
	interceptor.ActivityInboundInterceptorBase
}

func NewWorkerInterceptor() interceptor.WorkerInterceptor {
	return &workerInterceptor{}
}

func (w *workerInterceptor) InterceptActivity(ctx context.Context, next interceptor.ActivityInboundInterceptor) interceptor.ActivityInboundInterceptor {
	i := &activityInboundInterceptor{}
	i.Next = next
	return i
}

func (*workerInterceptor) InterceptWorkflow(ctx workflow.Context, next interceptor.WorkflowInboundInterceptor) interceptor.WorkflowInboundInterceptor {
	i := &workflowInboundInterceptor{}
	i.Next = next
	return i
}

func (a *activityInboundInterceptor) ExecuteActivity(ctx context.Context, in *interceptor.ExecuteActivityInput) (interface{}, error) {

	time := time.Now().String()
	span := trace.SpanFromContext(ctx)
	// Will get a noop span if there isn't one in context.
	fmt.Printf("%s --- ExecuteActivity Interceptor: TraceData: %s-%s-%s\n", time, span.SpanContext().TraceID(), span.SpanContext().SpanID(), span.SpanContext().TraceFlags())

	result, err := a.Next.ExecuteActivity(ctx, in)
	return result, err
}

func (w *workflowInboundInterceptor) ExecuteWorkflow(ctx workflow.Context, in *interceptor.ExecuteWorkflowInput) (interface{}, error) {
	time := time.Now().String()
	if span, ok := ctx.Value(CustomContextKey).(trace.Span); ok {
		fmt.Printf("%s --- ExecuteWorkflow Interceptor: TraceData: %s-%s-%s\n", time, span.SpanContext().TraceID(), span.SpanContext().SpanID(), span.SpanContext().TraceFlags())
	} else {
		fmt.Printf("%s --- ExecuteWorkflow Interceptor: Failed to get Span \n", time)
	}
	result, err := w.Next.ExecuteWorkflow(ctx, in)
	return result, err
}
