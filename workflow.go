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

	logger := workflow.GetLogger(ctx)

	journey := &userJourney{}

	workflow.GoNamed(ctx, "initialisedJourney", func(ctx workflow.Context) {
		workflow.Sleep(ctx, initWait)
		// Take whatever action is necessary. Maybe spawn an activity to do something
		if journey.initialised {
			logger.Info("Good value for init SLI")
		} else {
			logger.Info("Bad value for init SLI")
		}
	})

	ctx = workflow.WithActivityOptions(ctx, options)
	var result string

	err := workflow.ExecuteActivity(ctx, GreatSuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}
	journey.initialised = true

	err = workflow.ExecuteActivity(ctx, TwentySuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	journey.documentationComplete = true

	workflow.GoNamed(ctx, "drawDownJourney", func(ctx workflow.Context) {
		workflow.Sleep(ctx, drawdownWait)
		if journey.initialised {
			logger.Info("Good value for drawdown SLI")
		} else {
			logger.Info("Bad value for drawdown SLI")
		}
	})

	err = workflow.ExecuteActivity(ctx, FiftySuccess, name).Get(ctx, &result)
	if err != nil {
		return "", err
	}

	journey.drawdownComplete = true

	workflow.Sleep(ctx, 20*time.Second)

	return result, err
}
