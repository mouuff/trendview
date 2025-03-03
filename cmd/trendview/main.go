package main

import (
	"context"
	"fmt"

	"github.com/mouuff/TrendView/pkg/generator"
)

func main() {
	ctx := context.Background()
	pro, err := generator.NewOllamaGenerator()

	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}

	prompt := "Based solely on the news below, rate your confidence in investing in Bitcoin from 0 (no confidence, unwise) to 50 (neutral) to 100 (high confidence, good opportunity), considering market trends, regulations, or economic factors. If the news isn't relevant, score it 50. News: "

	newsList := []string{
		"Bitcoin will be banned in the US in 2026",
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

}
