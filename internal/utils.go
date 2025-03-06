package internal

import (
	"encoding/json"
	"os"

	"github.com/mouuff/TrendView/pkg/feed"
	"github.com/mouuff/TrendView/pkg/model"
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

func ConvertToFeedReaders(rssFeedReaders []feed.RssFeedReader) []model.FeedReader {
	feedReaders := make([]model.FeedReader, len(rssFeedReaders))
	for i, rssFeedReader := range rssFeedReaders {
		feedReaders[i] = &rssFeedReader
	}
	return feedReaders
}
