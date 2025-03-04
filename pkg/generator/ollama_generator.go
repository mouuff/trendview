package generator

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

type OllamaGenerator struct {
	Model  string
	Client *api.Client
}

func NewOllamaGenerator() (*OllamaGenerator, error) {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		return nil, err
	}

	return &OllamaGenerator{
		Model:  "mistral",
		Client: client,
	}, nil
}

func (c *OllamaGenerator) GenerateConfidence(ctx context.Context, prompt string) (*ConfidenceResult, error) {
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
			return fmt.Errorf("failed to parse response from ollama: %v", err)
		}
		return nil
	}

	err := c.generate(ctx, prompt, formatSchema, respFunc)

	if err != nil {
		return nil, fmt.Errorf("failed to generate confidence: %v", err)
	}

	return &result, nil
}

func (c *OllamaGenerator) generate(ctx context.Context, prompt string, formatSchema Schema, fn api.GenerateResponseFunc) error {
	format, err := json.Marshal(formatSchema)
	if err != nil {
		return fmt.Errorf("failed to marshal the format schema: %v", err)
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
		return fmt.Errorf("failed to generate response: %v", err)
	}

	return nil
}
