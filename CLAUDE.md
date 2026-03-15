# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Skalkaho is a construction quoting application for small-to-medium contractors. Named after Skalkaho Pass in Montana. Single-tenant: each deployment serves one contractor.

**Tech Stack**: Go backend, HTMX + Alpine.js + Tailwind (CDN) frontend, Svelte 5 estimate builder island, PostgreSQL only

## Build Commands

```bash
# Development
mage dev                # Run development server

# Build
mage build              # Build server + seed binaries to bin/

# Testing
mage test               # Run domain tests

# Code generation
mage sqlc               # Generate repository code from SQL queries

# User management (requires DATABASE_URL)
EMAIL=user@example.com PASSWORD=secret NAME="User Name" mage seed:create
mage seed:list

# Estimate Builder frontend
mage ui:install         # Install npm dependencies
mage ui:build           # Build Svelte bundle
mage ui:watch           # Watch mode

# Other
mage deps               # Download and tidy Go modules
mage clean              # Remove built binaries
```

## Architecture

```
cmd/
├── server/
│   ├── main.go             # Entry point, dependency wiring
│   └── migrations/         # Embedded Goose SQL migrations (001-008)
├── seed/
│   └── main.go             # CLI user management tool
internal/
├── auth/                   # Hashing, sessions (no org/role), context, middleware
├── config/                 # Environment configuration (simplified)
├── database/               # PostgreSQL connection (no SQLite)
├── domain/                 # EstimatePayload types, ResolveMarkup, cost calculation
├── middleware/              # Recover, RequestID, Logger
├── repository/             # sqlc-generated database code
sqlc/queries/               # SQL queries for sqlc
ui/                         # Svelte 5 Estimate Builder
static/                     # Built assets
```

## Domain Model

Core entities: **User** → **Project** (quote container) → **Section** → **Subcategory** → **LineItem**

Supporting: **Client**, **Supplier** → **Material**, **RateCategory** → **Rate**, **Quote** → **QuoteSignature**, **Template** → **TemplateSection** → **TemplateSubcategory** → **TemplateComponentGroup**

No `org_id` on any table. Single-tenant. No roles.

**Markup Resolution**: Per-subcategory overrides → global project defaults. 5 category types: materials, labor, equipment, subs, other. Each can be enabled/disabled per subcategory.

## Authentication

**Auth package** (`internal/auth/`):
- Argon2id password hashing
- Session tokens: crypto/rand generation, SHA-256 hashed storage
- Context helpers: `UserIDFromContext()`, `UserEmailFromContext()`, `UserNameFromContext()`, `IsAuthenticated()`
- Middleware: `SessionMiddleware` (loads session) → `RequireAuth` (redirects to /login)
- No roles, no org_id, no RequireRole

**User management**: Via CLI seed tool only. No `/register` page.

## Key Development Patterns

**Database**: Use sqlc for code generation. Define SQL queries in `sqlc/queries/`, run `mage sqlc`. PostgreSQL only via `DATABASE_URL`.

**Error handling**: Wrap with context, use early returns
```go
if err != nil {
    return fmt.Errorf("creating project: %w", err)
}
```

**HTTP routing** (Go 1.22+):
```go
mux.HandleFunc("GET /projects/{id}", handler.GetProject)
```

**Middleware order**: Recover → RequestID → Logger (global), then SessionMiddleware → RequireAuth (per-route via `protect()`)

**IDs**: TEXT primary keys. Use `uuid.New().String()[:20]` or nanoid.

**Frontend**: HTMX + Alpine.js + Tailwind CDN for all pages. Svelte 5 island for estimate builder only.

## Environment Variables

```bash
ADDR=:8080                          # Server address (default :8080)
ENVIRONMENT=development             # development or production
DATABASE_URL=postgres://...         # PostgreSQL connection string (required)
SESSION_SECRET=<random-string>      # Required in production
SESSION_DURATION=720h               # Session lifetime (default 30 days)
SESSION_COOKIE_NAME=skalkaho_session
POSTMARK_API_KEY=                   # Transactional email (Phase 9)
```

## Reference Documents

- `docs/skalkaho-spec.md` - Authoritative product specification
- `docs/IMPLEMENTATION_PLAN.md` - Greenfield build guide with phased delivery
