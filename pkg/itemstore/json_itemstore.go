package itemstore

import (
	"encoding/json"
	"io"
	"os"

	"github.com/mouuff/TrendView/pkg/model"
)

type JsonItemStore struct {
	Filename string
}

// Load loads enriched feed items from a JSON file.
func (tg *JsonItemStore) Load() (map[string]*model.ItemComposite, error) {
	file, err := os.Open(tg.Filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var items map[string]*model.ItemComposite
	if err := json.Unmarshal(bytes, &items); err != nil {
		return nil, err
	}

	return items, nil
}

// Save saves enriched feed items to a JSON file.
func (s *JsonItemStore) Save(items map[string]*model.ItemComposite) error {
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
func (s *JsonItemStore) Exists() bool {
	_, err := os.Stat(s.Filename)
	return !os.IsNotExist(err)
}
