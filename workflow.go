package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type userJourney struct {
	initialised           bool
	documentationComplete bool
	drawdownComplete      bool
}

// Would come from config
var initWait = 5 * time.Second
var drawdownWait = 10 * time.Second

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

	err = workflow.ExecuteActivity(ctx, TwentySuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	err = workflow.ExecuteActivity(ctx, FiftySuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	workflow.Sleep(ctx, 10*time.Second)

	return result, err
}
