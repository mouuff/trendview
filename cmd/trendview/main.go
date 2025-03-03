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
	Confidence int `json:"confidence"`
}

func main() {

	p := Schema{
		Type: "object",
		Properties: map[string]Property{
			"confidence": {
				Type: "integer",
			},
		},
		Required: []string{"confidence"},
	}

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

	news := "President Trump will host the first White House Crypto Summit on Friday March 7. Attendees will include prominent founders, CEOs, and investors from the crypto industry. Look forward to seeing everyone there!"

	req := &api.GenerateRequest{
		Model:  "mistral",
		Prompt: "Based solely on the news provided below, evaluate your confidence in investing in Bitcoin on a scale from 0 to 100, where 0 represents no confidence at all (strong belief that investing would be unwise), 50 represents a neutral or mitigated stance (uncertain or balanced risk), and 100 represents very high confidence (strong belief that investing is a good opportunity). Consider factors such as market trends, regulatory developments, economic indicators, or other relevant details mentioned in the news. Provide a specific confidence score. News: " + news,
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
