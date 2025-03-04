package trend

import (
	"encoding/json"
	"io"
	"os"
)

type TrendJsonStorage struct {
	Filename string
}

// Load loads enriched feed items from a JSON file.
func (tg *TrendJsonStorage) Load() (map[string]*EnrichedFeedItem, error) {
	file, err := os.Open(tg.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var items map[string]*EnrichedFeedItem
	if err := json.Unmarshal(bytes, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// Save saves enriched feed items to a JSON file.
func (s *TrendJsonStorage) Save(items map[string]*EnrichedFeedItem) error {
	bytes, err := json.Marshal(items)
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.Filename, bytes, 0644); err != nil {
		return err
	}

	return nil
}

// FileExists checks if the data file exists.
func (s *TrendJsonStorage) Exists() bool {
	_, err := os.Stat(s.Filename)
	return !os.IsNotExist(err)
}
