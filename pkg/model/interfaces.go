package model

import "context"

// ItemStore is the interface for storing items
type ItemStore interface {
	SaveItem(item *ItemComposite) error
	AddRating(articleGuid string, ratingResult *RatingResult) error
	FindItem(guid string) (*ItemComposite, error)
	FindItems() (ItemCompositeMap, error)
	GetSubjects() ([]string, error)
	GetItemsWithoutRating(subject, insight string) ([]string, error)
	RemoveAllRatings() error
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
