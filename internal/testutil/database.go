package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// TestDB creates an isolated PostgreSQL test database with migrations applied.
// Returns the database connection and a cleanup function.
func TestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()

	ctx := context.Background()

	// Start PostgreSQL container
	pgContainer, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	// Get connection string
	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("failed to get connection string: %v", err)
	}

	// Open database connection
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("failed to open test database: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		_ = db.Close()
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("failed to ping test database: %v", err)
	}

	// Run migrations
	migrationsDir := findMigrationsDir(t)

	// Set goose to use postgres dialect
	if err := goose.SetDialect("postgres"); err != nil {
		_ = db.Close()
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		_ = db.Close()
		_ = pgContainer.Terminate(ctx)
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Errorf("error closing test database: %v", err)
		}
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Errorf("error terminating postgres container: %v", err)
		}
	}

	return db, cleanup
}

// findMigrationsDir locates the migrations directory relative to the test.
// This works from any test location in the project.
// It looks for migrations/ (source migrations) which contains PostgreSQL migrations.
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
		// Try migrations/ directory (source PostgreSQL migrations)
		migrationsPath := filepath.Join(dir, "migrations")
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
