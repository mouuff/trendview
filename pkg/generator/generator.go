package generator

import (
	"context"
	"fmt"
	"log"

	"github.com/mouuff/TrendView/pkg/model"
)

// TrendGenerator is responsible for generating trends based on the provided context,
// brain, storage, and feeds. It also maintains an internal state of items.
type TrendGenerator struct {
	// Context: The context for managing request-scoped values, cancellation, and deadlines.
	Context context.Context

	// Brain: The brain component responsible for processing and analyzing data.
	Brain model.Brain

	// Storage: The item store for storing and retrieving items.
	Storage model.ItemStore

	// Feeds: A list of feed readers for reading data from various sources.
	Feeds []model.FeedReader

	// RatingPrompts: A base prompt used for generating rating levels.
	RatingPrompts []model.RatingPrompt

	// ReGenerate: A flag indicating whether to regenerate trends.
	ReGenerate bool

	// Internal state
	itemState model.ItemState
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) Execute() error {

	if tg.Storage.Exists() {
		items, err := tg.Storage.Load()

		if err != nil {
			return err
		}

		tg.itemState = items

		log.Printf("Loaded %d existing items", len(tg.itemState))
	} else {
		log.Printf("No existing data found, starting from scratch")
		tg.itemState = make(model.ItemState)
	}

	tg.readFeeds()
	tg.generateRatingScores(tg.Context)

	log.Printf("Saving %d items", len(tg.itemState))
	return tg.Storage.Save(tg.itemState)
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
			enrichedItem := model.ItemComposite{
				FeedItem: item,
			}

			if item.GUID != "" {
				if _, exists := tg.itemState[item.GUID]; !exists {
					tg.itemState[item.GUID] = &enrichedItem
				}
			}
		}
	}
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) generateRatingScores(ctx context.Context) {
	for _, item := range tg.itemState {

		if item.Results == nil || tg.ReGenerate {
			item.Results = make(map[string]*model.RatingResult)
		}

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
func (tg *TrendGenerator) generateSingleRatingScore(ctx context.Context, ratingPrompt model.RatingPrompt, item *model.ItemComposite) error {
	if ratingPrompt.BasePrompt == "" {
		return fmt.Errorf("variable BasePrompt is required for rating prompt")
	}
	if ratingPrompt.Identifier == "" {
		return fmt.Errorf("variable Identifier is required for rating prompt")
	}

	_, resultExists := item.Results[ratingPrompt.Identifier]

	if !resultExists || tg.ReGenerate {
		rating, err := tg.Brain.GenerateRating(ctx, ratingPrompt.BasePrompt+item.Content)

		if err != nil {
			return err
		}

		item.Results[ratingPrompt.Identifier] = &model.RatingResult{
			Subject: ratingPrompt.Subject,
			Value:   rating,
		}
	}

	return nil
}
