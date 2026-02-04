package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"

	_ "github.com/mattn/go-sqlite3"
)

// TestDB creates an isolated SQLite test database with migrations applied.
// Returns the database connection and a cleanup function.
func TestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	// Create temporary directory for test database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("failed to ping test database: %v", err)
	}

	// Run migrations
	migrationsDir := findMigrationsDir(t)

	// Set goose to use sqlite3 dialect
	if err := goose.SetDialect("sqlite3"); err != nil {
		db.Close()
		t.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		db.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Errorf("error closing test database: %v", err)
		}
	}

	return db, cleanup
}

// findMigrationsDir locates the migrations directory relative to the test.
// This works from any test location in the project.
// It looks for cmd/server/migrations (embedded migrations) or migrations/ (source migrations).
func findMigrationsDir(t *testing.T) string {
	t.Helper()

	// Try to find migrations directory by walking up from current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// Walk up the directory tree looking for migrations directory
	dir := cwd
	for {
		// Try cmd/server/migrations first (embedded migrations)
		migrationsPath := filepath.Join(dir, "cmd", "server", "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			return migrationsPath
		}

		// Try migrations/ directory (source migrations)
		migrationsPath = filepath.Join(dir, "migrations")
		if _, err := os.Stat(migrationsPath); err == nil {
			return migrationsPath
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root without finding migrations
			t.Fatalf("migrations directory not found starting from %s", cwd)
		}
		dir = parent
	}
}

// ExecSQL executes SQL statements on the test database.
// Useful for setting up test data.
func ExecSQL(t *testing.T, db *sql.DB, query string, args ...interface{}) {
	t.Helper()

	_, err := db.ExecContext(context.Background(), query, args...)
	if err != nil {
		t.Fatalf("failed to execute SQL: %v\nQuery: %s", err, query)
	}
}

// QueryRowSQL executes a query that returns a single row.
func QueryRowSQL(t *testing.T, db *sql.DB, query string, args ...interface{}) *sql.Row {
	t.Helper()
	return db.QueryRowContext(context.Background(), query, args...)
}

// MustCount executes a COUNT query and returns the result.
// Fails the test if the query fails.
func MustCount(t *testing.T, db *sql.DB, query string, args ...interface{}) int64 {
	t.Helper()

	var count int64
	err := db.QueryRowContext(context.Background(), query, args...).Scan(&count)
	if err != nil {
		t.Fatalf("failed to execute count query: %v\nQuery: %s", err, query)
	}
	return count
}

// AssertRowExists verifies that a row exists matching the given query.
func AssertRowExists(t *testing.T, db *sql.DB, query string, args ...interface{}) {
	t.Helper()

	count := MustCount(t, db, fmt.Sprintf("SELECT COUNT(*) FROM (%s)", query), args...)
	if count == 0 {
		t.Errorf("expected row to exist, but got 0 rows\nQuery: %s", query)
	}
}

// AssertRowNotExists verifies that no rows exist matching the given query.
func AssertRowNotExists(t *testing.T, db *sql.DB, query string, args ...interface{}) {
	t.Helper()

	count := MustCount(t, db, fmt.Sprintf("SELECT COUNT(*) FROM (%s)", query), args...)
	if count > 0 {
		t.Errorf("expected 0 rows, but got %d rows\nQuery: %s", count, query)
	}
}
