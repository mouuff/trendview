package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type string `json:"type"`
}

type DataOut struct {
	Temperature string `json:"temperature"`
}

func main() {

	p := Schema{
		Type: "object",
		Properties: map[string]Property{
			"capital": {
				Type: "string",
			},
			"confidence": {
				Type: "integer",
			},
		},
		Required: []string{"capital", "confidence"},
	}

	//p := DataOut{
	//	Temperature: "integer",
	//}

	// Step 2: Marshal the struct to JSON bytes
	data, err := json.Marshal(p)
	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}

	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(data))

	req := &api.GenerateRequest{
		Model:  "mistral",
		Prompt: "What is the capital of Norway?",
		Format: data,

		// set streaming to false
		Stream: new(bool),
	}

	ctx := context.Background()
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		fmt.Println(resp.Response)
		return nil
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		log.Fatal(err)
	}
}
