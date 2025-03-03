package provider

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

func NewOllamaProvider() (*OllamaProvider, error) {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		return nil, err
	}

	return &OllamaProvider{
		Model:  "mistral",
		Client: client,
	}, nil
}

func (c *OllamaProvider) PredictConfidence(ctx context.Context, prompt string) (*ConfidenceResult, error) {
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
		Model:  c.Model,
		Prompt: prompt,
		Format: format,

		// set streaming to false
		Stream: new(bool),
	}

	var result ConfidenceResult

	respFunc := func(resp api.GenerateResponse) error {
		// TODO remove println
		fmt.Println(resp.Response)
		err := json.Unmarshal([]byte(resp.Response), &result)
		if err != nil {
			return err
		}
		return nil
	}

	err = c.Client.Generate(ctx, req, respFunc)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
