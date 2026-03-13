# Skalkaho — Greenfield Implementation Plan

> Derived from `docs/skalkaho-spec.md` · March 2026
> This document is the build guide for an implementing agent.

---

## Table of Contents

1. [Context & Decisions](#1-context--decisions)
2. [What to Keep vs. Rebuild](#2-what-to-keep-vs-rebuild)
3. [Database Schema](#3-database-schema)
4. [Build Order](#4-build-order)
   - [Phase 0: Project scaffold + schema + CLI seed tool](#phase-0-project-scaffold--schema--cli-seed-tool)
   - [Phase 1: Auth (simplified)](#phase-1-auth-simplified)
   - [Phase 2: Navigation shell](#phase-2-navigation-shell)
   - [Phase 3: Projects module](#phase-3-projects-module)
   - [Phase 4: Clients module](#phase-4-clients-module)
   - [Phase 5: Materials Database](#phase-5-materials-database)
   - [Phase 6: Labor & Equipment Rates](#phase-6-labor--equipment-rates)
   - [Phase 7: Estimate Builder (port)](#phase-7-estimate-builder-port)
   - [Phase 8: Job Overview](#phase-8-job-overview)
   - [Phase 9: Quote & Signature system](#phase-9-quote--signature-system)
5. [File Structure (target)](#5-file-structure-target)
6. [Conventions](#6-conventions)

---

## 1. Context & Decisions

This is a **greenfield rebuild**. The existing codebase has accumulated significant legacy code (multi-tenancy, old quoting system, price imports) that does not match the spec. Starting fresh from the spec is cleaner than pruning.

### Architectural decisions (final, do not revisit)

| Decision | Detail |
|----------|--------|
| **Single-tenant** | No `org_id` on any table. No organizations table. No roles. Each deployment serves one contractor. |
| **User management** | CLI tool seeds users. No `/register` page. Multiple users allowed but managed via CLI, not UI. |
| **Legacy quoting system** | Dead code. The old `jobs → categories → line_items` with stacking/override surcharges is dropped entirely. The Svelte Estimate Builder is the only quoting interface. |
| **Theme** | Light Tailwind theme. Ignore the spec's dark GitHub-style design tokens. |
| **Timeline** | No rush. Build it right. |
| **Database** | PostgreSQL only. Drop SQLite fallback. |
| **IDs** | TEXT primary keys using nanoid (client-generated for estimate builder entities, server-generated elsewhere). |
| **Frontend** | HTMX + Alpine.js + Tailwind CDN for all pages except the Estimate Builder, which is a Svelte 5 island. |

### Reference spec

The authoritative spec is `docs/skalkaho-spec.md`. This plan references section numbers from that document (e.g., "spec §2.3" = Section 2.3 of the spec).

---

## 2. What to Keep vs. Rebuild

### Pre-port audit (completed)

Before porting, an audit was performed on all files listed below to identify multi-tenant assumptions that could break silently. Key findings:

| File | Org/Role refs | Action |
|------|---------------|--------|
| `auth/auth.go` | 0 | Port as-is |
| `auth/session.go` | 5 OrgID + 1 Role in Session struct, org validation in CreateSession (line ~73) | Remove OrgID & Role from Session struct, remove org validation |
| `auth/context.go` | `OrgIDFromContext()`, `UserRoleFromContext()` | Delete both functions; keep UserID/Email/Name/Session helpers |
| `auth/middleware.go` | `RequireRole()` (lines 67–89) | Delete function entirely |
| `domain/estimate.go` | 0 | Port as-is |
| `handler/keyboard/estimate_api.go` | 18+ `GetOrgID()` calls, all query params include OrgID | Remove all GetOrgID calls, update all query call signatures |
| `handler/keyboard/tenant.go` | `GetOrgID()`, `MustGetOrgID()`, `GetUserRole()` — 265 downstream call sites | Delete entire file |
| `config/config.go` | 5 legacy fields | Remove DatabasePath, SeedDemoUser, PriceImportToken, AutoApproveThreshold, AnthropicAPIKey |
| `database/database.go` | SQLite fallback block (lines 25–30), `Dialect()` function | Remove SQLite code, PG-only |
| `ui/src/**` | 0 explicit org_id in payloads | Route renaming only (`{jobID}` → `{projectID}`) |
| `sqlc/queries/*.sql` | 132 org_id refs across ~50 queries | All rewritten from scratch in greenfield |

**Critical note for Phase 7:** The estimate API handler (`estimate_api.go`) has the heaviest multi-tenant coupling: every query call passes `OrgID` as a parameter struct field (e.g., `repository.GetJobParams{ID: jobID, OrgID: orgID}`). These aren't just search-and-replace deletions — the sqlc-generated parameter structs change shape when the queries are rewritten without `org_id`. The implementing agent must rewrite the handler against the new query signatures, not attempt to patch the old one.

### Port from existing codebase (adapt, don't copy verbatim)

| Component | Source | Changes needed |
|-----------|--------|----------------|
| **Svelte Estimate Builder** | `ui/` (12 files) | Remove all `org_id` references from API calls. Rename `job` → `project` in API URLs (`/api/estimate/{projectID}`). No changes to component logic — markup engine, undo, autosave, autocomplete all stay. |
| **Markup engine (Go)** | `internal/domain/estimate.go` | Keep `ResolveMarkup()`, `EstimatePayload`, all Estimate* structs. Remove `org_id` from all struct fields. Rename references from "job" to "project" where they appear in JSON tags or field names. |
| **Auth primitives** | `internal/auth/auth.go` | Keep `HashPassword()`, `CheckPassword()`, `GenerateSessionToken()`, `HashSessionToken()`. These are pure functions with no org dependency. |
| **Session manager** | `internal/auth/session.go` | Simplify: remove `OrgID` and `Role` from `Session` struct and `CreateSessionParams`. Remove org validation from `CreateSession()`. Remove role population from `ValidateSession()`. |
| **Auth middleware** | `internal/auth/middleware.go` | Keep `SessionMiddleware` and `RequireAuth`. Delete `RequireRole` entirely. |
| **Auth context** | `internal/auth/context.go` | Delete `OrgIDFromContext()` and `UserRoleFromContext()`. Keep `UserIDFromContext()`, `UserEmailFromContext()`, `UserNameFromContext()`, `IsAuthenticated()`, `SessionFromContext()`. |
| **Middleware** | `internal/middleware/` | Keep Recover, RequestID, Logger as-is. |
| **Config** | `internal/config/config.go` | Remove `DatabasePath`, `SeedDemoUser`, `PriceImportToken`, `AutoApproveThreshold`, `AnthropicAPIKey`. Keep session config, `ADDR`, `ENVIRONMENT`, `DATABASE_URL`. Add `POSTMARK_API_KEY`. |
| **Database** | `internal/database/` | PostgreSQL only. Remove SQLite detection and fallback. Require `DATABASE_URL`; fail startup if missing. |
| **Estimate API handler** | `internal/handler/keyboard/estimate_api.go` | Rewrite (not patch) against new sqlc query signatures. Port `GetEstimateJSON`, `SaveEstimateJSON`, `buildEstimatePayload`, validation logic. Remove all `orgID`/`GetOrgID()` calls. Change `job` → `project` in route params. Autocomplete data source changes from `item_templates` to `materials` + `rates` tables. |

### Build new

| Component | Notes |
|-----------|-------|
| **CLI seed tool** | `cmd/seed/main.go` — creates users via CLI flags |
| **Fresh migrations** | New `001_*.sql` through `00N_*.sql`. See §3 below. |
| **sqlc queries** | All new. See each phase for query lists. |
| **Projects module** | Dashboard, CRUD, status workflow, stats. Spec §4.1. |
| **Clients module** | CRUD with search, expandable rows. Spec §4.2. |
| **Materials Database** | Supplier-organized price list with tabs. Spec §4.3. New entity not in old codebase. |
| **Rates module** | Category-organized rate cards with tabs. Spec §4.4. New entity not in old codebase. |
| **Job Overview** | Project summary page with cost breakdown, quote history. Spec §4.6. |
| **Quote system** | Token-based quotes, `/q/{token}` public page, Postmark email, e-signature. Spec §6. |
| **Navigation sidebar** | Fixed 220px sidebar with 7 pages. Spec §1 "Navigation Model". |
| **All templates** | Fresh html/template files. Light Tailwind theme. |

### Drop entirely (do not port)

- `internal/handler/keyboard/jobs.go` (legacy job/category/line-item handlers)
- `internal/handler/keyboard/category.go`
- `internal/handler/keyboard/jobs_api.go` (legacy JSON API)
- `internal/handler/keyboard/price_import.go`
- `internal/handler/keyboard/item_templates.go`
- `internal/handler/keyboard/job_item_types.go`
- `internal/handler/keyboard/settings.go`
- `internal/handler/keyboard/estimates.go` (old estimate snapshots)
- `internal/handler/keyboard/signatures.go` (old signature flow — rebuilt in Phase 9)
- `internal/service/claude/` (AI price matching — Phase 2 of spec)
- `internal/service/excel/` (Excel import — Phase 2 of spec)
- `internal/domain/domain.go` (old surcharge calculation, `Job`, `Category`, `LineItem` structs)
- `frontend/` directory (old quote editor)
- All legacy templates except as structural reference
- All old migrations
- All old sqlc queries

---

## 3. Database Schema

All migrations use Goose format (`-- +goose Up` / `-- +goose Down`), embedded in `cmd/server/migrations/`.

### Migration 001: Users and Sessions

```sql
-- +goose Up
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid()::text,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT NOT NULL UNIQUE,
    user_agent TEXT,
    ip_address TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token_hash ON sessions(token_hash);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);
```

### Migration 002: Clients

```sql
-- +goose Up
CREATE TABLE clients (
    id TEXT PRIMARY KEY,
    company_name TEXT NOT NULL,
    contact_name TEXT,
    email TEXT,
    phone TEXT,
    address TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### Migration 003: Projects

```sql
-- +goose Up
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    client_id TEXT REFERENCES clients(id) ON DELETE SET NULL,
    client_name TEXT,
    status TEXT NOT NULL DEFAULT 'Draft' CHECK (status IN ('Draft', 'In Review', 'Active', 'Completed')),
    total REAL NOT NULL DEFAULT 0,
    description TEXT,
    -- Global markup defaults
    materials_markup REAL NOT NULL DEFAULT 20,
    labor_markup REAL NOT NULL DEFAULT 25,
    equipment_markup REAL NOT NULL DEFAULT 15,
    subs_markup REAL NOT NULL DEFAULT 10,
    other_markup REAL NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_projects_client_id ON projects(client_id);
CREATE INDEX idx_projects_status ON projects(status);
```

### Migration 004: Materials

> **Ordering rationale:** Materials must be created before the Estimate Builder tables because `line_items.material_id` references `materials(id)`. This was originally numbered 005 but swapped to avoid a forward FK reference.

```sql
-- +goose Up
CREATE TABLE suppliers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE materials (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    supplier_id TEXT NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    unit_price REAL NOT NULL DEFAULT 0,
    unit TEXT NOT NULL DEFAULT 'ea',
    supplier_code TEXT,
    price_source TEXT NOT NULL DEFAULT 'Manual' CHECK (price_source IN ('Supplier', 'Manual')),
    last_updated TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_materials_supplier_id ON materials(supplier_id);
```

### Migration 005: Estimate Builder (4-level hierarchy)

```sql
-- +goose Up
CREATE TABLE sections (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_sections_project_id ON sections(project_id);

CREATE TABLE subcategories (
    id TEXT PRIMARY KEY,
    section_id TEXT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    lump_sum REAL NOT NULL DEFAULT 0,
    materials_markup REAL,
    labor_markup REAL,
    equipment_markup REAL,
    subs_markup REAL,
    other_markup REAL,
    materials_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    labor_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    equipment_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    subs_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    other_markup_enabled BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX idx_subcategories_section_id ON subcategories(section_id);

CREATE TABLE component_groups (
    id TEXT PRIMARY KEY,
    subcategory_id TEXT NOT NULL REFERENCES subcategories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_component_groups_subcategory_id ON component_groups(subcategory_id);

CREATE TABLE line_items (
    id TEXT PRIMARY KEY,
    subcategory_id TEXT NOT NULL REFERENCES subcategories(id) ON DELETE CASCADE,
    component_group_id TEXT REFERENCES component_groups(id) ON DELETE SET NULL,
    category_type TEXT NOT NULL CHECK (category_type IN ('materials', 'labor', 'equipment', 'subs', 'other')),
    item_name TEXT NOT NULL,
    quantity REAL NOT NULL DEFAULT 1,
    unit TEXT NOT NULL DEFAULT 'ea',
    unit_price REAL NOT NULL DEFAULT 0,
    is_custom BOOLEAN NOT NULL DEFAULT true,
    material_id TEXT REFERENCES materials(id) ON DELETE SET NULL,
    price_override BOOLEAN NOT NULL DEFAULT false,
    description TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_line_items_subcategory_id ON line_items(subcategory_id);
CREATE INDEX idx_line_items_component_group_id ON line_items(component_group_id);
```

### Migration 006: Rates

```sql
-- +goose Up
CREATE TABLE rate_categories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed default categories
INSERT INTO rate_categories (id, name, sort_order) VALUES
    ('cat_labor', 'Labor', 0),
    ('cat_equipment', 'Equipment Rentals', 1),
    ('cat_subs', 'Subcontractors', 2),
    ('cat_other', 'Other', 3);

CREATE TABLE rates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category_id TEXT NOT NULL REFERENCES rate_categories(id) ON DELETE CASCADE,
    supplier TEXT,
    rate REAL NOT NULL DEFAULT 0,
    unit TEXT NOT NULL DEFAULT 'hour',
    notes TEXT,
    last_updated TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_rates_category_id ON rates(category_id);
```

### Migration 007: Quotes and Signatures

```sql
-- +goose Up
CREATE TABLE quotes (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'signed', 'expired', 'superseded')),
    estimate_snapshot JSONB,
    totals_snapshot JSONB,
    token TEXT UNIQUE,
    expires_at TIMESTAMPTZ,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(project_id, version)
);

CREATE INDEX idx_quotes_project_id ON quotes(project_id);
CREATE INDEX idx_quotes_token ON quotes(token);

CREATE TABLE quote_signatures (
    id TEXT PRIMARY KEY,
    quote_id TEXT NOT NULL REFERENCES quotes(id) ON DELETE RESTRICT,
    signer_name TEXT NOT NULL,
    signer_ip TEXT,
    signed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_quote_signatures_quote_id ON quote_signatures(quote_id);

CREATE TABLE quote_emails (
    id TEXT PRIMARY KEY,
    quote_id TEXT NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    recipient TEXT NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    provider_id TEXT
);

CREATE INDEX idx_quote_emails_quote_id ON quote_emails(quote_id);
```

### Migration 008: Company Profile

```sql
-- +goose Up
CREATE TABLE company_profile (
    id TEXT PRIMARY KEY DEFAULT 'default',
    name TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    address TEXT,
    city TEXT,
    state TEXT,
    zip TEXT,
    logo_path TEXT,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## 4. Build Order

Each phase is independently testable and represents a logical commit boundary. Phases should be completed sequentially — each depends on prior phases.

---

### Phase 0: Project scaffold + schema + CLI seed tool

**Goal:** Fresh project compiles, migrations run, seed tool creates a user, server starts and serves health check.

#### Steps

1. **Delete all old migrations** from `cmd/server/migrations/`. Write migrations 001–008 as specified in §3.
   - Migration ordering is already resolved: materials is 004, estimate builder is 005 (see §3 rationale).

2. **Delete all old sqlc queries** from `sqlc/queries/`. Create fresh query files as needed per phase.

3. **Write initial sqlc queries** for this phase:
   - `sqlc/queries/users.sql`: `CreateUser`, `GetUser`, `GetUserByEmail`, `ListUsers`
   - `sqlc/queries/sessions.sql`: `CreateSession`, `GetSessionByTokenHash`, `DeleteSessionByTokenHash`, `DeleteExpiredSessions`, `DeleteUserSessions`, `UpdateSessionActivity`

4. **Update `sqlc.yaml`** if necessary for the new schema.

5. **Run `make sqlc`** to generate repository code.

6. **Create `cmd/seed/main.go`** — a CLI tool for user management:
   ```
   Usage:
     skalkaho-seed create --email <email> --password <password> --name <name>
     skalkaho-seed list
     skalkaho-seed delete --email <email>
   ```
   - Connects to PostgreSQL via `DATABASE_URL` env var
   - Runs migrations automatically before seeding
   - Uses `auth.HashPassword()` for password hashing
   - Generates nanoid for user ID

7. **Simplify `internal/config/config.go`:**
   - Remove: `DatabasePath`, `SeedDemoUser`, `AnthropicAPIKey`, `AutoApproveThreshold`, `PriceImportToken`
   - Keep: `Addr`, `DatabaseURL`, `Environment`, `SessionSecret`, `SessionDuration`, `SessionCookieName`, `SecureCookies`
   - Add: `PostmarkAPIKey` (for Phase 9)

8. **Simplify `internal/database/`:** PostgreSQL only. Remove SQLite detection.

9. **Simplify `internal/auth/`:**
   - `auth.go`: Keep `HashPassword`, `CheckPassword`, `GenerateSessionToken`, `HashSessionToken`. No changes needed.
   - `session.go`: Remove `OrgID` from `Session` struct and `CreateSessionParams`. Remove org validation from `CreateSession()`. Update `ValidateSession()` to not check org.
   - `context.go`: Remove `OrgIDFromContext`, `SetOrgID`. Keep `UserIDFromContext`, `SetUserID`, `SessionFromContext`, `SetSession`.
   - `middleware.go`: Keep `SessionMiddleware`, `RequireAuth`. Remove `RequireRole`.

10. **Simplify `internal/domain/`:**
    - Keep `estimate.go` (EstimatePayload, all Estimate* types, ResolveMarkup). No `org_id` fields exist on these types — no changes needed.
    - Delete `domain.go` (old Job, Category, LineItem, surcharge calculation).
    - Update `MaterialDBEntry` and `RateDBEntry` to match the new `materials`/`rates` table fields.

11. **Stub `cmd/server/main.go`:**
    - Remove demo seeding, SQLite support, old handler initialization
    - Wire: config → database → migrations → queries → session manager → router → middleware → serve
    - Start with just the health check route

12. **Verify:** `go build ./cmd/server && go build ./cmd/seed` compiles. Migrations run. Seed tool creates a user. Health check returns 200.

#### Acceptance criteria
- `cmd/seed create --email test@test.com --password test1234 --name "Test User"` creates a user row
- `cmd/seed list` shows the user
- `go run ./cmd/server` starts, runs migrations, `/health` returns 200
- No references to `org_id`, organizations, or SQLite anywhere in the codebase

---

### Phase 1: Auth (simplified)

**Goal:** Login page works. Session cookie set. Protected routes redirect to login.

#### Steps

1. **Write auth handler** (`internal/handler/auth/`):
   - `GetLogin` — render login form
   - `PostLogin` — validate credentials, create session, set cookie, redirect to `/`
   - `Logout` — destroy session, clear cookie, redirect to `/login`
   - No registration handler. No password reset (Phase 2 of spec).

2. **Write login template** (`internal/templates/pages/login.html`):
   - Email + password form
   - Error message display
   - Light Tailwind styling
   - No link to register

3. **Write sqlc queries** for auth:
   - `GetUserByEmail(email) → user`
   - `CreateSession(params) → session`
   - `GetSessionByTokenHash(hash) → session + user join`
   - `DeleteSessionByTokenHash(hash)`
   - `DeleteExpiredSessions()`
   - `UpdateSessionActivity(id)`

4. **Wire routes:**
   - `GET /login` → GetLogin (public)
   - `POST /login` → PostLogin (public)
   - `GET /logout` → Logout (public)
   - `POST /logout` → Logout (public)
   - `GET /` → redirect to login if not authenticated (placeholder)

5. **Wire `protect()` helper** in router — same pattern as existing: `SessionMiddleware → RequireAuth`.

#### Acceptance criteria
- Seed a user via CLI, log in via browser, see session cookie set
- Accessing `/` without session redirects to `/login`
- Logout clears session

---

### Phase 2: Navigation shell

**Goal:** Sidebar layout with 7 page links. All pages render stub content. Active page highlighted.

#### Steps

1. **Create base layout** (`internal/templates/layouts/base.html`):
   - Fixed 220px sidebar on left
   - Main content area on right
   - HTMX, Alpine.js, Tailwind CDN loaded in head
   - Svelte estimate builder script loaded conditionally (only on estimate builder page)

2. **Create sidebar partial** (`internal/templates/partials/sidebar.html`):
   - Logo/app name at top
   - Navigation links:
     - Projects (icon + label) → `/`
     - Clients → `/clients`
     - Materials → `/materials`
     - Rates → `/rates`
   - "Current Project" section (shown when a project is active):
     - Project name
     - Overview → `/projects/{id}`
     - Estimate Builder → `/projects/{id}/estimate`
     - Back to Projects → `/`
   - Active page state via Alpine.js or template conditional
   - User name + logout at bottom

3. **Create template renderer** (`internal/templates/`):
   - Port existing renderer pattern from `internal/templates/keyboard/`
   - Parse layouts, pages, and partials
   - `Render(w, pageName, data)` method

4. **Create stub pages** (each renders the layout with a heading):
   - `pages/projects.html` — "Projects"
   - `pages/clients.html` — "Clients"
   - `pages/materials.html` — "Materials Database"
   - `pages/rates.html` — "Labor & Equipment Rates"
   - `pages/project_overview.html` — "Job Overview"
   - `pages/estimate_builder.html` — Svelte mount point
   - `pages/login.html` — already built in Phase 1

5. **Wire stub routes** in router:
   ```
   GET /                           → Projects list (stub)
   GET /clients                    → Clients list (stub)
   GET /materials                  → Materials DB (stub)
   GET /rates                      → Rates (stub)
   GET /projects/{id}              → Job Overview (stub)
   GET /projects/{id}/estimate     → Estimate Builder (stub)
   ```

#### Acceptance criteria
- All 6 sidebar links navigate to stub pages
- Sidebar highlights the active page
- Layout is consistent across all pages
- Logged-out users are redirected to login

---

### Phase 3: Projects module

**Goal:** Full CRUD for projects. Dashboard with status tabs, stats, search, sortable columns.

Implements spec §4.1.

#### sqlc queries (`sqlc/queries/projects.sql`)

```
CreateProject(params) → project
GetProject(id) → project
ListProjects() → []project  (with optional status filter, search, sort)
UpdateProject(params)
UpdateProjectStatus(id, status)
UpdateProjectTotal(id, total)
DeleteProject(id)
CountProjectsByStatus() → {draft, in_review, active, completed counts + totals}
```

#### Handler (`internal/handler/projects.go`)

| Method | Route | Notes |
|--------|-------|-------|
| `ListProjects` | `GET /` | Status tab filter, search, sort. Renders projects list page. |
| `GetProjectForm` | `GET /project-form` | Returns create project modal HTML (htmx partial). |
| `CreateProject` | `POST /projects` | Create project, optionally link client. Redirect or htmx swap. |
| `DeleteProject` | `DELETE /projects/{id}` | Confirm modal, then delete. |
| `UpdateProjectStatus` | `PATCH /projects/{id}/status` | Status badge click → modal → update. |

#### Template (`pages/projects.html`)

- Stats row: 4 KPIs (Draft count, In Review $, Active $, Completed $)
- Status tabs: All / Draft / In Review / Active / Completed
- Search input (htmx: `hx-get="/?search=..."` with debounce)
- Sortable table: Name, Client, Status, Total, Created, Updated
- Each row links to `/projects/{id}` (Job Overview)
- "New Project" button → modal with name + client selector + description

#### Template (`partials/project_form.html`)

- Modal with: project name (required), client selector (dropdown from clients table), description (optional)

#### Acceptance criteria
- Create, view, delete projects
- Status tabs filter the list
- Search filters by name + client name
- Stats row shows correct counts/totals
- Clicking a project navigates to `/projects/{id}`

---

### Phase 4: Clients module

**Goal:** Full CRUD for clients. List with search, sortable columns, expandable row detail.

Implements spec §4.2.

#### sqlc queries (`sqlc/queries/clients.sql`)

```
CreateClient(params) → client
GetClient(id) → client
ListClients() → []client  (with optional search, sort)
UpdateClient(params)
DeleteClient(id)
CountClients() → count
CountProjectsByClient(client_id) → count
```

#### Handler (`internal/handler/clients.go`)

| Method | Route | Notes |
|--------|-------|-------|
| `ListClients` | `GET /clients` | Search, sort. Stats row. |
| `GetClientForm` | `GET /client-form` | Create modal. |
| `CreateClient` | `POST /clients` | |
| `GetClientEditForm` | `GET /clients/{id}/edit` | Edit modal (pre-filled). |
| `UpdateClient` | `PUT /clients/{id}` | |
| `DeleteClient` | `DELETE /clients/{id}` | Confirm modal. |

#### Template (`pages/clients.html`)

- Stats row: Total clients, Total projects (across all clients), Avg projects/client
- Search input
- Sortable table: Company, Contact, Email, Phone, Projects, Created
- Expandable row: click to show address + notes inline (Alpine.js `x-data` toggle)
- Add/Edit modals

#### Acceptance criteria
- Full CRUD
- Search filters all fields
- Row expands to show address + notes
- Stats row accurate

---

### Phase 5: Materials Database

**Goal:** Materials organized by supplier tabs. Full CRUD. Search, filter, sort.

Implements spec §4.3. This is a **new entity** — no equivalent in the old codebase.

#### sqlc queries (`sqlc/queries/materials.sql`, `sqlc/queries/suppliers.sql`)

```
-- Suppliers
CreateSupplier(params) → supplier
ListSuppliers() → []supplier
DeleteSupplier(id)  -- cascades to materials

-- Materials
CreateMaterial(params) → material
GetMaterial(id) → material
ListMaterials() → []material  (with optional supplier_id filter, search, sort, price_source filter)
ListMaterialsBySupplier(supplier_id) → []material
UpdateMaterial(params)
DeleteMaterial(id)
CountMaterialsBySupplier(supplier_id) → count
```

#### Handler (`internal/handler/materials.go`)

| Method | Route | Notes |
|--------|-------|-------|
| `ListMaterials` | `GET /materials` | Supplier tabs, search, price source filter, sort. |
| `GetMaterialForm` | `GET /material-form` | Create/edit modal. |
| `CreateMaterial` | `POST /materials` | |
| `GetMaterialEditForm` | `GET /materials/{id}/edit` | Pre-filled edit modal. |
| `UpdateMaterial` | `PUT /materials/{id}` | Sets `last_updated = now()`. |
| `DeleteMaterial` | `DELETE /materials/{id}` | Confirm modal. |
| `CreateSupplier` | `POST /suppliers` | Inline input in tab bar. |
| `DeleteSupplier` | `DELETE /suppliers/{id}` | Confirm modal showing material count. |

#### Template (`pages/materials.html`)

- Supplier tabs: dynamic, each tab shows count. "+" to add supplier. "×" to delete (with confirm showing item count).
- Price source filter: All / Supplier / Manual
- Search input
- Sortable table: Name, Supplier Code, Unit Price, Unit, Source, Last Updated
- Add/Edit modal: name, supplier (dropdown), unit, price, supplier code, price source

#### Acceptance criteria
- Supplier tabs filter materials
- Add/remove suppliers dynamically
- Full material CRUD
- Price source filter works
- `last_updated` auto-set on save

---

### Phase 6: Labor & Equipment Rates

**Goal:** Rates organized by category tabs. Full CRUD. Search, sort.

Implements spec §4.4. This is a **new entity** — no equivalent in the old codebase.

#### sqlc queries (`sqlc/queries/rates.sql`, `sqlc/queries/rate_categories.sql`)

```
-- Rate Categories
CreateRateCategory(params) → rate_category
ListRateCategories() → []rate_category
DeleteRateCategory(id)  -- cascades to rates

-- Rates
CreateRate(params) → rate
GetRate(id) → rate
ListRates() → []rate  (with optional category_id filter, search, sort)
ListRatesByCategory(category_id) → []rate
UpdateRate(params)
DeleteRate(id)
CountRatesByCategory(category_id) → count
```

#### Handler (`internal/handler/rates.go`)

| Method | Route | Notes |
|--------|-------|-------|
| `ListRates` | `GET /rates` | Category tabs, search, sort. |
| `GetRateForm` | `GET /rate-form` | Create/edit modal. |
| `CreateRate` | `POST /rates` | |
| `GetRateEditForm` | `GET /rates/{id}/edit` | Pre-filled. |
| `UpdateRate` | `PUT /rates/{id}` | Sets `last_updated = now()`. |
| `DeleteRate` | `DELETE /rates/{id}` | Confirm modal. |
| `CreateRateCategory` | `POST /rate-categories` | Inline input in tab bar. |
| `DeleteRateCategory` | `DELETE /rate-categories/{id}` | Confirm modal with count. |

#### Template (`pages/rates.html`)

- Category tabs: Labor, Equipment Rentals, Subcontractors, Other (seeded). Counts as badges. Add/remove tabs.
- Category count badges in header (spec §4.4: "labor / equipment / subs counts in header")
- Search input
- Sortable table: Name, Category, Supplier, Rate, Unit, Notes, Last Updated
- Add/Edit modal: name, category (dropdown), supplier (optional), rate, unit, notes

#### Acceptance criteria
- Category tabs filter rates
- Add/remove categories
- Full rate CRUD
- Count badges accurate

---

### Phase 7: Estimate Builder (port)

**Goal:** Existing Svelte estimate builder works with the new schema. Autocomplete pulls from `materials` + `rates` tables.

#### Steps

1. **Port `ui/` directory:**
   - Update `ui/src/main.js` if it references any old paths
   - Update API URLs in `EstimateBuilder.svelte`: `/api/estimate/{projectId}` (should already match)
   - No component logic changes needed — markup engine, undo, autosave, autocomplete all stay

2. **Update Go API handler** (`internal/handler/estimate_api.go`):
   - Port from `internal/handler/keyboard/estimate_api.go`
   - Remove all `orgID`/`GetOrgID()` references
   - Change route param from `jobID` to `projectID`
   - Change `GetJob` → `GetProject` queries
   - Change autocomplete data source: instead of `ListItemTemplates`, query `materials` + `rates` tables
   - Build `MaterialDBEntry` from `materials` table (add `supplier` field)
   - Build `RateDBEntry` from `rates` table (map `category_id` to category name)

3. **Update domain types** (`internal/domain/estimate.go`):
   - Update `MaterialDBEntry` to include supplier info:
     ```go
     type MaterialDBEntry struct {
         ID           string  `json:"id"`
         Name         string  `json:"name"`
         Supplier     string  `json:"supplier"`
         UnitPrice    float64 `json:"unit_price"`
         Unit         string  `json:"unit"`
         SupplierCode string  `json:"supplier_code,omitempty"`
     }
     ```
   - Update `RateDBEntry`:
     ```go
     type RateDBEntry struct {
         ID       string  `json:"id"`
         Name     string  `json:"name"`
         Category string  `json:"category"`
         Rate     float64 `json:"rate"`
         Unit     string  `json:"unit"`
     }
     ```
   - Update `EstimateLineItem.MaterialID` from `*int64` to `*string` (TEXT FK now)

4. **Update Svelte autocomplete** (`ui/src/lib/Autocomplete.svelte`):
   - Update field mappings if `MaterialDBEntry` or `RateDBEntry` field names changed
   - Materials: `unit_price` instead of `default_price`, `unit` instead of `default_unit`
   - Rates: `rate` instead of `default_price`, `unit` instead of `default_unit`

5. **Write sqlc queries** (`sqlc/queries/estimate_builder.sql`):
   ```
   -- Sections
   ListSectionsByProject(project_id) → []section
   CreateSection(params) → section
   DeleteSectionsByProject(project_id)

   -- Subcategories
   ListSubcategoriesByProject(project_id) → []subcategory (JOIN through sections)
   CreateSubcategory(params) → subcategory

   -- Component Groups
   ListComponentGroupsByProject(project_id) → []component_group (JOIN through sections → subcategories)
   CreateComponentGroup(params) → component_group

   -- Line Items
   ListLineItemsByProject(project_id) → []line_item (JOIN through sections → subcategories)
   CreateLineItem(params) → line_item

   -- Project markup globals
   GetProjectMarkupGlobals(project_id) → markup fields
   UpdateProjectMarkupGlobals(params)
   ```

6. **Write template** (`pages/estimate_builder.html`):
   - Sidebar layout with "Current Project" section active
   - `<div id="estimate-root" data-project-id="{{.Project.ID}}"></div>`
   - `<script type="module" src="/static/estimate-builder/estimate-builder.js"></script>`

7. **Build Svelte bundle:** `cd ui && npm run build` → outputs to `static/estimate-builder/`

8. **Create shared cost calculation function** (`internal/domain/cost.go`):
   - **This is a prerequisite for Phase 8, not an afterthought.** Build it now while the estimate save logic is fresh.
   - `CalculateProjectCosts(sections, subcategories, lineItems, globals) → ProjectCostSummary`
   - Inputs: the same data structures loaded by `buildEstimatePayload`
   - Uses `ResolveMarkup()` for each line item
   - Returns per-category totals, per-section breakdown, and grand total
   - The estimate builder's save handler (`SaveEstimateJSON`) should call this function to compute `projects.total`
   - Phase 8's overview handler will call the same function to render cost KPIs
   - **Single source of truth**: if markup logic changes, it changes in one place

#### Acceptance criteria
- Navigate to `/projects/{id}/estimate` → Svelte mounts, fetches data, renders hierarchy
- Add sections, subcategories, groups, line items
- Autocomplete shows materials and rates from new tables
- Auto-save persists to PostgreSQL
- Undo works (Ctrl+Z)
- Markup controls (global + per-subcategory) work
- Navigation guard warns on unsaved changes
- `projects.total` is updated on each save via `domain.CalculateProjectCosts()`

---

### Phase 8: Job Overview

**Goal:** Per-project summary page with client info, cost breakdown, status selector, and quote version list.

Implements spec §4.6.

#### Handler (`internal/handler/project_overview.go`)

| Method | Route | Notes |
|--------|-------|-------|
| `GetProjectOverview` | `GET /projects/{id}` | Full overview page. |
| `UpdateProjectStatus` | `PATCH /projects/{id}/status` | Status badge → modal → update. |
| `GetStatusModal` | `GET /projects/{id}/status-modal` | Returns status selector modal HTML. |

#### Template (`pages/project_overview.html`)

- **Client + project info card**: project name, client name/contact/email/phone, description, dates
- **Status badge**: clickable, opens modal with status options + descriptions
- **Cost summary KPIs**: Total + per-category cards. Only show categories with non-zero values.
  - Calculate from estimate builder data: query sections → subcategories → line items, apply markup, sum by category_type
- **Sections breakdown table**: per-section cost broken down by category (materials, labor, equipment, subs, other, total)
- **"Edit Estimate" button**: links to `/projects/{id}/estimate`
- **Quote version list** (stub for now — wired in Phase 9): shows v1, v2... with status badges

#### Cost calculation

The overview page uses the shared `domain.CalculateProjectCosts()` function created in Phase 7 step 8. **Do not reimplement cost calculation here.**

1. Load all sections, subcategories, line items for the project (same queries as estimate builder)
2. Call `domain.CalculateProjectCosts(sections, subcategories, lineItems, globals)` → `ProjectCostSummary`
3. Render the `ProjectCostSummary` into the template (per-category KPIs, per-section breakdown, grand total)

The `projects.total` column is already kept in sync by the estimate builder's save handler (Phase 7). The overview reads it directly for the dashboard but recalculates the full breakdown on demand for the detail view.

#### Acceptance criteria
- Overview shows project info, client info, cost breakdown
- Status can be changed via modal
- Cost summary matches what the estimate builder shows
- "Edit Estimate" navigates to builder
- Sidebar shows "Current Project" section with active project name

---

### Phase 9: Quote & Signature system

**Goal:** Generate quotes with tokens, send via link or Postmark email, public `/q/{token}` page, e-signature capture.

Implements spec §6.

#### sqlc queries (`sqlc/queries/quotes.sql`)

```
CreateQuote(params) → quote
GetQuote(id) → quote
GetQuoteByToken(token) → quote + project + client info
ListQuotesByProject(project_id) → []quote
UpdateQuoteStatus(id, status)
UpdateQuoteSent(id, sent_at, expires_at, status)  -- for resend
GetLatestQuoteVersion(project_id) → version number
SupersedeQuote(id)  -- set status = 'superseded'

CreateQuoteSignature(params) → quote_signature
GetQuoteSignature(quote_id) → quote_signature (nullable)

CreateQuoteEmail(params) → quote_email
ListQuoteEmails(quote_id) → []quote_email
```

#### Handler (split across two files)

**`internal/handler/quotes.go`** (protected routes):

| Method | Route | Notes |
|--------|-------|-------|
| `ListQuotes` | `GET /projects/{id}/quotes` | Returns quote version list partial (for overview page). |
| `CreateQuote` | `POST /projects/{id}/quotes` | Snapshot estimate + totals, create quote row (status: draft). |
| `SendQuote` | `POST /quotes/{id}/send` | Generate token (nanoid 10 chars), set sent_at + expires_at, transition to 'sent'. Show modal with link + email option. |
| `ResendQuote` | `POST /quotes/{id}/resend` | Reset sent_at + expires_at. Same token. |
| `GetSendModal` | `GET /quotes/{id}/send-modal` | Modal with copy-link button + email form. |

**`internal/handler/public_quote.go`** (public routes, no auth):

| Method | Route | Notes |
|--------|-------|-------|
| `GetQuotePage` | `GET /q/{token}` | Token lookup. Render based on state: not found (404), expired, already signed, or active. |
| `SubmitSignature` | `POST /q/{token}` | Validate token, write `quote_signatures` row, update project status to "In Review". |

#### Postmark email service (`internal/service/email/postmark.go`)

- Simple HTTP client, no SDK
- `SendQuoteEmail(to, quoterURL, projectName, contractorName) error`
- Uses `POSTMARK_API_KEY` from config
- Falls back gracefully if no API key configured (copy-link only)

#### Quote snapshot

When creating a quote from a finalized estimate:
1. Load the full estimate payload (same as `GetEstimateJSON`)
2. Calculate totals server-side (same as Job Overview calculation)
3. Serialize both as JSONB into `estimate_snapshot` and `totals_snapshot`
4. These are **immutable** once status transitions from draft to sent

#### Templates

**`pages/quote_public.html`** (no sidebar, standalone page):
- Contractor name (prominent, from `company_profile`)
- Project name, client name
- Section breakdown: section name + section total
- Grand total breakdown: Materials, Labor, Equipment, Subs, Other, Total
- Quote version + sent date
- Expiry notice if applicable
- **Active state**: signature form (typed name + "I Accept This Quote" button)
- **Signed state**: "Thank you" confirmation with signer name + date
- **Expired state**: "This quote has expired" message
- Print/save as PDF via native browser print (CSS `@media print`)

**`partials/quote_list.html`** (for Job Overview page):
- Version list with status badges
- "Send New Quote" button
- Per-quote: version number, status badge, sent date, actions (resend/view)

#### Acceptance criteria
- Create quote from overview page → snapshots estimate
- Send quote → generates token, shows copyable link
- `/q/{token}` shows quote content
- Client can type name + accept → signature recorded
- Project status transitions to "In Review" on signature
- Expired quotes show expired state
- New version supersedes previous
- Resend resets expiry without creating new version
- Postmark email sends if API key configured

---

## 5. File Structure (target)

```
cmd/
├── server/
│   ├── main.go                  # Entry point
│   └── migrations/              # Goose SQL migrations (001-008)
├── seed/
│   └── main.go                  # CLI user management tool
internal/
├── auth/
│   ├── auth.go                  # HashPassword, CheckPassword, token generation
│   ├── session.go               # SessionManager (simplified, no org)
│   ├── context.go               # UserID context helpers (no org)
│   └── middleware.go            # SessionMiddleware, RequireAuth (no RequireRole)
├── config/
│   └── config.go                # Environment configuration (simplified)
├── database/
│   └── database.go              # PostgreSQL connection (no SQLite)
├── domain/
│   └── estimate.go              # EstimatePayload types, ResolveMarkup, cost calculation
├── handler/
│   ├── auth/
│   │   └── handler.go           # Login, Logout
│   ├── clients.go               # Client CRUD
│   ├── materials.go             # Materials + Suppliers CRUD
│   ├── rates.go                 # Rates + Rate Categories CRUD
│   ├── projects.go              # Projects CRUD, dashboard
│   ├── project_overview.go      # Job Overview page
│   ├── estimate_api.go          # Svelte JSON API (GET/POST estimate)
│   ├── quotes.go                # Quote management (protected)
│   ├── public_quote.go          # /q/{token} (public)
│   └── handler.go               # Shared handler struct, helpers
├── middleware/
│   ├── recover.go
│   ├── request_id.go
│   └── logger.go
├── repository/                  # sqlc-generated code
├── router/
│   └── router.go                # All route definitions
├── service/
│   └── email/
│       └── postmark.go          # Postmark API client
└── templates/
    ├── layouts/
    │   └── base.html            # Sidebar + main content shell
    ├── pages/
    │   ├── login.html
    │   ├── projects.html        # Dashboard
    │   ├── clients.html
    │   ├── materials.html
    │   ├── rates.html
    │   ├── project_overview.html
    │   ├── estimate_builder.html
    │   └── quote_public.html    # Public quote page (no sidebar)
    ├── partials/
    │   ├── sidebar.html
    │   ├── project_form.html
    │   ├── client_form.html
    │   ├── material_form.html
    │   ├── rate_form.html
    │   ├── status_modal.html
    │   ├── quote_list.html
    │   ├── send_quote_modal.html
    │   └── delete_confirm.html  # Reusable confirm modal
    └── renderer.go              # Template rendering logic
sqlc/
├── sqlc.yaml
└── queries/
    ├── users.sql
    ├── sessions.sql
    ├── projects.sql
    ├── clients.sql
    ├── materials.sql
    ├── suppliers.sql
    ├── rates.sql
    ├── rate_categories.sql
    ├── estimate_builder.sql
    └── quotes.sql
ui/                              # Svelte 5 Estimate Builder (existing, ported)
├── package.json
├── vite.config.ts
└── src/
    ├── main.js
    ├── EstimateBuilder.svelte
    └── lib/
        ├── SectionBlock.svelte
        ├── SubcategoryBlock.svelte
        ├── ComponentGroupBlock.svelte
        ├── LineItemRow.svelte
        ├── Autocomplete.svelte
        ├── FooterSummary.svelte
        ├── SaveStatus.svelte
        ├── autosave.svelte.js
        ├── undo.svelte.js
        └── markup.js
static/
└── estimate-builder/            # Vite build output
    └── estimate-builder.js
```

---

## 6. Conventions

### Go patterns

- **HTTP routing** (Go 1.22+): `mux.HandleFunc("GET /projects/{id}", handler.GetProject)`
- **Error handling**: Wrap with context, early return: `return fmt.Errorf("creating project: %w", err)`
- **Handler struct**: Single `Handler` struct in `handler.go` with `db`, `queries`, `renderer`, `logger` fields. All handlers are methods on this struct.
- **Auth handler**: Separate struct in `handler/auth/` with its own dependencies.
- **Database access**: Always via sqlc-generated `*repository.Queries`. No raw SQL in handlers.
- **Transactions**: `h.db.BeginTx()` → `h.queries.WithTx(tx)` → operations → `tx.Commit()`
- **IDs**: nanoid (TEXT) for all entities. Use a small Go helper or inline `nanoid.New()`.
- **No org_id**: Never scope queries by org_id. No tenant isolation. Single-tenant.

### Frontend patterns

- **HTMX**: `hx-get`, `hx-post`, `hx-swap`, `hx-target` for all server interactions except estimate builder
- **Alpine.js**: `x-data` for local UI state (modals, tabs, toggles, search). No heavy state management.
- **Tailwind CDN**: No build step for CSS. Use utility classes directly.
- **Modals**: Backdrop blur, rendered as htmx partials swapped into a modal container
- **Search**: `hx-get` with `hx-trigger="input changed delay:300ms"` for debounced search
- **Sort**: Query params `?sort=name&dir=asc`, htmx replaces table body
- **Tabs**: Alpine.js `x-data="{ tab: 'all' }"` for client-side tab switching, or htmx for server-filtered tabs

### Template patterns

- Base layout defines `{{block "content" .}}{{end}}`
- Pages define `{{define "content"}}...{{end}}`
- Partials are `{{template "partial_name" .}}`
- Data passed as `map[string]interface{}` from handlers
- Active page passed as `"ActivePage": "projects"` for sidebar highlighting

---

*End of Implementation Plan*
