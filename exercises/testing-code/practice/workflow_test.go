package translation

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"go.temporal.io/sdk/testsuite"
)

func TestSuccessfulCompleteFrenchTranslation(t *testing.T) {
	s := testsuite.WorkflowTestSuite{}

	env := s.NewTestWorkflowEnvironment()
	env.RegisterActivity(TranslateTerm)

	workflowInput := TranslationWorkflowInput{
		Name:         "Pierre",
		LanguageCode: "fr",
	}

	env.ExecuteWorkflow(SayHelloGoodbye, workflowInput)

	assert.True(t, env.IsWorkflowCompleted())

	var result TranslationWorkflowOutput
	err := env.GetWorkflowResult(&result)

	require.NoError(t, err)
	assert.Equal(t, "Bonjour, Pierre", result.HelloMessage)
	assert.Equal(t, "Au revoir, Pierre", result.GoodbyeMessage)
}
