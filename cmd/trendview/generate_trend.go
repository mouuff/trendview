package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/mouuff/TrendView/internal"
	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/feed"
	"github.com/mouuff/TrendView/pkg/itemstore"
	"github.com/mouuff/TrendView/pkg/trend"
)

type TrendGeneratorConfig struct {
	RssFeedReaders       []feed.RssFeedReader
	ConfidenceBasePrompt string
}

// Ms describes the generate-trend subcommand
// This command is used to generate trend
type GenerateTrend struct {
	flagSet *flag.FlagSet

	config   string
	datafile string
}

// Name gets the name of the command
func (cmd *GenerateTrend) Name() string {
	return "generate-trend"
}

// Init initializes the command
func (cmd *GenerateTrend) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.config, "config", "", "configuration file (required)")
	cmd.flagSet.StringVar(&cmd.datafile, "datafile", "", "file used to load and store data (required)")
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

	storage := &itemstore.JsonItemStore{
		Filename: cmd.datafile,
	}

	tg := &trend.TrendGenerator{
		Context:              context.Background(),
		Brain:                brain,
		Storage:              storage,
		Feeds:                internal.ConvertToFeedReaders(config.RssFeedReaders),
		ConfidenceBasePrompt: config.ConfidenceBasePrompt,
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
		ConfidenceBasePrompt: "Based solely on the news below, rate your confidence in investing in Bitcoin from 0 (no confidence, unwise) to 50 (neutral) to 100 (high confidence, good opportunity), considering market trends, regulations, or economic factors. If the news isn't relevant, score it 50. News: ",
	}

	jsonData, err := json.MarshalIndent(configTemplate, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonData))
}
