package database

import (
	"testing"
)

func TestOpen_SQLite(t *testing.T) {
	// Test SQLite connection (backward compatibility)
	db, err := Open("", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping SQLite database: %v", err)
	}
}

func TestOpen_PostgreSQL(t *testing.T) {
	// Skip if no DATABASE_URL is set
	t.Skip("Skipping PostgreSQL test - requires DATABASE_URL environment variable")

	// This test can be run manually with:
	// DATABASE_URL=postgres://user:password@localhost:5432/dbname go test -v ./internal/database
}
