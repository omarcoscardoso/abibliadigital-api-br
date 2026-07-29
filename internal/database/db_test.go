package database_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"abibliadigital/internal/database"
)

func getTestStore(t *testing.T) *database.Store {
	t.Helper()
	dbPath := filepath.Join("..", "..", "biblia.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("biblia.db not found at %s", dbPath)
	}

	store, err := database.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	return store
}

func TestDatabaseQueryPerformance(t *testing.T) {
	store := getTestStore(t)
	defer store.Close()

	start := time.Now()
	chapter, err := store.GetChapter("nvi", "tg", 1)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if chapter == nil || len(chapter.Verses) == 0 {
		t.Fatalf("expected chapter verses, got empty")
	}

	t.Logf("GetChapter execution time: %s", duration)

	if duration > 50*time.Millisecond {
		t.Errorf("query took too long: %s (expected < 50ms in test environment)", duration)
	}
}

func TestDatabaseSearchPerformance(t *testing.T) {
	store := getTestStore(t)
	defer store.Close()

	start := time.Now()
	res, err := store.Search("nvi", "Deus")
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error searching: %v", err)
	}

	if res == nil || res.Occurrence == 0 {
		t.Fatalf("expected search results for 'Deus', got 0")
	}

	t.Logf("Search execution time: %s for %d occurrences", duration, res.Occurrence)

	if duration > 250*time.Millisecond {
		t.Errorf("search took too long: %s (expected < 250ms in CI environment)", duration)
	}
}
