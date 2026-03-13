package main

import (
	"database/sql"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"github.com/dukerupert/skalkaho/internal/auth"
	"github.com/dukerupert/skalkaho/internal/config"
	"github.com/dukerupert/skalkaho/internal/database"
	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Skalkaho starting", "environment", cfg.Environment)

	// Validate session secret in production
	if cfg.Environment == "production" && cfg.SessionSecret == "" {
		log.Fatal("SESSION_SECRET is required in production")
	}

	sessionSecret := cfg.SessionSecret
	if sessionSecret == "" {
		sessionSecret = "dev-secret-do-not-use-in-production"
		logger.Warn("Using default session secret - DO NOT use in production")
	}

	// Open database (PostgreSQL only)
	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Run migrations
	if err := runMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repository
	queries := repository.New(db)
	_ = queries // Used by handlers in later phases

	// Initialize session manager
	sessionManager := auth.NewSessionManager(
		db,
		sessionSecret,
		cfg.SessionDuration,
		cfg.SessionCookieName,
		cfg.SecureCookies,
	)
	_ = sessionManager // Used by router in later phases

	// Setup router with health check
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// Static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Apply middleware
	handler := middleware.Chain(mux,
		middleware.Recover,
		middleware.RequestID,
		middleware.Logger(logger),
	)

	// Start server
	logger.Info("Starting server", "addr", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func runMigrations(db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
