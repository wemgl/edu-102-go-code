package translation

import (
	"fmt"
	"go.temporal.io/sdk/log"
	"time"

	"go.temporal.io/sdk/workflow"
)

func SayHelloGoodbye(ctx workflow.Context, input TranslationWorkflowInput) (TranslationWorkflowOutput, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 45,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	logger := workflow.GetLogger(ctx)
	logger = log.With(logger, "version", "v0.1.0")
	logger.Info("Starting SayHelloGoodbye workflow")

	helloInput := TranslationActivityInput{
		Term:         "Hello",
		LanguageCode: input.LanguageCode,
	}
	var helloResult TranslationActivityOutput
	err := workflow.ExecuteActivity(ctx, TranslateTerm, helloInput).Get(ctx, &helloResult)
	if err != nil {
		return TranslationWorkflowOutput{}, err
	}
	helloMessage := fmt.Sprintf("%s, %s", helloResult.Translation, input.Name)

	goodbyeInput := TranslationActivityInput{
		Term:         "Goodbye",
		LanguageCode: input.LanguageCode,
	}
	var goodbyeResult TranslationActivityOutput
	err = workflow.ExecuteActivity(ctx, TranslateTerm, goodbyeInput).Get(ctx, &goodbyeResult)
	if err != nil {
		return TranslationWorkflowOutput{}, err
	}
	goodbyeMessage := fmt.Sprintf("%s, %s", goodbyeResult.Translation, input.Name)

	output := TranslationWorkflowOutput{
		HelloMessage:   helloMessage,
		GoodbyeMessage: goodbyeMessage,
	}
	logger.Info("Completed SayHelloGoodbye workflow")

	return output, nil
}
