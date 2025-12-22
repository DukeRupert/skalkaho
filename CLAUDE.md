# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Skalkaho is a construction quoting SaaS application for small-to-medium contractors. Named after Skalkaho Pass in Montana. Currently in MVP development phase.

**Tech Stack**: Go backend, HTMX + Alpine.js + Tailwind frontend, SQLite (MVP) → PostgreSQL (production)

## Build Commands

```bash
# Database (Docker Compose + Goose)
make db-up              # Start PostgreSQL
make db-migrate         # Run migrations
make db-rollback        # Rollback one migration
make db-reset           # Reset and re-apply all migrations
make db-status          # Check migration status
make db-new             # Create new migration

# Code generation
make sqlc               # Generate repository code from SQL queries
```

## Architecture

```
cmd/server/main.go              # Entry point, dependency wiring
internal/
├── config/                     # Environment configuration
├── database/                   # Connection and migrations
├── domain/                     # Business logic, validation, errors
├── handler/                    # HTTP handlers (by user context, not entity)
│   ├── admin/                  # Contractor admin operations
│   ├── storefront/             # Customer-facing pages
│   └── webhook/                # External integrations
├── middleware/                 # HTTP middleware
├── repository/                 # sqlc-generated database code
├── router/                     # Route definitions
└── templates/                  # html/template files
migrations/                     # Goose SQL migrations
sqlc/queries/                   # SQL queries for sqlc
```

## Domain Model

Core entities: **Settings** (app defaults) → **Job** (quote container) → **Category** (nested up to 3 levels) → **LineItem** (material or labor)

**Surcharge Inheritance**: LineItem → Category → Job hierarchy with two modes:
- `stacking`: All surcharges add together
- `override`: Use most specific (lowest-level) surcharge

## Key Development Patterns

**Handlers**: Organize by user context (`/handler/admin/`), not by entity (`/handler/products/`)

**Database**: Use sqlc for code generation. Define SQL queries in `sqlc/queries/`, run `make sqlc`

**Interfaces**: Define where they're used, not at the repository layer. Accept interfaces, return structs

**Error handling**: Wrap with context, use early returns
```go
if err != nil {
    return fmt.Errorf("creating job: %w", err)
}
```

**HTTP routing** (Go 1.22+):
```go
mux.HandleFunc("GET /products/{id}", handler.GetProduct)
```

**Middleware order**: Panics → RequestID → Logging → Security → RateLimit → Tenant → Route-specific Auth

**Frontend**: HTMX for real-time updates, Alpine.js for lightweight interactivity. UI must update immediately on any data change

**Spam protection**: Honeypot fields + Cloudflare Turnstile

## Reference Documents

- `development/MVP_GUIDE.md` - Complete MVP specification, data model, calculation logic
- `development/GO_STYLE_GUIDE.md` - Comprehensive coding standards for the entire stack
