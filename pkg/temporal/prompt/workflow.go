package prompt

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type Message struct {
	User    string `json:"user"`
	Content string `json:"content"`
}

// ChatBot *Temporal* workflow will trigger the [Assistant.Acknowledge] activity with the
// input request [Message] and it will return the AI chat bot response from the configured
// activity under the hood
func ChatBot(ctx workflow.Context, req *Message) (string, error) {
	logger := workflow.GetLogger(ctx)
	logger.Debug("Execute workflow for chat bot acknowledge", "User", req.User)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	ack, err := ollama(ctx, req)
	if err != nil {
		logger.Error("Get acknowledge failed.", "Error", err)
		return "", nil
	}
	return ack, nil
}

func ollama(ctx workflow.Context, req *Message) (ack string, err error) {
	var a *OLLama
	err = workflow.ExecuteActivity(ctx, a.Acknowledge, req.Content).Get(ctx, &ack)
	return
}
