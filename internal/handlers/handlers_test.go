package handlers_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

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
	r.Use(chimiddleware.GetHead)
	r.Use(middleware.CORS)
	r.Use(middleware.CacheControl)

	r.Get("/health", h.Check)
	r.Head("/health", h.Check)
	r.Get("/healthz", h.Check)
	r.Head("/healthz", h.Check)

	r.Route("/api", func(r chi.Router) {
		r.Get("/check", h.Check)
		r.Head("/check", h.Check)
		r.Get("/health", h.Check)
		r.Head("/health", h.Check)
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

	endpoints := []string{"/api/check", "/api/health", "/health", "/healthz"}

	for _, endpoint := range endpoints {
		t.Run("GET "+endpoint, func(t *testing.T) {
			resp, err := http.Get(ts.URL + endpoint)
			if err != nil {
				t.Fatalf("failed to make request to %s: %v", endpoint, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status 200 for %s, got %d", endpoint, resp.StatusCode)
			}

			if cc := resp.Header.Get("Cache-Control"); cc != "no-cache, no-store, must-revalidate" {
				t.Errorf("expected no-cache Cache-Control for %s, got '%s'", endpoint, cc)
			}

			var res models.CheckResponse
			if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
				t.Fatalf("failed to decode response for %s: %v", endpoint, err)
			}

			if res.Result != "success" {
				t.Errorf("expected result 'success', got '%s'", res.Result)
			}
			if res.Status != "ok" {
				t.Errorf("expected status 'ok', got '%s'", res.Status)
			}
			if res.Database != "connected" {
				t.Errorf("expected database 'connected', got '%s'", res.Database)
			}
			if res.Timestamp == "" {
				t.Errorf("expected non-empty timestamp")
			}
		})

		t.Run("HEAD "+endpoint, func(t *testing.T) {
			resp, err := http.Head(ts.URL + endpoint)
			if err != nil {
				t.Fatalf("failed to make HEAD request to %s: %v", endpoint, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("expected status 200 for HEAD %s, got %d", endpoint, resp.StatusCode)
			}

			if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
				t.Errorf("expected Content-Type application/json; charset=utf-8, got '%s'", ct)
			}

			if cc := resp.Header.Get("Cache-Control"); cc != "no-cache, no-store, must-revalidate" {
				t.Errorf("expected no-cache Cache-Control for HEAD %s, got '%s'", endpoint, cc)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed reading HEAD response body: %v", err)
			}
			if len(body) != 0 {
				t.Errorf("expected empty body for HEAD %s, got %d bytes", endpoint, len(body))
			}
		})
	}
}

func TestCheckEndpointDegraded(t *testing.T) {
	dbPath := filepath.Join("..", "..", "biblia.db")
	store, err := database.NewStore(dbPath)
	if err != nil {
		t.Fatalf("failed to connect to test db: %v", err)
	}
	_ = store.Close() // Close store immediately to induce Ping error

	hClosed := handlers.NewAPIHandler(store)
	rClosed := chi.NewRouter()
	rClosed.Use(chimiddleware.GetHead)
	rClosed.Get("/health", hClosed.Check)
	rClosed.Head("/health", hClosed.Check)

	tsClosed := httptest.NewServer(rClosed)
	defer tsClosed.Close()

	t.Run("GET degraded", func(t *testing.T) {
		resp, err := http.Get(tsClosed.URL + "/health")
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("expected status 503 Service Unavailable for closed DB, got %d", resp.StatusCode)
		}

		var res models.CheckResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res.Status != "degraded" {
			t.Errorf("expected status 'degraded', got '%s'", res.Status)
		}
		if res.Database != "disconnected" {
			t.Errorf("expected database 'disconnected', got '%s'", res.Database)
		}
	})

	t.Run("HEAD degraded", func(t *testing.T) {
		resp, err := http.Head(tsClosed.URL + "/health")
		if err != nil {
			t.Fatalf("failed to make HEAD request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Errorf("expected status 503 for HEAD degraded, got %d", resp.StatusCode)
		}
	})
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

func TestGetBooksHead(t *testing.T) {
	ts, cleanup := setupTestServer(t)
	defer cleanup()

	resp, err := http.Head(ts.URL + "/api/books")
	if err != nil {
		t.Fatalf("failed to make HEAD request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 for HEAD /api/books, got %d", resp.StatusCode)
	}

	if cc := resp.Header.Get("Cache-Control"); cc != "public, max-age=31536000" {
		t.Errorf("expected Cache-Control public header for HEAD /api/books, got '%s'", cc)
	}

	if ct := resp.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Errorf("expected Content-Type application/json; charset=utf-8, got '%s'", ct)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed reading HEAD response body: %v", err)
	}
	if len(body) != 0 {
		t.Errorf("expected empty body for HEAD /api/books, got %d bytes", len(body))
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
