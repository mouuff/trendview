package internal

import (
	"encoding/json"
	"os"

	"github.com/mouuff/TrendView/pkg/feed"
)

func ReadFromJson(path string, dataOut interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), dataOut); err != nil {
		return err
	}

	return nil
}

func ConvertToFeedReaders(rssFeedReaders []feed.RssFeedReader) []feed.FeedReader {
	feedReaders := make([]feed.FeedReader, len(rssFeedReaders))
	for i, rssFeedReader := range rssFeedReaders {
		feedReaders[i] = &rssFeedReader
	}
	return feedReaders
}
