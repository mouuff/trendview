package generator_test

import (
	"context"
	"testing"

	"github.com/mouuff/TrendView/pkg/generator"
)

func TestOllamaGeneratorGenerateConfidence(t *testing.T) {
	ctx := context.Background()
	gen, err := generator.NewOllamaGenerator()

	if err != nil {
		t.Fatal(err)
	}

	baseprompt := "Based solely on the news below, rate your confidence in investing in Bitcoin from 0 (no confidence, unwise) to 50 (neutral) to 100 (high confidence, good opportunity), considering market trends, regulations, or economic factors. If the news isn't relevant, score it 50. News: "

	testGenerateConfidence(t, ctx, gen, baseprompt+"Bitcoin will be banned in the US in 2026 - CNN", 0, 20)
	testGenerateConfidence(t, ctx, gen, baseprompt+"Bitcoin Falls 10%, Solana Crashes 20%: Why Is Crypto Market Down Today? - News18", 0, 40)
	testGenerateConfidence(t, ctx, gen, baseprompt+"El Salvador Increases Bitcoin Holdings, How Much BTC Does It Own Now? - Watcher Guru", 60, 90)
	testGenerateConfidence(t, ctx, gen, baseprompt+"El Salvador buys Bitcoin dip, adding 5 BTC amid price plunge to $83,000 - Crypto Briefing", 60, 90)
	testGenerateConfidence(t, ctx, gen, baseprompt+"President Trump will host the first White House Crypto Summit on Friday March 7. Attendees will include prominent founders, CEOs, and investors from the crypto industry. Look forward to seeing everyone there!", 60, 90)
}

func testGenerateConfidence(
	t *testing.T,
	ctx context.Context,
	gen generator.ConfidenceGenerator,
	prompt string,
	expectedMinConfidence, expectedMaxConfidence int) {
	result, err := gen.GenerateConfidence(ctx, prompt)

	if err != nil {
		t.Fatal(err)
	}

	if result.Confidence > expectedMaxConfidence {
		t.Errorf("Confidence should be less than %d for prompt: '%s'", expectedMaxConfidence, prompt)
	}

	if result.Confidence < expectedMinConfidence {
		t.Errorf("Confidence should be more than %d for prompt: '%s'", expectedMinConfidence, prompt)
	}
}
