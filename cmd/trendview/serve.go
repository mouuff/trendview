package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

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

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET requests
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Fetch all items
		items, err := storage.FindItems()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to fetch items: %v", err), http.StatusInternalServerError)
			return
		}

		// Set JSON content type
		w.Header().Set("Content-Type", "application/json")

		// Encode and send the response
		if err := json.NewEncoder(w).Encode(items); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
			return
		}
	})

	// Start the server
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	return nil
}
