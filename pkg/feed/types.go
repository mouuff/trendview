package feed

import "time"

type FeedItem struct {
	Title    string
	Content  string
	DateTime time.Time
	Link     string
	GUID     string
	Source   string
}

type FeedReader interface {
	GetFeedItems() ([]FeedItem, error)
}
