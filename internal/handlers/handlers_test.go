package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"

	"abibliadigital/internal/database"
	"abibliadigital/internal/handlers"
	"abibliadigital/internal/middleware"
	"abibliadigital/internal/models"
)

func setupTestServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	// Locate biblia.db relative to repository root
	dbPath := filepath.Join("..", "..", "biblia.db")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatalf("test database biblia.db does not exist at %s", dbPath)
	}

	store, err := database.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}

	h := handlers.NewAPIHandler(store)

	r := chi.NewRouter()
	r.Use(middleware.CORS)
	r.Use(middleware.CacheControl)

	r.Route("/api", func(r chi.Router) {
		r.Get("/check", h.Check)
		r.Get("/books", h.GetBooks)
		r.Get("/books/{abbrev}", h.GetBook)
		r.Get("/versions", h.GetVersions)
		r.Get("/verses/{version}/random", h.GetRandomVerse)
		r.Get("/verses/{version}/{abbrev}/random", h.GetRandomVerse)
		r.Get("/verses/{version}/{abbrev}/{chapter}", h.GetChapter)
		r.Get("/verses/{version}/{abbrev}/{chapter}/{number}", h.GetVerse)
		r.Post("/verses/search", h.Search)
	})

	ts := httptest.NewServer(r)
	cleanup := func() {
		ts.Close()
		store.Close()
	}

	return ts, cleanup
}

func TestCheckEndpoint(t *testing.T) {
	ts, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Get(ts.URL + "/api/check")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var res models.CheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Result != "success" {
		t.Errorf("expected result 'success', got '%s'", res.Result)
	}
}

func TestGetBooks(t *testing.T) {
	ts, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Get(ts.URL + "/api/books")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	if cc := resp.Header.Get("Cache-Control"); cc != "public, max-age=31536000" {
		t.Errorf("expected Cache-Control header, got '%s'", cc)
	}

	var books []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(books) == 0 {
		t.Errorf("expected non-empty list of books, got 0")
	}
}

func TestGetBook(t *testing.T) {
	ts, cleanup := setupTestServer(t)
	defer cleanup()

	// Existing book
	resp, err := http.Get(ts.URL + "/api/books/tg")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var book models.Book
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		t.Fatalf("failed to decode book: %v", err)
	}

	if book.Name != "Tiago" {
		t.Errorf("expected book name 'Tiago', got '%s'", book.Name)
	}

	// Non-existent book
	resp404, err := http.Get(ts.URL + "/api/books/invalidbook")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp404.Body.Close()

	if resp404.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", resp404.StatusCode)
	}
}

func TestGetChapter(t *testing.T) {
	ts, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Get(ts.URL + "/api/verses/nvi/tg/1")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var ch models.ChapterResponse
	if err := json.NewDecoder(resp.Body).Decode(&ch); err != nil {
		t.Fatalf("failed to decode chapter response: %v", err)
	}

	if ch.Book.Name != "Tiago" {
		t.Errorf("expected book Tiago, got %s", ch.Book.Name)
	}

	if ch.Chapter.Number != 1 {
		t.Errorf("expected chapter 1, got %d", ch.Chapter.Number)
	}

	if len(ch.Verses) != 27 {
		t.Errorf("expected 27 verses in Tiago chapter 1, got %d", len(ch.Verses))
	}
}

func TestGetVerse(t *testing.T) {
	ts, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Get(ts.URL + "/api/verses/nvi/tg/1/1")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var v models.SingleVerseResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("failed to decode verse response: %v", err)
	}

	if v.Number != 1 || v.Chapter != 1 {
		t.Errorf("expected chapter 1 number 1, got chapter %d number %d", v.Chapter, v.Number)
	}
}

func TestGetRandomVerse(t *testing.T) {
	ts, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Get(ts.URL + "/api/verses/nvi/random")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var v models.SingleVerseResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		t.Fatalf("failed to decode random verse: %v", err)
	}

	if v.Text == "" {
		t.Errorf("expected non-empty text in random verse")
	}
}

func TestSearch(t *testing.T) {
	ts, cleanup := setupTestServer(t)
	defer cleanup()

	body := []byte(`{"version":"nvi","search":"Deus"}`)
	resp, err := http.Post(ts.URL+"/api/verses/search", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("failed to make search request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var s models.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&s); err != nil {
		t.Fatalf("failed to decode search response: %v", err)
	}

	if s.Occurrence == 0 {
		t.Errorf("expected occurrences for search 'Deus', got 0")
	}
}

func TestGetVersions(t *testing.T) {
	ts, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Get(ts.URL + "/api/versions")
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var versions []models.VersionResponse
	if err := json.NewDecoder(resp.Body).Decode(&versions); err != nil {
		t.Fatalf("failed to decode versions: %v", err)
	}

	if len(versions) == 0 {
		t.Errorf("expected at least 1 version, got 0")
	}
}
