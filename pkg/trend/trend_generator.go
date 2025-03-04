package trend

import (
	"encoding/json"
	"io"
	"os"

	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/feed"
)

type EnrichedFeedItem struct {
	feed.FeedItem
	Confidence *brain.ConfidenceResult
}

type TrendGenerator struct {
	Feeds            []feed.RssFeedReader
	Brain            brain.Brain
	ConfidencePrompt string

	// Internal variables:
	items map[string]EnrichedFeedItem
}

// Load loads enriched feed items from a JSON file.
func (tg *TrendGenerator) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var items map[string]EnrichedFeedItem
	if err := json.Unmarshal(bytes, &items); err != nil {
		return err
	}

	tg.items = items
	return nil
}

// Save saves enriched feed items to a JSON file.
func (tg *TrendGenerator) Save(filename string) error {
	bytes, err := json.Marshal(tg.items)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, bytes, 0644); err != nil {
		return err
	}

	return nil
}
