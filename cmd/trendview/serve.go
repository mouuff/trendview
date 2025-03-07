package main

import (
	"errors"
	"flag"
	"log"

	"github.com/mouuff/TrendView/pkg/itemstore"
)

// Ms describes the generate-trend subcommand
// This command is used to generate trend
type Serve struct {
	flagSet *flag.FlagSet

	datafile string
}

// Name gets the name of the command
func (cmd *Serve) Name() string {
	return "serve"
}

// Init initializes the command
func (cmd *Serve) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.datafile, "datafile", "", "file used to load and store data (required)")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *Serve) Run() error {
	log.Println("Starting server...")

	if cmd.datafile == "" {
		log.Println("Please specify a data file using -datafile (e.g. -datafile data.json)")
		return errors.New("-datafile parameter required")
	}

	storage, err := itemstore.NewSQLiteItemStore(cmd.datafile)
	if err != nil {
		return err
	}
	defer storage.Close()

	return nil
}
