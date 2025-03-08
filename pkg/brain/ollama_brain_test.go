package brain_test

import (
	"context"
	"testing"

	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/model"
)

func TestOllamaBrainGenerateRating(t *testing.T) {
	ctx := context.Background()
	gen, err := brain.NewOllamaBrain("mistral")

	if err != nil {
		t.Fatal(err)
	}

	baseprompt := "Based solely on the news below, rate your confidence in investing in Bitcoin from 0 (no confidence, unwise) to 50 (neutral) to 100 (very confident, good opportunity), considering market trends, regulations, or economic factors. If the news isn't relevant, score it 50. News: "
	propertyName := "BitcoinConfidenceRating"
	testGenerateRating(t, ctx, gen, propertyName, baseprompt+"Top 10 AI Tools in 2023 That Will Make Your Life Easier", 50, 50)
	testGenerateRating(t, ctx, gen, propertyName, baseprompt+"Bitcoin will be banned in the US in 2026 - CNN", 0, 35)
	testGenerateRating(t, ctx, gen, propertyName, baseprompt+"Bitcoin Falls 10%, Solana Crashes 20%: Why Is Crypto Market Down Today? - News18", 0, 45)
	testGenerateRating(t, ctx, gen, propertyName, baseprompt+"El Salvador Increases Bitcoin Holdings, How Much BTC Does It Own Now? - Watcher Guru", 60, 90)
	testGenerateRating(t, ctx, gen, propertyName, baseprompt+"El Salvador buys Bitcoin dip, adding 5 BTC amid price plunge to $83,000 - Crypto Briefing", 40, 90)
	testGenerateRating(t, ctx, gen, propertyName, baseprompt+"President Trump will host the first White House Crypto Summit on Friday March 7. Attendees will include prominent founders, CEOs, and investors from the crypto industry. Look forward to seeing everyone there!", 60, 90)
}

func testGenerateRating(
	t *testing.T,
	ctx context.Context,
	gen model.Brain,
	propertyName string,
	prompt string,
	expectedMinRating, expectedMaxRating int) {
	result, err := gen.GenerateRating(ctx, propertyName, prompt)

	if err != nil {
		t.Fatal(err)
	}

	if result > expectedMaxRating {
		t.Errorf("Rating is %d should be less than %d for prompt: '%s'", result, expectedMaxRating, prompt)
	}

	if result < expectedMinRating {
		t.Errorf("Rating is %d should be more than %d for prompt: '%s'", result, expectedMinRating, prompt)
	}
}
