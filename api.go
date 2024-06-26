package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ahbarrios/mybeaconlabs/pkg/temporal/prompt"

	"go.temporal.io/sdk/client"
)

type serverError struct {
	WorkflowID string `json:"operation_id"`
	Error      string `json:"error"`
}

// PromptHandler will interact with a *Temporal* workflow to retrieve AI generated responses
// using [prompt.Message] data as input from the returned [http.HandlerFunc].
func PromptHandler(c client.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var msg prompt.Message
		if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("unable to decode request: %v", err),
			})
			return
		}

		wo := client.StartWorkflowOptions{
			// Create and idempotent ID for this request
			ID:        "prompt_" + msg.User + "_" + strconv.FormatInt(time.Now().UnixMilli(), 10),
			TaskQueue: "prompts",
		}

		we, err := c.ExecuteWorkflow(context.Background(), wo, prompt.ChatBot, &msg)
		if err != nil {
			log.Fatalln("Unable to execute workflow", err)
		}
		log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

		// Synchronously wait for the workflow completion.
		var ack string
		if err := we.Get(r.Context(), &ack); err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(serverError{
				WorkflowID: we.GetID(),
				Error:      err.Error(),
			})
			return
		}

		json.NewEncoder(w).Encode(prompt.Message{
			User:    "Bot",
			Content: ack,
		})
	}
}
