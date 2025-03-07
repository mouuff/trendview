package model

import (
	"time"
)

// The processing item state
type ItemState map[string]*ItemComposite
type RatingResultMap map[string]*RatingResult

type FeedItem struct {
	Title    string
	Content  string
	DateTime time.Time
	Link     string
	GUID     string
	Source   string
}

type RatingResult struct {
	SubjectName string
	InsightName string
	Value       int
}

// ItemComposite is a feed item with all the generated results
type ItemComposite struct {
	FeedItem
	Results RatingResultMap
}

type RatingPrompt struct {
	Identifier  string
	SubjectName string
	InsightName string
	BasePrompt  string
}
