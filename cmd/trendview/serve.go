package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"time"

	"github.com/mouuff/TrendView/internal"
	"github.com/mouuff/TrendView/pkg/brain"
	"github.com/mouuff/TrendView/pkg/generator"
	"github.com/mouuff/TrendView/pkg/itemstore"
)

// Ms describes the generate-trend subcommand
// This command is used to generate trend
type Serve struct {
	flagSet *flag.FlagSet

	config   string
	datafile string
	loop     bool
	regen    bool
}

// Name gets the name of the command
func (cmd *Serve) Name() string {
	return "serve"
}

// Init initializes the command
func (cmd *Serve) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.config, "config", "", "configuration file (required)")
	cmd.flagSet.StringVar(&cmd.datafile, "datafile", "", "file used to load and store data (required)")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *Serve) Run() error {
	log.Println("Starting server...")

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

	tg := &generator.TrendGenerator{
		Context:       context.Background(),
		Brain:         brain,
		Storage:       storage,
		Feeds:         internal.ConvertToFeedReaders(config.RssFeedReaders),
		RatingPrompts: config.RatingPrompts,
		ReGenerate:    cmd.regen,
	}

	if cmd.loop {
		for {
			tg.Execute()
			time.Sleep(5 * time.Minute)
		}
	}

	return tg.Execute()
}
