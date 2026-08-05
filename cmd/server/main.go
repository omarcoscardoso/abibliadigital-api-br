package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"abibliadigital/internal/database"
	"abibliadigital/internal/handlers"
	"abibliadigital/internal/middleware"
)

func main() {
	var logger *slog.Logger
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	} else {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	slog.SetDefault(logger)

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "biblia.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Connecting to SQLite database", "path", dbPath, "mode", "ro")
	store, err := database.NewStore(dbPath)
	if err != nil {
		slog.Error("Database connection error", "error", err)
		os.Exit(1)
	}
	defer store.Close()

	h := handlers.NewAPIHandler(store)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS)
	r.Use(middleware.CacheControl)

	// Health check endpoints for monitoring and probes
	r.Get("/health", h.Check)
	r.Get("/healthz", h.Check)

	// API Routes
	r.Route("/api", func(r chi.Router) {
		r.Get("/check", h.Check)
		r.Get("/health", h.Check)

		r.Get("/books", h.GetBooks)
		r.Get("/books/{abbrev}", h.GetBook)

		r.Get("/versions", h.GetVersions)

		r.Get("/verses/{version}/random", h.GetRandomVerse)
		r.Get("/verses/{version}/{abbrev}/random", h.GetRandomVerse)
		r.Get("/verses/{version}/{abbrev}/{chapter}", h.GetChapter)
		r.Get("/verses/{version}/{abbrev}/{chapter}/{number}", h.GetVerse)
		r.Post("/verses/search", h.Search)

		r.Get("/users/stats", h.UserStats)
		r.Get("/users/{email}", h.GetUser)
		r.Post("/users/password/{email}", h.ResendPassword)
		r.Post("/users", h.CreateUser)
		r.Put("/users/token", h.UpdateToken)
		r.Delete("/users", h.RemoveUser)

		r.Get("/requests/{period}", h.GetRequests)
		r.Get("/requests/amount/{period}", h.GetRequestsNumber)
	})

	// Serve Landing Page & Static Assets
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})
	r.Get("/pt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/pt/index.html")
	})
	r.Get("/en", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/en/index.html")
	})

	// Serve OpenAPI Spec & Interactive API Documentation (Scalar)
	r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/openapi.yaml")
	})
	r.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/docs/index.html")
	})
	r.Get("/docs/*", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/docs/index.html")
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "public"))
	fileServer(r, "/", filesDir)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Server listening", "port", port, "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped cleanly.")
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.Contains(path, ":") {
		panic("fileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	r.Get(path+"*", func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	})
}
