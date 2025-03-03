package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/ollama/ollama/api"
)

type DataOut struct {
	Temperature int `json:"temperature"`
}

func main() {

	p := DataOut{
		Temperature: 0,
	}

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

	req := &api.GenerateRequest{
		Model:  "mistral",
		Prompt: "What is the temperature of a monkey?",
		Format: "json",

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
