package main

import (
	"database/sql"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"

	"github.com/dukerupert/skalkaho/internal/auth"
	"github.com/dukerupert/skalkaho/internal/config"
	"github.com/dukerupert/skalkaho/internal/database"
	authhandler "github.com/dukerupert/skalkaho/internal/handler/auth"
	"github.com/dukerupert/skalkaho/internal/handler/keyboard"
	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/router"
	keyboardtemplates "github.com/dukerupert/skalkaho/internal/templates/keyboard"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	// Load .env file if present (ignore error if not found)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Skalkaho starting", "environment", cfg.Environment)

	// Validate session secret in production
	if cfg.Environment == "production" && cfg.SessionSecret == "" {
		log.Fatal("SESSION_SECRET is required in production")
	}

	// Use a default session secret for development if not set
	sessionSecret := cfg.SessionSecret
	if sessionSecret == "" {
		sessionSecret = "dev-secret-do-not-use-in-production"
		logger.Warn("Using default session secret - DO NOT use in production")
	}

	// Open database (PostgreSQL or SQLite)
	db, err := database.Open(cfg.DatabaseURL, cfg.DatabasePath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() { _ = db.Close() }()

	// Determine dialect for migrations
	dialect := "sqlite3"
	if cfg.DatabaseURL != "" {
		dialect = "postgres"
	}

	// Run migrations
	if err := runMigrations(db, dialect); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repository
	queries := repository.New(db)

	// Initialize template renderer
	renderer, err := keyboardtemplates.NewRenderer()
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
	handler := keyboard.NewHandler(queries, renderer, logger, cfg)
	authHandler := authhandler.NewHandler(db, renderer, sessionManager)

	// Setup router
	mux := http.NewServeMux()
	router.Register(mux, handler, authHandler, sessionManager)

	// Apply middleware
	httpHandler := middleware.Chain(mux,
		middleware.Recover,
		middleware.RequestID,
		middleware.Logger(logger),
	)

	// Start server
	logger.Info("Starting server", "addr", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, httpHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func runMigrations(db *sql.DB, dialect string) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
