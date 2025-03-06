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

type RatingResult struct {
	Subject string
	Value   int
}

// ItemComposite is a feed item with all the generated results
type ItemComposite struct {
	FeedItem
	Results map[string]*RatingResult
}

type RatingPrompt struct {
	Identifier string
	BasePrompt string
}
