package brain

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ollama/ollama/api"
)

type ollamaRatingResult struct {
	Rating int `json:"rating"`
}

type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required"`
}

type Property struct {
	Type string `json:"type"`
}

type OllamaBrain struct {
	Model  string
	Client *api.Client
}

func NewOllamaBrain() (*OllamaBrain, error) {
	client, err := api.ClientFromEnvironment()

	if err != nil {
		return nil, err
	}

	return &OllamaBrain{
		Model:  "mistral",
		Client: client,
	}, nil
}

func (c *OllamaBrain) GenerateRating(ctx context.Context, prompt string) (int, error) {
	formatSchema := Schema{
		Type: "object",
		Properties: map[string]Property{
			"rating": {
				Type: "integer",
			},
		},
		Required: []string{"rating"},
	}

	var result ollamaRatingResult

	respFunc := func(resp api.GenerateResponse) error {
		err := json.Unmarshal([]byte(resp.Response), &result)
		if err != nil {
			return fmt.Errorf("failed to parse response: %v", err)
		}
		return nil
	}

	err := c.generate(ctx, prompt, formatSchema, respFunc)

	if err != nil {
		return 0, fmt.Errorf("failed to generate rating: %v", err)
	}

	return result.Rating, nil
}

func (c *OllamaBrain) generate(ctx context.Context, prompt string, formatSchema Schema, fn api.GenerateResponseFunc) error {
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
