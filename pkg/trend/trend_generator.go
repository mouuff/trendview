package trend

import (
	"context"
	"log"

	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/feed"
	"github.com/mouuff/TrendView/pkg/itemstore"
)

// TrendGenerator is responsible for generating trends based on the provided context,
// brain, storage, and feeds. It also maintains an internal state of items.
type TrendGenerator struct {
	// Context: The context for managing request-scoped values, cancellation, and deadlines.
	Context context.Context

	// Brain: The brain component responsible for processing and analyzing data.
	Brain brain.Brain

	// Storage: The item store for storing and retrieving items.
	Storage itemstore.ItemStore

	// Feeds: A list of feed readers for reading data from various sources.
	Feeds []feed.FeedReader

	// ConfidenceBasePrompt: A base prompt used for generating confidence levels.
	ConfidenceBasePrompt string

	// ReGenerate: A flag indicating whether to regenerate trends.
	ReGenerate bool

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

	log.Printf("Saving %d items", len(tg.items))
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
		if item.ConfidenceResult == nil || tg.ReGenerate {
			confidence, err := tg.Brain.GenerateConfidence(ctx, tg.ConfidenceBasePrompt+item.Content)

			if err != nil {
				log.Printf("Error generating confidence: %v\n", err)
				continue
			}

			item.ConfidenceResult = confidence
		}
	}
}
