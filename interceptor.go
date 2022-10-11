package app

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
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
	activityName := activity.GetInfo(ctx).ActivityType.Name
	activityStart := time.Now()
	logger := activity.GetLogger(ctx)
	logger.Info(fmt.Sprintf("Activity %s Started", activityName))
	result, err := a.Next.ExecuteActivity(ctx, in)
	activityDuration := time.Since(activityStart).Milliseconds()
	contextDuration := time.Since(activity.GetInfo(ctx).StartedTime).Milliseconds()
	if err != nil {
		logger.Error(fmt.Sprintf("Activity %s Ended with Error", activityName), "error_string", err.Error(), activityDurationKey, activityDuration, testKey, contextDuration)
	} else {
		logger.Info(fmt.Sprintf("Activity %s Ended", activityName), activityDurationKey, activityDuration, testKey, contextDuration)
	}
	return result, err
}

func (w *workflowInboundInterceptor) ExecuteWorkflow(ctx workflow.Context, in *interceptor.ExecuteWorkflowInput) (interface{}, error) {
	workflowName := workflow.GetInfo(ctx).WorkflowType.Name
	workflowStart := time.Now()
	logger := workflow.GetLogger(ctx)
	logger.Info(fmt.Sprintf("Workflow %s Started", workflowName))
	result, err := w.Next.ExecuteWorkflow(ctx, in)
	workflowDuration := time.Since(workflowStart).Milliseconds()
	contextDuration := time.Since(workflow.GetInfo(ctx).WorkflowStartTime).Milliseconds()
	if err != nil {
		logger.Error(fmt.Sprintf("%s Ended with Error", workflowName), "error_string", err.Error(), workflowDurationKey, workflowDuration, testKey, contextDuration)
	} else {
		logger.Info(fmt.Sprintf("%s Ended", workflowName), workflowDurationKey, workflowDuration, testKey, contextDuration)
	}
	return result, err
}
