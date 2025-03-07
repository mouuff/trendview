package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/mouuff/TrendView/internal"
	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/feed"
	"github.com/mouuff/TrendView/pkg/generator"
	"github.com/mouuff/TrendView/pkg/itemstore"
	"github.com/mouuff/TrendView/pkg/model"
)

type TrendGeneratorConfig struct {
	RssFeedReaders []feed.RssFeedReader
	RatingPrompts  []model.RatingPrompt
}

// Ms describes the generate-trend subcommand
// This command is used to generate trend
type GenerateTrend struct {
	flagSet *flag.FlagSet

	config   string
	datafile string
	loop     bool
	regen    bool
}

// Name gets the name of the command
func (cmd *GenerateTrend) Name() string {
	return "generate"
}

// Init initializes the command
func (cmd *GenerateTrend) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.config, "config", "", "configuration file (required)")
	cmd.flagSet.StringVar(&cmd.datafile, "datafile", "", "file used to load and store data (required)")
	cmd.flagSet.BoolVar(&cmd.loop, "loop", false, "should we loop forever")
	cmd.flagSet.BoolVar(&cmd.regen, "regen", false, "should we regenerate trends")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *GenerateTrend) Run() error {
	log.Println("Generating trend...")

	if cmd.config == "" {
		log.Println("Please pass the configuration file using -config")
		log.Println("Here is an example configuration:")
		printConfigurationTemplate()
		return errors.New("-config parameter required")
	}
	if cmd.datafile == "" {
		log.Println("Please specify a data file using -datafile (e.g. -datafile data.json)")
		return errors.New("-datafile parameter required")
	}

	var config TrendGeneratorConfig
	err := internal.ReadFromJson(cmd.config, &config)
	if err != nil {
		return err
	}

	brain, err := brain.NewOllamaBrain()
	if err != nil {
		return err
	}

	storage, err := itemstore.NewSQLiteItemStore(cmd.datafile)
	if err != nil {
		return err
	}
	defer storage.Close()

	tg := &generator.TrendGenerator{
		Context:       context.Background(),
		Brain:         brain,
		Storage:       storage,
		Feeds:         internal.ConvertToFeedReaders(config.RssFeedReaders),
		RatingPrompts: config.RatingPrompts,
	}

	if cmd.regen {
		log.Println("Re-Generating...")
		tg.ReGenerate()
	}

	if cmd.loop {
		for {
			tg.Execute()
			time.Sleep(7 * time.Minute)
		}
	}

	return tg.Execute()
}

func printConfigurationTemplate() {
	configTemplate := &TrendGeneratorConfig{
		RssFeedReaders: []feed.RssFeedReader{
			{
				Url:             "https://news.google.com/rss/search?q=BTC+Bitcoin+news+when:1h&hl=en-US&gl=US&ceid=US:en",
				ShouldCleanHtml: true,
			},
		},
		RatingPrompts: []model.RatingPrompt{
			{
				SubjectName: "Bitcoin",
				InsightName: "Confidence",
				BasePrompt:  "Based solely on the news below, rate your confidence in investing in Bitcoin from 0 (no confidence, unwise) to 50 (neutral) to 100 (very confident, good opportunity), considering market trends, regulations, or economic factors. If the news isn't relevant, score it 50. News: ",
			},
			{
				SubjectName: "Bitcoin",
				InsightName: "Relevance",
				BasePrompt:  "Based exclusively on the news provided below, evaluate the potential connection to Bitcoin's price. Assign a rating on a scale from 0 to 100, where:  - 0 = completely unrelated - 50 = somewhat related - 100 = very much related If there is any uncertainty or insufficient information to determine relevance, default to a rating of 0. News: ",
			},
		},
	}

	jsonData, err := json.MarshalIndent(configTemplate, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonData))
}
