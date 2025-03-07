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
			log.Printf("Error reading feed %s: %v\n", feed.GetSource(), err)
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

	if tg.ReGenerate {
		tg.Storage.RemoveAllRatings()
	}

	for _, ratingPrompt := range tg.RatingPrompts {
		guids, err := tg.Storage.GetItemsWithoutRating(ratingPrompt.SubjectName, ratingPrompt.InsightName)

		if err != nil {
			log.Printf("Error getting articles without rating: %v\n", err)
			continue
		}

		for _, guid := range guids {
			item, err := tg.Storage.FindItem(guid)

			if err != nil {
				log.Printf("Error getting item by guid: %v\n", err)
				continue
			}

			if item == nil {
				log.Printf("Item not found.")
				continue
			}

			ratingResult, err := tg.generateSingleRatingScore(ctx, ratingPrompt, &item.FeedItem)

			if err != nil {
				log.Printf("Error generating rating: %v\n", err)
				continue
			}

			err = tg.Storage.AddRating(guid, ratingResult)
			if err != nil {
				log.Printf("Error saving rating: %v\n", err)
			}

			log.Printf("Generated rating for item: %s\n", item.Title)
		}
	}
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) generateSingleRatingScore(
	ctx context.Context,
	ratingPrompt model.RatingPrompt,
	item *model.FeedItem) (*model.RatingResult, error) {
	if ratingPrompt.BasePrompt == "" {
		return nil, fmt.Errorf("variable BasePrompt is required for rating prompt")
	}
	if ratingPrompt.SubjectName == "" {
		return nil, fmt.Errorf("variable SubjectName is required for rating prompt")
	}
	if ratingPrompt.InsightName == "" {
		return nil, fmt.Errorf("variable InsightName is required for rating prompt")
	}

	// TODO make this configurable
	prompt := ratingPrompt.BasePrompt + item.Title

	if item.Content != "" {
		prompt = prompt + "\n\n" + item.Content
	}

	ratingValue, err := tg.Brain.GenerateRating(ctx, prompt)

	if err != nil {
		return nil, err
	}

	return &model.RatingResult{
		SubjectName: ratingPrompt.SubjectName,
		InsightName: ratingPrompt.InsightName,
		Value:       ratingValue,
	}, nil
}
