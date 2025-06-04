package translation

import (
	"fmt"
	"go.temporal.io/sdk/log"
	"time"

	"go.temporal.io/sdk/workflow"
)

func SayHelloGoodbye(ctx workflow.Context, input TranslationWorkflowInput) (TranslationWorkflowOutput, error) {
	logger := workflow.GetLogger(ctx)
	logger = log.With(logger, "LanguageCode", input.LanguageCode, "Name", input.Name)
	logger.Info("A workflow to translate the name has been invoked")

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 45,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	logger.Debug("About to execute Hello activity")
	helloInput := TranslationActivityInput{
		Term:         "Hello",
		LanguageCode: input.LanguageCode,
	}
	var helloResult TranslationActivityOutput
	err := workflow.ExecuteActivity(ctx, TranslateTerm, helloInput).Get(ctx, &helloResult)
	if err != nil {
		logger.Error("Encountered an error during Hello translation", "Error", err)
		return TranslationWorkflowOutput{}, err
	}
	helloMessage := fmt.Sprintf("%s, %s", helloResult.Translation, input.Name)

	logger.Debug("Sleeping between translation calls")
	err = workflow.Sleep(ctx, time.Duration(10)*time.Second)
	if err != nil {
		return TranslationWorkflowOutput{}, err
	}

	goodbyeInput := TranslationActivityInput{
		Term:         "Goodbye",
		LanguageCode: input.LanguageCode,
	}
	var goodbyeResult TranslationActivityOutput
	logger.Debug("About to Goodbye execute activity")
	err = workflow.ExecuteActivity(ctx, TranslateTerm, goodbyeInput).Get(ctx, &goodbyeResult)
	if err != nil {
		logger.Error("Encountered an error during Goodbye translation", "Error", err)
		return TranslationWorkflowOutput{}, err
	}
	goodbyeMessage := fmt.Sprintf("%s, %s", goodbyeResult.Translation, input.Name)

	output := TranslationWorkflowOutput{
		HelloMessage:   helloMessage,
		GoodbyeMessage: goodbyeMessage,
	}

	logger = log.With(logger, "HelloMsg", output.HelloMessage, "GoodbyeMsg", output.GoodbyeMessage)
	logger.Debug("Workflow completed successfully")

	return output, nil
}
