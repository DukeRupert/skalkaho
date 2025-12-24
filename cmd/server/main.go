package main

import (
	"database/sql"
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"

	"github.com/dukerupert/skalkaho/internal/config"
	"github.com/dukerupert/skalkaho/internal/handler/quote"
	"github.com/dukerupert/skalkaho/internal/middleware"
	"github.com/dukerupert/skalkaho/internal/repository"
	"github.com/dukerupert/skalkaho/internal/router"
	"github.com/dukerupert/skalkaho/internal/templates"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("Skalkaho starting", "environment", cfg.Environment)

	// Open database
	db, err := sql.Open("sqlite3", cfg.DatabasePath+"?_foreign_keys=on")
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

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

	// Initialize handlers
	quoteHandler := quote.NewHandler(queries, renderer, logger)

	// Setup router
	mux := http.NewServeMux()
	router.Register(mux, quoteHandler)

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

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
