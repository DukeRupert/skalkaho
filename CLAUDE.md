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
│   └── migrations/         # Embedded Goose SQL migrations (001-017)
├── seed/
│   └── main.go             # CLI user management tool
internal/
├── auth/                   # Hashing, sessions (no org/role), context, middleware
├── config/                 # Environment configuration (simplified)
├── database/               # PostgreSQL connection (no SQLite)
├── domain/                 # EstimatePayload types, ResolveMarkup, cost calculation
├── handler/
│   ├── app/                # App handlers (projects, clients, materials, imports, etc.)
│   └── auth/               # Login/logout handlers
├── middleware/              # Recover, RequestID, Logger
├── repository/             # sqlc-generated database code
├── worker/                 # Background job harness (poll loop, panic recovery)
│   └── importer/           # Material import pipeline (Excel parser, Claude matcher)
sqlc/queries/               # SQL queries for sqlc
ui/                         # Svelte 5 Estimate Builder
static/                     # Built assets
```

## Domain Model

Core entities: **User** → **Project** (quote container) → **Section** → **Subcategory** → **LineItem**

Supporting: **Client**, **Supplier** → **Material**, **RateCategory** → **Rate**, **Quote** → **QuoteSignature**, **Template** → **TemplateSection** → **TemplateSubcategory** → **TemplateComponentGroup**

Supply: **SupplyList** (scoped to Quote, optionally to Supplier) → **SupplyListItem** (snapshot of line item + material at generation time)

Materials: **MaterialsList** (scoped to Quote, optionally to Section) — metadata only, items queried live from quote at render time

Import: **MaterialImportJob** → **MaterialImportFlag** (linked to Supplier and optionally to a candidate Material)

QuickBooks: **QBConnection** (single-row), **Invoice** → **InvoicePayment**. Clients have optional `qb_customer_id`/`qb_customer_name` for manual QBO linking.

No `org_id` on any table. Single-tenant. No roles.

**Markup Resolution**: Per-subcategory overrides → global project defaults. 5 category types: materials, labor, equipment, subs, other. Each can be enabled/disabled per subcategory.

## Background Worker

A single goroutine polls `material_import_jobs` every 5 seconds using `FOR UPDATE SKIP LOCKED` (safe for future multi-worker scaling). Wired in `main.go` with `signal.NotifyContext` for clean shutdown on SIGTERM.

**Material Import Pipeline** (`internal/worker/importer/`):
1. Parse Excel file (auto-detect columns from header row, normalize units)
2. Split rows into 3 tiers by `supplier_code` presence/match
3. Tier 1: code matches existing material → update price (skip if `price_source = 'Manual'`)
4. Tier 2: code present but no match → normalize name via Claude → create new material
5. Tier 3: no code → batch to Claude API for fuzzy name matching + name normalization → auto-update (≥0.98), flag for review (0.75–0.97), or create new (<0.75)
6. All DB writes in a single transaction
7. HTMX polls `/imports/{id}/status` every 2s; flag review queue rendered inline

**Material Naming Convention**: All new materials created by the import pipeline are normalized to `Category, Dimensions/Weight, Details` format via Claude. Category comes first so searching by what-it-is returns all sizes. Examples: `Stud, 2x4x8, Douglas Fir` · `Plywood, 4x8, 3/4 CDX` · `Concrete Mix, 80lb, Quikrete` · `Screw, #10x3, GRK RSS`. Existing materials (Tier 1 price updates) keep their current names.

## Supply Lists

Persistent supply lists generated from a signed quote. Lists are snapshots — subsequent quote/estimate edits do not affect existing lists. Three generation modes:

- **All suppliers** — one list covering every supplier + ad hoc items grouped under "No Supplier"
- **Per supplier (all)** — one list per supplier in a single transaction, excludes ad hoc items
- **Per supplier (individual)** — owner picks a specific supplier from a dropdown

**Data flow**: `line_items` → JOIN `materials` → JOIN `suppliers` → snapshot into `supply_list_items`. Each item stores snapshotted `material_name`, `supplier_name`, `supplier_code`, `unit_price`, and computed `extended_price`. Supplier name is stored as a string (not FK) so the snapshot is independent of future supplier record changes.

**Key constraints**: Generation gated to `quote.status = 'signed'`. `supply_lists.quote_id` uses `ON DELETE RESTRICT` (can't delete a quote that has supply lists). Items CASCADE delete with their parent list.

**Routes**: All under `/quotes/{id}/supply-lists/...`. Accessible from the quote row on the project overview page (signed quotes show a "Supply Lists" action link).

**Templates**: `supply_lists.html` (index), `supply_list_show.html` (detail/print view with `@media print` styles). Dates and notes are editable inline via HTMX.

## Materials Lists

Job-site materials lists generated from a signed quote. Unlike supply lists, these are **live** — items are queried from the quote at render time, not snapshotted. No prices or supplier info; just item name, quantity, and unit for the foreman. Three generation modes:

- **Entire quote** — one list covering all categories and subcategories
- **Per category (all)** — one list per top-level section
- **Single category** — owner picks a specific section from a dropdown

**Data model**: `materials_lists` stores only metadata (quote_id, optional section_id, label, notes). Items are queried live via `GetMaterialsListItems` (full quote, grouped by section → subcategory) or `GetMaterialsListItemsBySection` (single section, grouped by subcategory). A nullable `supply_list_id` FK is reserved for future supply association features.

**Key constraints**: Generation gated to `quote.status = 'signed'`. `materials_lists.quote_id` uses `ON DELETE RESTRICT`. `section_id` uses `ON DELETE RESTRICT`. `supply_list_id` uses `ON DELETE SET NULL`.

**Routes**: All under `/quotes/{id}/materials-lists/...`. Accessible from the quote row on the project overview page alongside Supply Lists.

**Templates**: `materials_lists.html` (index), `materials_list_show.html` (detail/print view with `@media print` styles, section → subcategory hierarchy). Notes editable inline via HTMX.

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
ANTHROPIC_API_KEY=                  # Claude API for material import fuzzy matching (Tier 3)
QBO_CLIENT_ID=                      # QuickBooks OAuth client ID
QBO_CLIENT_SECRET=                  # QuickBooks OAuth client secret
QBO_REDIRECT_URI=                   # QuickBooks OAuth redirect URI (must match Intuit dashboard)
QBO_WEBHOOK_TOKEN=                  # QuickBooks webhook verifier token (HMAC-SHA256)
QBO_SANDBOX=true                    # Use QBO sandbox URLs (omit or "false" for production)
```

## Reference Documents

- `docs/skalkaho-spec.md` - Authoritative product specification
- `docs/IMPLEMENTATION_PLAN.md` - Greenfield build guide with phased delivery
- `docs/material-import-feature-plan.md` - Material import feature design (Excel upload, Claude matching, flag review)
- `docs/quickbooks-integration-plumbing.md` - QuickBooks Online integration plan (OAuth, invoicing, payment sync)
- `docs/supply-list-feature-plan.md` - Supply list feature design (snapshot generation, print view, delivery tracking)
- `docs/materials-list-feature-plan.md` - Materials list feature design (live queries, section grouping, foreman print view)
