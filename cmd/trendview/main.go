package main

import (
	"context"
	"fmt"

	"github.com/mouuff/TrendView/pkg/feedreader"
	"github.com/mouuff/TrendView/pkg/generator"
)

func main() {
	provider := feedreader.NewGoogleRssFeedReader("BTC+Bitcoin+news+when:1h")
	feeditems, err := provider.GetFeedItems()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	ctx := context.Background()
	pro, err := generator.NewOllamaGenerator()

	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}

	// Based solely on the relevant news articles provided below, rate your confidence in investing in Bitcoin on a scale from 0 (no confidence, unwise) to 50 (neutral) to 100 (high confidence, good opportunity). Consider only market trends, regulatory changes, and significant economic factors that directly impact Bitcoin. If a news article is not relevant to these factors, score it as 50.
	prompt := "Based solely on the news below, rate your confidence in investing in Bitcoin from 0 (no confidence, unwise) to 50 (neutral) to 100 (high confidence, good opportunity), considering market trends, regulations, or economic factors. If the news isn't relevant, score it 50. News: "

	// Print first few reports as example
	for i, feeditem := range feeditems {
		result, err := pro.GenerateConfidence(ctx, prompt+feeditem.Title)
		if err != nil {
			fmt.Println("Error generating confidence:", err)
			return
		}

		fmt.Printf("Report %d:\n", i+1)
		fmt.Printf("Title: %s\n", feeditem.Title)
		fmt.Printf("Date: %s\n", feeditem.DateTime)
		fmt.Printf("Confidence: %d\n", result.Confidence)
		fmt.Println("---")
	}
}
