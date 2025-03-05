package trend

import (
	"context"
	"fmt"
	"log"

	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/feed"
	"github.com/mouuff/TrendView/pkg/itemstore"
)

type RatingPrompt struct {
	Identifier string
	BasePrompt string
}

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

	// RatingPrompts: A base prompt used for generating rating levels.
	RatingPrompts []RatingPrompt

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
	tg.generateRatingScores(tg.Context)

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
func (tg *TrendGenerator) generateRatingScores(ctx context.Context) {
	for _, item := range tg.items {
		for _, ratingPrompt := range tg.RatingPrompts {
			err := tg.generateSingleRatingScore(ctx, ratingPrompt, item)
			if err != nil {
				log.Printf("Error generating rating: %v\n", err)
				continue
			}
		}
	}
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) generateSingleRatingScore(ctx context.Context, ratingPrompt RatingPrompt, item *itemstore.ItemComposite) error {
	if ratingPrompt.BasePrompt == "" {
		return fmt.Errorf("variable BasePrompt is required for rating prompt")
	}
	if ratingPrompt.Identifier == "" {
		return fmt.Errorf("variable Identifier is required for rating prompt")
	}

	if item.Results == nil || tg.ReGenerate {
		item.Results = make(map[string]*brain.RatingResult)
	}

	_, resultExists := item.Results[ratingPrompt.Identifier]

	if !resultExists || tg.ReGenerate {
		rating, err := tg.Brain.GenerateRating(ctx, ratingPrompt.BasePrompt+item.Content)

		if err != nil {
			return err
		}

		item.Results[ratingPrompt.Identifier] = rating
	}

	return nil
}
