package itemstore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/mouuff/TrendView/pkg/model"
	_ "modernc.org/sqlite"
)

// SQLiteItemStore implements ItemStoreV2 using SQLite
type SQLiteItemStore struct {
	db *sql.DB
}

// NewSQLiteItemStore creates a new SQLiteItemStore
func NewSQLiteItemStore(dbPath string) (*SQLiteItemStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	store := &SQLiteItemStore{db: db}
	if err := store.createTables(); err != nil {
		return nil, err
	}

	return store, nil
}

// createTables sets up the database schema
func (s *SQLiteItemStore) createTables() error {
	_, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS feed_items (
            guid TEXT PRIMARY KEY,
            title TEXT NOT NULL,
            content TEXT,
            datetime TEXT NOT NULL,
            link TEXT NOT NULL,
            source TEXT NOT NULL
        );
        CREATE TABLE IF NOT EXISTS results (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            article_guid TEXT NOT NULL,
            subject_name TEXT NOT NULL,
            insight_name TEXT NOT NULL,
            value INTEGER NOT NULL,
            FOREIGN KEY (article_guid) REFERENCES feed_items(guid) ON DELETE CASCADE
        );
    `)
	return err
}

// SaveItem saves or updates an ItemComposite in the database
func (s *SQLiteItemStore) SaveItem(item *model.ItemComposite) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        INSERT OR REPLACE INTO feed_items (guid, title, content, datetime, link, source)
        VALUES (?, ?, ?, ?, ?, ?)
    `, item.GUID, item.Title, item.Content, item.DateTime.Format(time.RFC3339), item.Link, item.Source)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM results WHERE article_guid = ?`, item.GUID)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
        INSERT INTO results (article_guid, subject_name, insight_name, value)
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, result := range item.Results {
		_, err = stmt.Exec(item.GUID, result.SubjectName, result.InsightName, result.Value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// UpdateResults updates only the results for an existing item by GUID
func (s *SQLiteItemStore) UpdateResults(item *model.ItemComposite) error {
	// Check if the item exists in feed_items first
	var exists bool
	err := s.db.QueryRow(`SELECT EXISTS(SELECT 1 FROM feed_items WHERE guid = ?)`, item.GUID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("no item found with GUID: %s", item.GUID)
	}

	// Use a transaction to update results
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete existing results for this GUID
	_, err = tx.Exec(`DELETE FROM results WHERE article_guid = ?`, item.GUID)
	if err != nil {
		return err
	}

	// Insert new results
	stmt, err := tx.Prepare(`
        INSERT INTO results (article_guid, subject_name, insight_name, value)
        VALUES (?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, result := range item.Results {
		_, err = stmt.Exec(item.GUID, result.SubjectName, result.InsightName, result.Value)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// FindItem retrieves an ItemComposite by GUID
func (s *SQLiteItemStore) FindItem(guid string) (*model.ItemComposite, error) {
	var item model.ItemComposite
	var datetimeStr string // Temporary string to hold the datetime value

	err := s.db.QueryRow(`
        SELECT title, content, datetime, link, guid, source
        FROM feed_items
        WHERE guid = ?
    `, guid).Scan(&item.Title, &item.Content, &datetimeStr, &item.Link, &item.GUID, &item.Source)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// Parse the datetime string into time.Time
	item.DateTime, err = time.Parse(time.RFC3339, datetimeStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse datetime: %v", err)
	}

	rows, err := s.db.Query(`
        SELECT subject_name, insight_name, value
        FROM results
        WHERE article_guid = ?
    `, guid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	item.Results = make(model.RatingResultMap)
	for rows.Next() {
		var result model.RatingResult
		err := rows.Scan(&result.SubjectName, &result.InsightName, &result.Value)
		if err != nil {
			return nil, err
		}
		key := fmt.Sprintf("%s-%s", result.SubjectName, result.InsightName)
		item.Results[key] = &result
	}

	return &item, rows.Err()
}

// FindItems retrieves all ItemComposites as a map keyed by GUID
func (s *SQLiteItemStore) FindItems() (map[string]*model.ItemComposite, error) {
	items := make(map[string]*model.ItemComposite)

	rows, err := s.db.Query(`
        SELECT guid, title, content, datetime, link, source
        FROM feed_items
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.ItemComposite
		var datetimeStr string // Temporary string to hold the datetime value
		err := rows.Scan(&item.GUID, &item.Title, &item.Content, &datetimeStr, &item.Link, &item.Source)
		if err != nil {
			return nil, err
		}
		// Parse the datetime string into time.Time
		item.DateTime, err = time.Parse(time.RFC3339, datetimeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse datetime: %v", err)
		}
		item.Results = make(model.RatingResultMap)
		items[item.GUID] = &item
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	rows, err = s.db.Query(`
        SELECT article_guid, subject_name, insight_name, value
        FROM results
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var guid, subjectName, insightName string
		var value int
		err := rows.Scan(&guid, &subjectName, &insightName, &value)
		if err != nil {
			return nil, err
		}
		if item, exists := items[guid]; exists {
			key := fmt.Sprintf("%s-%s", subjectName, insightName)
			item.Results[key] = &model.RatingResult{SubjectName: subjectName, InsightName: insightName, Value: value}
		}
	}

	return items, rows.Err()
}

// Close shuts down the database connection
func (s *SQLiteItemStore) Close() {
	s.db.Close()
}
