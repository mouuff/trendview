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
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) Execute() error {
	tg.readFeeds()
	tg.generateRatingScores(tg.Context)
	return nil
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
				foundItem, err := tg.Storage.FindItem(item.GUID)

				if err != nil {
					log.Printf("Error reading database: %v\n", err)
				}

				if err == nil && foundItem == nil {
					// If item is not found, save it
					err = tg.Storage.SaveItem(&enrichedItem)
					if err != nil {
						log.Printf("Error writing to database: %v\n", err)
					}
				}
			}
		}
	}
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) generateRatingScores(ctx context.Context) {
	// TODO query only items that need to be updated
	// Don't forget to take "ReGenerate" into account
	items, err := tg.Storage.FindItems()
	if err != nil {
		log.Printf("Error reading database: %v\n", err)
	}

	for _, item := range items {

		if item.Results == nil || tg.ReGenerate {
			item.Results = make(map[string]*model.RatingResult)
		}

		log.Printf("Generating rating for item: %s\n", item.Title)

		for _, ratingPrompt := range tg.RatingPrompts {
			err := tg.generateSingleRatingScore(ctx, ratingPrompt, item)
			if err != nil {
				log.Printf("Error generating rating: %v\n", err)
				continue
			}
		}

		err = tg.Storage.UpdateResults(item)
		if err != nil {
			log.Printf("Error saving rating: %v\n", err)
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
	if ratingPrompt.SubjectName == "" {
		return fmt.Errorf("variable SubjectName is required for rating prompt")
	}
	if ratingPrompt.InsightName == "" {
		return fmt.Errorf("variable InsightName is required for rating prompt")
	}

	_, resultExists := item.Results[ratingPrompt.Identifier]

	if !resultExists || tg.ReGenerate {
		ratingValue, err := tg.Brain.GenerateRating(ctx, ratingPrompt.BasePrompt+item.Content)

		if err != nil {
			return err
		}

		item.Results[ratingPrompt.Identifier] = &model.RatingResult{
			SubjectName: ratingPrompt.SubjectName,
			InsightName: ratingPrompt.InsightName,
			Value:       ratingValue,
		}

	}

	return nil
}
