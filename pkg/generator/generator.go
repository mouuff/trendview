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
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) Execute() error {
	for _, feed := range tg.Feeds {
		feedItems, err := feed.GetFeedItems()
		if err != nil {
			log.Printf("Error reading feed %s: %v\n", feed.GetSource(), err)
			continue
		}

		for _, item := range feedItems {
			err := tg.ProcessItem(&item)
			if err != nil {
				log.Printf("Error processing item: %v\n", err)
				continue
			}

			log.Printf("Saved item: %s\n", item.Title)
		}
	}
	return nil
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) ReGenerate() {
	items, err := tg.Storage.FindItems()

	if err != nil {
		log.Printf("Error getting items: %v\n", err)
		return
	}

	for _, item := range items {
		err := tg.ProcessItem(&item.FeedItem)
		if err != nil {
			log.Printf("Error processing item: %v\n", err)
			continue
		}

		log.Printf("Saved item: %s\n", item.Title)
	}
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) ProcessItem(item *model.FeedItem) error {

	if item.GUID == "" {
		return fmt.Errorf("empty item GUID")
	}

	ratingResultMap, err := tg.generateRatingResultMap(item)

	if err != nil {
		return fmt.Errorf("error generating rating %s: %v", item.Source, err)
	}

	// If item is not found, save it
	err = tg.Storage.SaveItem(&model.ItemComposite{
		FeedItem: *item,
		Results:  ratingResultMap,
	})
	if err != nil {
		return fmt.Errorf("error writing to database: %v", err)
	}

	return nil
}

// ReadFeeds reads feed items from the feeds.
func (tg *TrendGenerator) generateRatingResultMap(item *model.FeedItem) (model.RatingResultMap, error) {
	ratingResultMap := make(model.RatingResultMap)

	for _, ratingPrompt := range tg.RatingPrompts {

		ratingResult, err := tg.generateSingleRatingScore(tg.Context, ratingPrompt, item)

		if err != nil {
			return nil, err
		}

		ratingResultMap[ratingResult.GetKey()] = ratingResult
	}

	return ratingResultMap, nil
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
