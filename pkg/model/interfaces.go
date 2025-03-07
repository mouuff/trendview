package model

import "context"

// ItemStore is the interface for storing items
type ItemStore interface {
	SaveItem(item *ItemComposite) error
	FindItem(guid string) (*ItemComposite, error)
	FindItems() (map[string]*ItemComposite, error)
	UpdateResults(item *ItemComposite) error
	Close()
}

type FeedReader interface {
	GetFeedItems() ([]FeedItem, error)
}

type Brain interface {
	GenerateRating(ctx context.Context, prompt string) (int, error)
}
