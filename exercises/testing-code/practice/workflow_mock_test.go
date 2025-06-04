package translation

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"

	"go.temporal.io/sdk/testsuite"
)

func TestSuccessfulTranslationWithMocks(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}

	env := s.NewTestWorkflowEnvironment()

	workflowInput := TranslationWorkflowInput{
		Name:         "Pierre",
		LanguageCode: "fr",
	}

	helloInput := TranslationActivityInput{
		Term:         "Hello",
		LanguageCode: "fr",
	}
	helloOutput := TranslationActivityOutput{Translation: "Bonjour"}

	goodbyeInput := TranslationActivityInput{
		Term:         "Goodbye",
		LanguageCode: "fr",
	}
	goodbyeOutput := TranslationActivityOutput{Translation: "Au revoir"}

	env.OnActivity(TranslateTerm, mock.Anything, helloInput).Return(helloOutput, nil)
	env.OnActivity(TranslateTerm, mock.Anything, goodbyeInput).Return(goodbyeOutput, nil)

	env.ExecuteWorkflow(SayHelloGoodbye, workflowInput)

	assert.True(t, env.IsWorkflowCompleted())

	var result TranslationWorkflowOutput
	err := env.GetWorkflowResult(&result)

	require.NoError(t, err)
	assert.Equal(t, "Bonjour, Pierre", result.HelloMessage)
	assert.Equal(t, "Au revoir, Pierre", result.GoodbyeMessage)
}
