package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func DodgeyWorkflow(ctx workflow.Context, name string) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 120,
	}

	ctx = workflow.WithActivityOptions(ctx, options)
	var result string

	err := workflow.ExecuteActivity(ctx, GreatSuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	err = workflow.ExecuteActivity(ctx, FiftySuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	err = workflow.ExecuteActivity(ctx, TwentySuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	err = workflow.ExecuteActivity(ctx, FlipSuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	err = workflow.ExecuteActivity(ctx, UnlikelySuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	return result, err
}
