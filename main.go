package main

import (
	"log"
	"net/http"

	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create temporal client", err)
	}
	defer c.Close()

	http.HandleFunc("/chat", PromptHandler(c))
	log.Println("Server started at http://localhost:3002")
	log.Fatal(http.ListenAndServe(":3002", nil))
}
