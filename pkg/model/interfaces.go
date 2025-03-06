package model

import "context"

// ItemStore is the interface for storing and loading items
type ItemStore interface {
	Load() (map[string]*ItemComposite, error)
	Save(items map[string]*ItemComposite) error
	Exists() bool
}

type FeedReader interface {
	GetFeedItems() ([]FeedItem, error)
}

type Brain interface {
	GenerateRating(ctx context.Context, prompt string) (*RatingResult, error)
}
