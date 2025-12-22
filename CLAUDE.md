# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Skalkaho is a construction quoting SaaS application for small-to-medium contractors. Named after Skalkaho Pass in Montana. Currently in MVP development phase.

**Tech Stack**: Go backend, HTMX + Alpine.js + Tailwind (CDN) frontend, SQLite (MVP) → PostgreSQL (production)

## Build Commands

```bash
# Development
make dev                # Run development server
go run ./cmd/server     # Alternative: run directly

# Build
make build              # Build binary to bin/server
go build -o bin/server ./cmd/server

# Testing
make test               # Run domain tests
go test ./internal/domain/... -v

# Code generation
make sqlc               # Generate repository code from SQL queries

# Database
make db-reset           # Delete quotes.db and re-run migrations
```

## Architecture

```
cmd/server/
├── main.go             # Entry point, dependency wiring
└── migrations/         # Embedded Goose SQL migrations
internal/
├── config/             # Environment configuration
├── database/           # SQLite connection
├── domain/             # Business logic, validation, surcharge calculation
├── handler/quote/      # HTTP handlers for quotes, categories, line items
├── middleware/         # Recover, RequestID, Logger
├── repository/         # sqlc-generated database code
├── router/             # Route definitions
└── templates/          # html/template files (layouts, pages, partials)
migrations/             # Source Goose SQL migrations
sqlc/queries/           # SQL queries for sqlc
```

## Domain Model

Core entities: **Settings** (app defaults) → **Job** (quote container) → **Category** (nested up to 3 levels) → **LineItem** (material or labor)

**Surcharge Inheritance**: LineItem → Category → Job hierarchy with two modes:
- `stacking`: All surcharges add together (Job 15% + Category 10% + Item 5% = 30%)
- `override`: Use most specific (lowest-level) surcharge (Item 5% wins)

## Key Development Patterns

**Handlers**: Organized in `/handler/quote/` for the single-user MVP context

**Database**: Use sqlc for code generation. Define SQL queries in `sqlc/queries/`, run `make sqlc`

**Templates**: Each page template (jobs_list, job, settings) is self-contained with full HTML structure. Partials for category and line_item.

**Error handling**: Wrap with context, use early returns
```go
if err != nil {
    return fmt.Errorf("creating job: %w", err)
}
```

**HTTP routing** (Go 1.22+):
```go
mux.HandleFunc("GET /jobs/{id}", handler.GetJob)
```

**Middleware order**: Recover → RequestID → Logger

**Frontend**: HTMX for real-time updates, Alpine.js for UI state (expand/collapse, forms), Tailwind via CDN

## Environment Variables

```bash
ADDR=:8080              # Server address (default :8080)
DATABASE_PATH=quotes.db # SQLite database path (default quotes.db)
ENVIRONMENT=development # Environment mode
```

## Reference Documents

- `development/MVP_GUIDE.md` - Complete MVP specification, data model, calculation logic
- `development/GO_STYLE_GUIDE.md` - Comprehensive coding standards for the entire stack
