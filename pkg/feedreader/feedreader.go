package feedreader

import "time"

type FeedItem struct {
	Title    string
	Content  string
	DateTime time.Time
	Link     string
	GUID     string
}

type FeedReader interface {
	GetFeedItems() ([]FeedItem, error)
}
