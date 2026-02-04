package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

// Open creates a database connection.
// If databaseURL is provided (starts with postgres://), uses PostgreSQL.
// Otherwise falls back to SQLite with the provided path.
func Open(databaseURL, databasePath string) (*sql.DB, error) {
	var db *sql.DB
	var err error

	if databaseURL != "" && strings.HasPrefix(databaseURL, "postgres://") {
		// PostgreSQL connection
		db, err = sql.Open("pgx", databaseURL)
		if err != nil {
			return nil, fmt.Errorf("opening postgres database: %w", err)
		}
	} else {
		// SQLite connection (backward compatibility)
		db, err = sql.Open("sqlite3", databasePath+"?_foreign_keys=on")
		if err != nil {
			return nil, fmt.Errorf("opening sqlite database: %w", err)
		}
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return db, nil
}

// Dialect returns the database dialect based on the driver name.
func Dialect(db *sql.DB) string {
	// This is a simple heuristic - in production you might want a more robust solution
	if db.Driver() == nil {
		return "sqlite3"
	}
	return "postgres"
}
