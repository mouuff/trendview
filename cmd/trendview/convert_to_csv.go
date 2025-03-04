package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/mouuff/TrendView/pkg/itemstore"
)

// Ms describes the generate-trend subcommand
// This command is used to generate trend
type ConvertToCsv struct {
	flagSet *flag.FlagSet

	datafile string
}

// Name gets the name of the command
func (cmd *ConvertToCsv) Name() string {
	return "convert-to-csv"
}

// Init initializes the command
func (cmd *ConvertToCsv) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.datafile, "datafile", "", "file used to load data (required)")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *ConvertToCsv) Run() error {
	if cmd.datafile == "" {
		log.Println("Please specify a data file using -datafile (e.g. -datafile data.json)")
		return errors.New("-datafile parameter required")
	}

	storage := &itemstore.JsonItemStore{
		Filename: cmd.datafile,
	}

	data, err := storage.Load()

	if err != nil {
		return err
	}

	for _, item := range data {
		formatedDate := item.DateTime.Format("2006-01-02 15:04:05")
		formatedTitle := strings.ReplaceAll(item.Title, ",", ";")
		fmt.Printf("%s,%d,%s\n", formatedDate, item.ConfidenceResult.Confidence, formatedTitle)
	}

	return nil
}
