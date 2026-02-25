package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
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

	// Seed demo user if enabled
	if err := seedDemoUser(db, queries, cfg, logger); err != nil {
		log.Fatalf("Failed to seed demo user: %v", err)
	}

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
	handler := keyboard.NewHandler(db, queries, renderer, logger, cfg)
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

func seedDemoUser(db *sql.DB, queries *repository.Queries, cfg *config.Config, logger *slog.Logger) error {
	if !cfg.SeedDemoUser {
		return nil
	}

	ctx := context.Background()

	// Idempotent: check if demo org already exists
	_, err := queries.GetOrganizationBySubdomain(ctx, "demo")
	if err == nil {
		logger.Info("Demo user already exists, skipping seed")
		return nil
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("checking for demo org: %w", err)
	}

	// Hash demo password
	passwordHash, err := auth.HashPassword("demo1234")
	if err != nil {
		return fmt.Errorf("hashing demo password: %w", err)
	}

	// Create org, settings, and user in a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	qtx := queries.WithTx(tx)

	org, err := qtx.CreateOrganization(ctx, repository.CreateOrganizationParams{
		Name:      "Demo Construction Co",
		Subdomain: "demo",
		Plan:      "free",
		Status:    "active",
	})
	if err != nil {
		return fmt.Errorf("creating demo org: %w", err)
	}

	_, err = qtx.CreateSettings(ctx, repository.CreateSettingsParams{
		OrgID:                   org.ID,
		DefaultSurchargeMode:    "stacking",
		DefaultSurchargePercent: 0,
	})
	if err != nil {
		return fmt.Errorf("creating demo settings: %w", err)
	}

	_, err = qtx.CreateUser(ctx, repository.CreateUserParams{
		OrgID:        org.ID,
		Email:        "demo@skalkaho.com",
		PasswordHash: passwordHash,
		Name:         "Demo User",
		Role:         "owner",
		Status:       "active",
	})
	if err != nil {
		return fmt.Errorf("creating demo user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	logger.Info("Demo user seeded", "email", "demo@skalkaho.com", "password", "demo1234")
	return nil
}
