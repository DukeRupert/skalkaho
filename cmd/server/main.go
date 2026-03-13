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
	authhandler "github.com/dukerupert/skalkaho/internal/handler/auth"
	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/router"
	"github.com/dukerupert/skalkaho/internal/templates"
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

	// Initialize template renderer
	renderer, err := templates.NewRenderer()
	if err != nil {
		log.Fatalf("Failed to initialize templates: %v", err)
	}

	// Initialize session manager
	sessionManager := auth.NewSessionManager(
		db,
		sessionSecret,
		cfg.SessionDuration,
		cfg.SessionCookieName,
		cfg.SecureCookies,
	)

	// Initialize handlers
	authHandler := authhandler.NewHandler(queries, renderer, sessionManager, logger)

	// Setup router
	mux := http.NewServeMux()
	router.Register(mux, authHandler, sessionManager)

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

	return goose.Up(db, "migrations")
}
