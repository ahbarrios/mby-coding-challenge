package main

import (
	"log"

	"github.com/ahbarrios/mybeaconlabs/pkg/temporal/prompt"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Start client & worker to register challenge associated workflows and activities
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "prompts", worker.Options{})

	w.RegisterWorkflow(prompt.ChatBot)
	activity := prompt.FAQCanadaImmigration(prompt.AssistantOptions{
		URL: "https://ahbarrios-faq-canada-immigration.hf.space/generate",
	})
	w.RegisterActivity(activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
