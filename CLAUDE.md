# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Skalkaho is a construction quoting SaaS application for small-to-medium contractors. Named after Skalkaho Pass in Montana. Multi-tenant with authentication, transitioning from MVP to production SaaS.

**Tech Stack**: Go backend, HTMX + Alpine.js + Tailwind (CDN) frontend, PostgreSQL (primary) / SQLite (dev fallback)

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

# Integration tests (auth)
go test ./tests/integration/... -v

# Code generation
make sqlc               # Generate repository code from SQL queries

# Database
make db-migrate         # Run migrations
make db-rollback        # Roll back last migration
make db-status          # Show migration status
make db-reset           # Delete DB and re-run migrations
```

## Architecture

```
cmd/server/
├── main.go             # Entry point, dependency wiring, session manager init
└── migrations/         # Embedded Goose SQL migrations (19 total)
internal/
├── auth/               # Authentication: hashing, sessions, context, middleware
├── config/             # Environment configuration
├── database/           # PostgreSQL + SQLite connection (auto-detects via DATABASE_URL)
├── domain/             # Business logic, validation, surcharge calculation
├── handler/
│   ├── auth/           # Login, register, logout handlers
│   └── keyboard/       # Quote/job handlers, tenant context helpers
├── middleware/          # Recover, RequestID, Logger
├── repository/          # sqlc-generated database code
├── router/              # Route definitions with auth protection
├── service/
│   ├── claude/          # AI-powered price matching (Anthropic API)
│   └── excel/           # Excel import service
├── templates/keyboard/  # html/template files (layouts, pages, partials)
└── testutil/            # Test helpers for auth, fixtures, database setup
migrations/              # Source Goose SQL migrations
sqlc/queries/            # SQL queries for sqlc
tests/
├── integration/         # Auth flow integration tests
└── AUTH_TEST_GUIDE.md   # Guide for running auth tests
```

## Domain Model

Core entities: **Organization** (tenant) → **User** (member) → **Job** (quote container) → **Category** (nested up to 3 levels) → **LineItem** (material or labor)

All tenant-scoped tables include `org_id` for multi-tenancy isolation.

**Surcharge Inheritance**: LineItem → Category → Job hierarchy with two modes:
- `stacking`: All surcharges add together (Job 15% + Category 10% + Item 5% = 30%)
- `override`: Use most specific (lowest-level) surcharge (Item 5% wins)

## Authentication & Multi-Tenancy

**Auth package** (`internal/auth/`):
- Argon2id password hashing
- Session tokens: crypto/rand generation, SHA-256 hashed storage
- Context helpers: `OrgIDFromContext()`, `UserIDFromContext()`, etc.
- Middleware: `SessionMiddleware` (loads session) → `RequireAuth` (redirects to /login) → `RequireRole` (role check)

**Tenant isolation** (`internal/handler/keyboard/tenant.go`):
- `GetOrgID()` / `MustGetOrgID()` extract org_id from session context
- All queries scoped by org_id

**Routing**: Protected routes wrapped with `protect()` helper (SessionMiddleware + RequireAuth). Public routes: login, register, logout, health, e-signature links.

**User roles**: owner, admin, member

## Key Development Patterns

**Handlers**: Organized in `/handler/keyboard/` (app) and `/handler/auth/` (authentication)

**Database**: Use sqlc for code generation. Define SQL queries in `sqlc/queries/`, run `make sqlc`. PostgreSQL detected via `DATABASE_URL`; falls back to SQLite via `DATABASE_PATH`.

**Templates**: Located in `internal/templates/keyboard/` with `layouts/`, `pages/`, and `partials/` subdirectories.

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

**Middleware order**: Recover → RequestID → Logger (global), then SessionMiddleware → RequireAuth (per-route via `protect()`)

**Frontend**: HTMX for real-time updates, Alpine.js for UI state (expand/collapse, forms), Tailwind via CDN

## Environment Variables

```bash
# Server
ADDR=:8080                          # Server address (default :8080)
ENVIRONMENT=development             # Environment mode (development|production)

# Database (set DATABASE_URL for PostgreSQL, otherwise SQLite fallback)
DATABASE_URL=postgres://...         # PostgreSQL connection string (primary)
DATABASE_PATH=quotes.db             # SQLite path (dev fallback, default quotes.db)

# Authentication
SESSION_SECRET=<random-string>      # Required in production
SESSION_DURATION=720h               # Session lifetime (default 30 days)
SESSION_COOKIE_NAME=skalkaho_session

# Services
ANTHROPIC_API_KEY=                  # AI price matching
AUTO_APPROVE_THRESHOLD=0.9          # Price match auto-approval threshold
PRICE_IMPORT_TOKEN=                 # Secret token for price import access
```

## Reference Documents

- `development/MVP_GUIDE.md` - Complete MVP specification, data model, calculation logic
- `development/GO_STYLE_GUIDE.md` - Comprehensive coding standards for the entire stack
- `development/ROADMAP.md` - SaaS development roadmap
- `tests/AUTH_TEST_GUIDE.md` - Guide for running authentication tests
