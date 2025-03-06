package model

import (
	"time"
)

// The processing item state
type ItemState map[string]*ItemComposite

type FeedItem struct {
	Title    string
	Content  string
	DateTime time.Time
	Link     string
	GUID     string
	Source   string
}

// ItemComposite is a feed item with all the generated results
type ItemComposite struct {
	FeedItem
	Results map[string]int
}

type RatingPrompt struct {
	Identifier string
	BasePrompt string
}
