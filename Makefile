.PHONY: dev build test sqlc db-migrate db-rollback db-status db-new clean

# Development
dev:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

# Testing
test:
	go test ./internal/domain/... -v

test-short:
	go test ./internal/domain/...

# Database
DB_PATH ?= quotes.db

db-migrate:
	go run ./cmd/migrate up

db-rollback:
	go run ./cmd/migrate down

db-status:
	go run ./cmd/migrate status

db-reset:
	rm -f $(DB_PATH)
	go run ./cmd/migrate up

# Code generation
sqlc:
	sqlc generate

# Cleanup
clean:
	rm -f bin/server
	rm -f $(DB_PATH)

# Install dependencies
deps:
	go mod download
	go mod tidy
