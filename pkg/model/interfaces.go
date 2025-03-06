package model

import "context"

// ItemStore is the interface for storing and loading items
type ItemStore interface {
	Load() (ItemState, error)
	Save(items ItemState) error
	Exists() bool
}

type FeedReader interface {
	GetFeedItems() ([]FeedItem, error)
}

type Brain interface {
	GenerateRating(ctx context.Context, prompt string) (*RatingResult, error)
}
