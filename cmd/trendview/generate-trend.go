package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mouuff/TrendView/pkg/feedreader"
)

type TrendGeneratorConfig struct {
	RssFeedReaders       []feedreader.RssFeedReader
	ConfidenceBasePrompt string
}

// Ms describes the generate-trend subcommand
// This command is used to generate trend
type GenerateTrend struct {
	flagSet *flag.FlagSet

	config string
	outDir string
}

func ReadFromJson(path string, dataOut interface{}) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(data), dataOut); err != nil {
		return err
	}

	return nil
}

// Name gets the name of the command
func (cmd *GenerateTrend) Name() string {
	return "generate-trend"
}

// Init initializes the command
func (cmd *GenerateTrend) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.config, "config", "", "configuration file (required)")
	cmd.flagSet.StringVar(&cmd.outDir, "outdir", "", "output directory (required)")
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
	if cmd.outDir == "" {
		log.Println("Please specify an output directory using -outdir")
		return errors.New("-outdir parameter required")
	}

	var config TrendGeneratorConfig
	err := ReadFromJson(cmd.config, &config)
	if err != nil {
		return err
	}

	return nil
}

func printConfigurationTemplate() {
	configTemplate := &TrendGeneratorConfig{
		RssFeedReaders: []feedreader.RssFeedReader{
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
