package itemstore_test

import (
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/mouuff/TrendView/pkg/itemstore"
	"github.com/mouuff/TrendView/pkg/model"
)

func TestSQLiteItemStore(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "testdb.db")

	newStore := func(t *testing.T) model.ItemStore {
		store, err := itemstore.NewSQLiteItemStore(dbPath)
		if err != nil {
			t.Fatalf("Failed to create store: %v", err)
		}
		return store
	}

	sampleItem := &model.ItemComposite{
		FeedItem: model.FeedItem{
			Title:    "Test Article",
			Content:  "Test content",
			DateTime: time.Date(2025, 3, 4, 12, 0, 0, 0, time.UTC),
			Link:     "https://test.com/article",
			GUID:     "https://test.com/article",
			Source:   "test.com",
		},
		Results: model.RatingResultMap{
			"MicrosoftConfidence": {SubjectName: "Microsoft", InsightName: "Confidence", Value: 50},
			"MicrosoftRelevance":  {SubjectName: "Microsoft", InsightName: "Relevance", Value: 20},
		},
	}

	t.Run("GetSubjects", func(t *testing.T) {
		store := newStore(t)
		defer store.Close()

		// Test with empty results table
		subjects, err := store.GetSubjects()
		if err != nil {
			t.Errorf("GetSubjects failed on empty table: %v", err)
		}
		if len(subjects) != 0 {
			t.Errorf("Expected 0 subjects on empty table, got %d", len(subjects))
		}

		// Add sample data
		if err := store.SaveItem(sampleItem); err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		// Test with populated data
		subjects, err = store.GetSubjects()
		if err != nil {
			t.Errorf("GetSubjects failed: %v", err)
		}
		expectedSubjects := []string{"Microsoft"}
		if !reflect.DeepEqual(subjects, expectedSubjects) {
			t.Errorf("Subjects mismatch: got %v, want %v", subjects, expectedSubjects)
		}
	})

	t.Run("SaveItem", func(t *testing.T) {
		store := newStore(t)
		defer store.Close()

		if err := store.SaveItem(sampleItem); err != nil {
			t.Errorf("SaveItem failed: %v", err)
		}

		found, err := store.FindItem(sampleItem.GUID)
		if err != nil {
			t.Errorf("FindItem failed after save: %v", err)
		}
		if found == nil {
			t.Fatal("Item not found after save")
		}
		if !reflect.DeepEqual(found, sampleItem) {
			t.Errorf("Saved item differs: got %v, want %v", found, sampleItem)
		}
	})

	t.Run("FindItem", func(t *testing.T) {
		store := newStore(t)
		defer store.Close()

		if err := store.SaveItem(sampleItem); err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		found, err := store.FindItem(sampleItem.GUID)
		if err != nil {
			t.Errorf("FindItem failed: %v", err)
		}
		if found == nil {
			t.Fatal("Item not found")
		}
		if !reflect.DeepEqual(found, sampleItem) {
			t.Errorf("Found item mismatch: got %v, want %v", found, sampleItem)
		}

		nonExistent, err := store.FindItem("non-existent-guid")
		if err != nil {
			t.Errorf("FindItem failed for non-existent: %v", err)
		}
		if nonExistent != nil {
			t.Errorf("Expected nil for non-existent GUID, got %v", nonExistent)
		}
	})

	t.Run("FindItems", func(t *testing.T) {
		store := newStore(t)
		defer store.Close()

		secondItem := &model.ItemComposite{
			FeedItem: model.FeedItem{
				Title:    "Second Article",
				Content:  "Second content",
				DateTime: time.Date(2025, 3, 5, 12, 0, 0, 0, time.UTC),
				Link:     "https://test.com/article2",
				GUID:     "https://test.com/article2",
				Source:   "test.com",
			},
			Results: model.RatingResultMap{
				"GoogleConfidence": {SubjectName: "Google", InsightName: "Confidence", Value: 80},
			},
		}

		if err := store.SaveItem(sampleItem); err != nil {
			t.Fatalf("Setup failed for first item: %v", err)
		}
		if err := store.SaveItem(secondItem); err != nil {
			t.Fatalf("Setup failed for second item: %v", err)
		}

		items, err := store.FindItems()
		if err != nil {
			t.Errorf("FindItems failed: %v", err)
		}
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}

		if !reflect.DeepEqual(items[sampleItem.GUID], sampleItem) {
			t.Errorf("First item mismatch: got %v, want %v", items[sampleItem.GUID], sampleItem)
		}
		if !reflect.DeepEqual(items[secondItem.GUID], secondItem) {
			t.Errorf("Second item mismatch: got %v, want %v", items[secondItem.GUID], secondItem)
		}
	})

	t.Run("Close", func(t *testing.T) {
		store := newStore(t)

		store.Close()

		_, err := store.FindItems()
		if err == nil {
			t.Error("Expected error after closing database, got nil")
		}
	})

	t.Run("EmptyTable", func(t *testing.T) {
		store := newStore(t)
		defer store.Close()

		count, err := store.GetResultsCount()
		if err != nil {
			t.Fatalf("Failed to count results: %v", err)
		}
		if count != 3 {
			t.Errorf("Expected 3 ratings before removal, got %d", count)
		}

		// Test removing ratings from an empty table
		if err := store.RemoveAllRatings(); err != nil {
			t.Errorf("RemoveAllRatings failed on empty table: %v", err)
		}

		// Verify no ratings exist (should already be true)
		count, err = store.GetResultsCount()
		if err != nil {
			t.Fatalf("Failed to count results: %v", err)
		}
		if count != 0 {
			t.Errorf("Expected 0 ratings after removal, got %d", count)
		}
	})

	t.Run("AddToExistingArticle", func(t *testing.T) {
		store := newStore(t)
		defer store.Close()

		// Add an article without ratings
		if err := store.SaveItem(sampleItem); err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		count, err := store.GetResultsCount()
		if err != nil {
			t.Fatalf("Failed to count results: %v", err)
		}
		if count != 2 {
			t.Errorf("Expected 2 ratings before adding rating, got %d", count)
		}

		// Add a rating
		ratingResult := &model.RatingResult{SubjectName: "Subject", InsightName: "Insight", Value: 40}
		if err := store.AddRating(sampleItem.GUID, ratingResult); err != nil {
			t.Errorf("AddRating failed: %v", err)
		}

		// Verify count
		count, err = store.GetResultsCount()
		if err != nil {
			t.Fatalf("Failed to count results: %v", err)
		}
		if count != 3 {
			t.Errorf("Expected 3 rating, got %d", count)
		}
	})

	t.Run("SaveItem_Update", func(t *testing.T) {
		store := newStore(t)
		defer store.Close()

		if err := store.SaveItem(sampleItem); err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		updatedItem := &model.ItemComposite{
			FeedItem: model.FeedItem{
				Title:    "Updated Test Article",
				Content:  sampleItem.Content,
				DateTime: sampleItem.DateTime,
				Link:     sampleItem.Link,
				GUID:     sampleItem.GUID,
				Source:   sampleItem.Source,
			},
			Results: model.RatingResultMap{
				"MicrosoftConfidence": {SubjectName: "Microsoft", InsightName: "Confidence", Value: 75},
				"MicrosoftRelevance":  {SubjectName: "Microsoft", InsightName: "Relevance", Value: 30},
				"GoogleConfidence":    {SubjectName: "Google", InsightName: "Confidence", Value: 60},
			},
		}
		if err := store.SaveItem(updatedItem); err != nil {
			t.Errorf("SaveItem failed for update: %v", err)
		}

		found, err := store.FindItem(sampleItem.GUID)
		if err != nil {
			t.Errorf("FindItem failed after update: %v", err)
		}
		if found == nil {
			t.Fatal("Item not found after update")
		}
		if !reflect.DeepEqual(found, updatedItem) {
			t.Errorf("Updated item mismatch: got %v, want %v", found, updatedItem)
		}
	})
}
