package trend

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/feed"
)

type EnrichedFeedItem struct {
	feed.FeedItem
	Confidence *brain.ConfidenceResult
}

type TrendGenerator struct {
	Feeds            []feed.FeedReader
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

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) ReadFeeds() error {
	for _, feed := range tg.Feeds {
		feedItems, err := feed.GetFeedItems()
		if err != nil {
			log.Printf("Error reading feed: %v\n", err)
			continue
		}

		for _, item := range feedItems {
			enrichedItem := EnrichedFeedItem{
				FeedItem: item,
			}

			if item.GUID != "" {
				if _, exists := tg.items[item.GUID]; !exists {
					tg.items[item.GUID] = enrichedItem
				}
			}
		}
	}

	return nil
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) GenerateConfidenceScores(ctx context.Context) error {
	if tg.ConfidencePrompt == "" {
		return nil
	}

	for _, item := range tg.items {
		if item.Confidence == nil && tg.ConfidencePrompt != "" {
			confidence, err := tg.Brain.GenerateConfidence(ctx, tg.ConfidencePrompt)

			if err != nil {
				log.Printf("Error generating confidence: %v\n", err)
				continue
			}

			item.Confidence = confidence

		}
	}
	return nil
}
