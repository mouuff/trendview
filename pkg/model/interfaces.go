package model

import "context"

// ItemStore is the interface for storing items
type ItemStore interface {
	SaveItem(item *ItemComposite) error
	FindItem(guid string) (*ItemComposite, error)
	FindItems() (ItemCompositeMap, error)
	UpdateResults(item *ItemComposite) error
	GetSubjects() ([]string, error)
	Close()
}

type FeedReader interface {
	GetFeedItems() ([]FeedItem, error)
	GetSource() string
}

type Brain interface {
	GenerateRating(ctx context.Context, prompt string) (int, error)
}
