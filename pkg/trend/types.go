package trend

import (
	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/feed"
)

type EnrichedFeedItem struct {
	feed.FeedItem
	Confidence *brain.ConfidenceResult
}

// TrendJsonStorage is the interface for storing and loading trend data
type TrendStorage interface {
	Load() (map[string]*EnrichedFeedItem, error)
	Save(items map[string]*EnrichedFeedItem) error
	Exists() bool
}
