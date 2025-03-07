package model

import "context"

// ItemStore is the interface for storing and loading items
type ItemStore interface {
	Load() (ItemState, error)
	Save(items ItemState) error
	Exists() bool
}

/* TODO implement a proper database!!!

type ItemStoreV2 interface {
	SaveItem(item *ItemComposite) error
	FindItem(guid string) (*ItemComposite, error)
	FindItems() (map[string]*ItemComposite, error)
	Close()
}
*/

type FeedReader interface {
	GetFeedItems() ([]FeedItem, error)
}

type Brain interface {
	GenerateRating(ctx context.Context, prompt string) (int, error)
}
