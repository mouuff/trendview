package brain_test

import (
	"context"
	"testing"

	"github.com/mouuff/TrendView/pkg/brain"
)

func TestOllamaBrainGenerateRating(t *testing.T) {
	ctx := context.Background()
	gen, err := brain.NewOllamaBrain()

	if err != nil {
		t.Fatal(err)
	}

	baseprompt := "Based solely on the news below, rate your rating in investing in Bitcoin from 0 (no rating, unwise) to 50 (neutral) to 100 (high rating, good opportunity), considering market trends, regulations, or economic factors. If the news isn't relevant, score it 50. News: "

	testGenerateRating(t, ctx, gen, baseprompt+"Ways to Stabilize Your Earnings with HTXmining's Liquidity Staking Service Amidst Recent BTC Price Drop - ABC Money", 50, 50)
	testGenerateRating(t, ctx, gen, baseprompt+"Top 10 AI Tools in 2023 That Will Make Your Life Easier", 50, 50)
	testGenerateRating(t, ctx, gen, baseprompt+"Bitcoin will be banned in the US in 2026 - CNN", 0, 20)
	testGenerateRating(t, ctx, gen, baseprompt+"Bitcoin Falls 10%, Solana Crashes 20%: Why Is Crypto Market Down Today? - News18", 0, 45)
	testGenerateRating(t, ctx, gen, baseprompt+"El Salvador Increases Bitcoin Holdings, How Much BTC Does It Own Now? - Watcher Guru", 60, 90)
	testGenerateRating(t, ctx, gen, baseprompt+"El Salvador buys Bitcoin dip, adding 5 BTC amid price plunge to $83,000 - Crypto Briefing", 60, 90)
	testGenerateRating(t, ctx, gen, baseprompt+"President Trump will host the first White House Crypto Summit on Friday March 7. Attendees will include prominent founders, CEOs, and investors from the crypto industry. Look forward to seeing everyone there!", 60, 90)
}

func testGenerateRating(
	t *testing.T,
	ctx context.Context,
	gen brain.Brain,
	prompt string,
	expectedMinRating, expectedMaxRating int) {
	result, err := gen.GenerateRating(ctx, prompt)

	if err != nil {
		t.Fatal(err)
	}

	if result.Rating > expectedMaxRating {
		t.Errorf("Rating is %d should be less than %d for prompt: '%s'", result.Rating, expectedMaxRating, prompt)
	}

	if result.Rating < expectedMinRating {
		t.Errorf("Rating is %d should be more than %d for prompt: '%s'", result.Rating, expectedMinRating, prompt)
	}
}
