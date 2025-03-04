package itemstore

import (
	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/feed"
)

// ItemComposite is a feed item with all the generated results
type ItemComposite struct {
	feed.FeedItem
	ConfidenceResult *brain.ConfidenceResult
}

// ItemStore is the interface for storing and loading items
type ItemStore interface {
	Load() (map[string]*ItemComposite, error)
	Save(items map[string]*ItemComposite) error
	Exists() bool
}
