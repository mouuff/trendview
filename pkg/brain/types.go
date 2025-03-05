package brain

import "context"

type RatingResult struct {
	Rating int `json:"rating"`
}

type Brain interface {
	GenerateRating(ctx context.Context, prompt string) (*RatingResult, error)
}
