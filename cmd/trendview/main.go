package main

import (
	"context"
	"fmt"

	"github.com/mouuff/TrendView/pkg/feedreader"
	"github.com/mouuff/TrendView/pkg/generator"
)

func main() {
	provider := feedreader.NewGoogleRssProvider("BTC+Bitcoin+news+when:1h")
	feeditems, err := provider.GetFeedItems()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print first few reports as example
	for i, feeditem := range feeditems {
		if i >= 5 {
			break
		}
		fmt.Printf("Report %d:\n", i+1)
		fmt.Printf("Title: %s\n", feeditem.Title)
		fmt.Printf("Date: %s\n", feeditem.DateTime)
		fmt.Printf("Link: %s\n", feeditem.Link)
		fmt.Printf("Content: %s\n", feeditem.Content)
		fmt.Println("---")
	}

	ctx := context.Background()
	pro, err := generator.NewOllamaGenerator()

	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}

	// Based solely on the relevant news articles provided below, rate your confidence in investing in Bitcoin on a scale from 0 (no confidence, unwise) to 50 (neutral) to 100 (high confidence, good opportunity). Consider only market trends, regulatory changes, and significant economic factors that directly impact Bitcoin. If a news article is not relevant to these factors, score it as 50.
	prompt := "Based solely on the news below, rate your confidence in investing in Bitcoin from 0 (no confidence, unwise) to 50 (neutral) to 100 (high confidence, good opportunity), considering market trends, regulations, or economic factors. If the news isn't relevant, score it 50. News: "

	newsList := []string{
		"Bitcoin will be banned in the US in 2026",
		"MSFT STOCKS ARE UP",
		"HODL BITCOIN",
		"BITCOIN PI CYCLE TOP PREDICTION SHOWDOWN! 🚨 @BitcoinProMag predicts $236,598 on Sep 28, 2025! @BitboBTC x @PositiveCrypto calls for $258,263 on Dec 2, 2025! WHO WILL BE RIGHT? Will history repeat with the Pi Cycle’s perfect track record?",
		"President Trump will host the first White House Crypto Summit on Friday March 7. Attendees will include prominent founders, CEOs, and investors from the crypto industry. Look forward to seeing everyone there!",
	}

	for _, news := range newsList {
		result, err := pro.GenerateConfidence(ctx, prompt+news)
		if err != nil {
			fmt.Println("Error generating confidence:", err)
			return
		}
		fmt.Println(result.Confidence)
	}

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
