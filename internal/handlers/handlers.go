package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"abibliadigital/internal/database"
	"abibliadigital/internal/models"
)

type APIHandler struct {
	store *database.Store
}

func NewAPIHandler(store *database.Store) *APIHandler {
	return &APIHandler{store: store}
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, models.ErrorResponse{Msg: msg})
}

// GET /api/books
func (h *APIHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.store.GetBooks()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve books")
		return
	}

	type bookListItem struct {
		Abbrev    models.Abbrev `json:"abbrev"`
		Author    string        `json:"author"`
		Chapters  int           `json:"chapters"`
		Group     string        `json:"group"`
		Name      string        `json:"name"`
		Testament string        `json:"testament"`
	}

	result := make([]bookListItem, 0, len(books))
	for _, b := range books {
		result = append(result, bookListItem{
			Abbrev:    b.Abbrev,
			Author:    b.Author,
			Chapters:  b.Chapters,
			Group:     b.Group,
			Name:      b.Name,
			Testament: b.Testament,
		})
	}

	respondJSON(w, http.StatusOK, result)
}

// GET /api/books/:abbrev
func (h *APIHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	abbrev := chi.URLParam(r, "abbrev")
	book, err := h.store.GetBook(abbrev)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error retrieving book")
		return
	}
	if book == nil {
		respondError(w, http.StatusNotFound, "Book not found")
		return
	}

	respondJSON(w, http.StatusOK, book)
}

// GET /api/verses/:version/:abbrev/:chapter
func (h *APIHandler) GetChapter(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	abbrev := chi.URLParam(r, "abbrev")
	chapStr := chi.URLParam(r, "chapter")

	chapNum, err := strconv.Atoi(chapStr)
	if err != nil {
		respondError(w, http.StatusNotFound, "Chapter not found")
		return
	}

	book, err := h.store.GetBook(abbrev)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error checking book")
		return
	}
	if book == nil {
		respondError(w, http.StatusNotFound, "Book not found")
		return
	}

	chapterResp, err := h.store.GetChapter(version, abbrev, chapNum)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error retrieving chapter")
		return
	}
	if chapterResp == nil {
		respondError(w, http.StatusNotFound, "Chapter not found")
		return
	}

	respondJSON(w, http.StatusOK, chapterResp)
}

// GET /api/verses/:version/:abbrev/:chapter/:number
func (h *APIHandler) GetVerse(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	abbrev := chi.URLParam(r, "abbrev")
	chapStr := chi.URLParam(r, "chapter")
	numStr := chi.URLParam(r, "number")

	chapNum, err1 := strconv.Atoi(chapStr)
	numVal, err2 := strconv.Atoi(numStr)
	if err1 != nil || err2 != nil {
		respondError(w, http.StatusNotFound, "Verse not found")
		return
	}

	book, err := h.store.GetBook(abbrev)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error checking book")
		return
	}
	if book == nil {
		respondError(w, http.StatusNotFound, "Book not found")
		return
	}

	verseResp, err := h.store.GetVerse(version, abbrev, chapNum, numVal)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error retrieving verse")
		return
	}
	if verseResp == nil {
		respondError(w, http.StatusNotFound, "Verse not found")
		return
	}

	respondJSON(w, http.StatusOK, verseResp)
}

// GET /api/verses/:version/random & /api/verses/:version/:abbrev/random
func (h *APIHandler) GetRandomVerse(w http.ResponseWriter, r *http.Request) {
	version := chi.URLParam(r, "version")
	abbrev := chi.URLParam(r, "abbrev")

	verseResp, err := h.store.GetRandomVerse(version, abbrev)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error generating random verse")
		return
	}
	if verseResp == nil {
		respondError(w, http.StatusNotFound, "Verse not found")
		return
	}

	respondJSON(w, http.StatusOK, verseResp)
}

// POST /api/verses/search
func (h *APIHandler) Search(w http.ResponseWriter, r *http.Request) {
	var req models.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if req.Version == "" {
		respondError(w, http.StatusNotFound, "Version not found")
		return
	}

	res, err := h.store.Search(req.Version, req.Search)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error searching verses")
		return
	}

	respondJSON(w, http.StatusOK, res)
}

// GET /api/versions
func (h *APIHandler) GetVersions(w http.ResponseWriter, r *http.Request) {
	versions, err := h.store.GetVersions()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Error retrieving versions")
		return
	}

	if versions == nil {
		versions = []models.VersionResponse{}
	}

	respondJSON(w, http.StatusOK, versions)
}

// GET /api/check
func (h *APIHandler) Check(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, models.CheckResponse{Result: "success"})
}

// Transparent User and Request Mocks for Legacy Compatibility
func (h *APIHandler) UserStats(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"lastLogin":        "2026-01-01T00:00:00.000Z",
		"requestsPerMonth": []interface{}{},
	})
}

func (h *APIHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"name":          "User",
		"email":         r.URL.Path,
		"token":         "mock-jwt-token",
		"notifications": true,
	})
}

func (h *APIHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"name":          "User",
		"email":         "user@abibliadigital.com.br",
		"token":         "mock-jwt-token",
		"notifications": true,
	})
}

func (h *APIHandler) UpdateToken(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"name":  "User",
		"email": "user@abibliadigital.com.br",
		"token": "mock-jwt-token",
	})
}

func (h *APIHandler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"msg": "User successfully removed",
	})
}

func (h *APIHandler) ResendPassword(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")
	respondJSON(w, http.StatusOK, map[string]string{
		"msg": "New password successfully sent to email " + email,
	})
}

func (h *APIHandler) GetRequests(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, []interface{}{})
}

func (h *APIHandler) GetRequestsNumber(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]interface{}{
		"total":    0,
		"requests": []interface{}{},
	})
}
