package main

import (
	"context"
	"fmt"

	"github.com/mouuff/TrendView/pkg/provider"
)

func main() {
	ctx := context.Background()
	pro, err := provider.NewOllamaProvider()

	if err != nil {
		fmt.Println("Error marshaling:", err)
		return
	}

	prompt := "Based solely on the news provided below, evaluate your confidence in investing in Bitcoin on a scale from 0 to 100, where 0 represents no confidence at all (strong belief that investing would be unwise), 50 represents a neutral or mitigated stance (uncertain or balanced risk), and 100 represents very high confidence (strong belief that investing is a good opportunity). Consider factors such as market trends, regulatory developments, economic indicators, or other relevant details mentioned in the news. Provide a specific confidence score. News: "

	newsA := "BITCOIN PI CYCLE TOP PREDICTION SHOWDOWN! 🚨 @BitcoinProMag predicts $236,598 on Sep 28, 2025! @BitboBTC x @PositiveCrypto calls for $258,263 on Dec 2, 2025! WHO WILL BE RIGHT? Will history repeat with the Pi Cycle’s perfect track record?"
	newsB := "President Trump will host the first White House Crypto Summit on Friday March 7. Attendees will include prominent founders, CEOs, and investors from the crypto industry. Look forward to seeing everyone there!"

	result, err := pro.PredictConfidence(ctx, prompt+newsA)

	if err != nil {
		fmt.Println("Error predicting confidence:", err)
		return
	}
	fmt.Println(result.Confidence)

	result, err = pro.PredictConfidence(ctx, prompt+newsB)

	if err != nil {
		fmt.Println("Error predicting confidence:", err)
		return
	}
	fmt.Println(result.Confidence)

}
