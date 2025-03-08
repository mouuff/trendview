package brain_test

import (
	"context"
	"testing"

	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/model"
)

func TestOllamaBrainGenerateRating(t *testing.T) {
	ctx := context.Background()
	gen, err := brain.NewOllamaBrain("llama3.2")

	if err != nil {
		t.Fatal(err)
	}

	baseprompt := "Based solely on the news provided below, give a rating on how it might affect Bitcoin's price on a scale from 0 to 100, where: 0 indicates a very negative confidence (likely price drop), 50 indicates a neutral confidence (no significant change), and 100 indicates a positive confidence (likely price increase). For the rating, consider market trends, regulations, economic factors, and any other relevant information. News: "
	propertyName := "BitcoinConfidenceRating"
	testGenerateRating(t, ctx, gen, propertyName, baseprompt+"Top 10 AI Tools in 2023 That Will Make Your Life Easier", 0, 50)
	testGenerateRating(t, ctx, gen, propertyName, baseprompt+"Bitcoin will be banned in the US in 2026 - CNN", 0, 35)
	testGenerateRating(t, ctx, gen, propertyName, baseprompt+"Bitcoin Falls 10%, Solana Crashes 20%: Why Is Crypto Market Down Today? - News18", 0, 50)
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
