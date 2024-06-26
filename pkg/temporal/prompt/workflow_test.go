package prompt

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func (s *UnitTestSuite) TestOLLama_ChatBotWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	var a *OLLama
	env.OnActivity(a.Acknowledge, "Hello").Return("Hello! How are you today?", nil)
	env.ExecuteWorkflow(ChatBot, &Message{
		User:    "Test",
		Content: "Hello",
	})

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) TestTransformers_ChatBotWorkflow() {
	env := s.NewTestWorkflowEnvironment()

	var a *Transformers
	env.OnActivity(a.Acknowledge, "Visa").Return("Visa. Welcome to Canada!", nil)
	env.ExecuteWorkflow(ChatBot, &Message{
		User:    "Test",
		Content: "Visa",
	})

	s.True(env.IsWorkflowCompleted())
	s.NoError(env.GetWorkflowError())

	env.AssertExpectations(s.T())
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
