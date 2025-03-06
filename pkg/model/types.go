package model

import (
	"time"
)

type FeedItem struct {
	Title    string
	Content  string
	DateTime time.Time
	Link     string
	GUID     string
	Source   string
}

type RatingResult struct {
	Rating int `json:"rating"`
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
