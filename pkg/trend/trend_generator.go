package trend

import (
	"context"
	"log"

	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/feed"
	"github.com/mouuff/TrendView/pkg/itemstore"
)

type TrendGenerator struct {
	Context              context.Context
	Brain                brain.Brain
	Storage              itemstore.ItemStore
	Feeds                []feed.FeedReader
	ConfidenceBasePrompt string

	// Internal state
	items map[string]*itemstore.ItemComposite
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) Execute() error {

	if tg.Storage.Exists() {
		items, err := tg.Storage.Load()

		if err != nil {
			return err
		}

		tg.items = items

		log.Printf("Loaded %d existing items", len(tg.items))
	} else {
		log.Printf("No existing data found, starting from scratch")
		tg.items = make(map[string]*itemstore.ItemComposite)
	}

	tg.readFeeds()
	tg.generateConfidenceScores(tg.Context)

	return tg.Storage.Save(tg.items)
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) readFeeds() {
	for _, feed := range tg.Feeds {
		feedItems, err := feed.GetFeedItems()
		if err != nil {
			log.Printf("Error reading feed: %v\n", err)
			continue
		}

		for _, item := range feedItems {
			enrichedItem := itemstore.ItemComposite{
				FeedItem: item,
			}

			if item.GUID != "" {
				if _, exists := tg.items[item.GUID]; !exists {
					tg.items[item.GUID] = &enrichedItem
				}
			}
		}
	}
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) generateConfidenceScores(ctx context.Context) {
	if tg.ConfidenceBasePrompt == "" {
		return
	}

	for _, item := range tg.items {
		if item.Confidence == nil {
			confidence, err := tg.Brain.GenerateConfidence(ctx, tg.ConfidenceBasePrompt+item.Title)

			if err != nil {
				log.Printf("Error generating confidence: %v\n", err)
				continue
			}

			item.Confidence = confidence
		}
	}
}
