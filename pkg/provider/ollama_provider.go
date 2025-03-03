package provider

import (
	"context"
	"encoding/json"

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

	var result ConfidenceResult

	respFunc := func(resp api.GenerateResponse) error {
		err := json.Unmarshal([]byte(resp.Response), &result)
		if err != nil {
			return err
		}
		return nil
	}

	err := c.Generate(ctx, prompt, formatSchema, respFunc)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *OllamaProvider) Generate(ctx context.Context, prompt string, formatSchema Schema, fn api.GenerateResponseFunc) error {
	format, err := json.Marshal(formatSchema)
	if err != nil {
		return err
	}

	req := &api.GenerateRequest{
		Model:  c.Model,
		Prompt: prompt,
		Format: format,

		// set streaming to false
		Stream: new(bool),
	}

	err = c.Client.Generate(ctx, req, fn)
	if err != nil {
		return err
	}

	return nil
}
