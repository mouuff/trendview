package pkg

import (
	"context"
	"encoding/json"
	"fmt"

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

type OllamaProvider struct {
	Model  string
	Client *api.Client
}

func NewOllamaProvider(confidencePrompt string) (*OllamaProvider, error) {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		return nil, err
	}

	return &OllamaProvider{
		Model:  "mistral",
		Client: client,
	}, nil
}

func (c *OllamaProvider) PredictConfidence(prompt string) (*ConfidenceResult, error) {
	formatSchema := Schema{
		Type: "object",
		Properties: map[string]Property{
			"confidence": {
				Type: "integer",
			},
		},
		Required: []string{"confidence"},
	}

	format, err := json.Marshal(formatSchema)
	if err != nil {
		return nil, err
	}

	req := &api.GenerateRequest{
		Model:  "mistral",
		Prompt: prompt,
		Format: format,

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

	err = c.Client.Generate(ctx, req, respFunc)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
