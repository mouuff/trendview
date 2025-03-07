package model

import "context"

// ItemStore is the interface for storing items
type ItemStore interface {
	SaveItem(item *ItemComposite) error
	GetItem(guid string) (*ItemComposite, error)
	GetItems() (ItemCompositeMap, error)
	GetSubjects() ([]string, error)
	GetItemsWithoutRating(subject, insight string) ([]string, error)
	GetResultsCount() (int, error)
	Close()
}

type FeedReader interface {
	GetFeedItems() ([]FeedItem, error)
	GetSource() string
}

type Brain interface {
	GenerateRating(ctx context.Context, prompt string) (int, error)
}
