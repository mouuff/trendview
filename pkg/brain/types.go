package brain

import "context"

type ConfidenceResult struct {
	Confidence int `json:"confidence"`
}

type ConfidenceGenerator interface {
	GenerateConfidence(ctx context.Context, news string) (*ConfidenceResult, error)
}
