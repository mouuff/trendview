package main

import (
	"errors"
	"flag"
	"log"

	"github.com/mouuff/TrendView/pkg/itemstore"
)

// Ms describes the generate-trend subcommand
// This command is used to generate trend
type Migrate struct {
	flagSet *flag.FlagSet

	jsonfile string
	dbfile   string
}

// Name gets the name of the command
func (cmd *Migrate) Name() string {
	return "migrate"
}

// Init initializes the command
func (cmd *Migrate) Init(args []string) error {
	cmd.flagSet = flag.NewFlagSet(cmd.Name(), flag.ExitOnError)
	cmd.flagSet.StringVar(&cmd.jsonfile, "jsonfile", "", "json file (required)")
	cmd.flagSet.StringVar(&cmd.dbfile, "dbfile", "", "file used to load and store data (required)")
	return cmd.flagSet.Parse(args)
}

// Run runs the command
func (cmd *Migrate) Run() error {
	log.Println("Generating trend...")

	if cmd.jsonfile == "" {
		log.Println("Please pass the jsonfile file using -jsonfile")
		return errors.New("-jsonfile parameter required")
	}
	if cmd.dbfile == "" {
		log.Println("Please specify a data file using -dbfile (e.g. -dbfile data.json)")
		return errors.New("-dbfile parameter required")
	}

	storage, err := itemstore.NewSQLiteItemStore(cmd.dbfile)

	if err != nil {
		return err
	}

	defer storage.Close()

	jsonStorage := &itemstore.JsonItemStore{
		Filename: cmd.jsonfile,
	}

	items, err := jsonStorage.Load()

	if err != nil {
		return err
	}

	for _, item := range items {
		err = storage.SaveItem(item)
		if err != nil {
			return err
		}
	}

	return nil
}
