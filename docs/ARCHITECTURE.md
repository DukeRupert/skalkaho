# Skalkaho Architecture

Construction quoting application for small-to-medium contractors. Single-tenant: each deployment serves one contractor.

**Stack**: Go 1.22+ backend, PostgreSQL, HTMX + Alpine.js + Tailwind CDN pages, Svelte 5 island (estimate builder only).

---

## Core Entities

```
User
├── Sessions (1:N)
└── Quotes created (1:N)

Client ←── Project (N:1 optional)
              ├── Sections (1:N CASCADE)
              │   └── Subcategories (1:N CASCADE)
              │       ├── ComponentGroups (1:N CASCADE, optional)
              │       │   └── LineItems (1:N)
              │       └── LineItems (1:N, ungrouped)
              │           ├── Material (N:1 optional)
              │           └── Subcontractor (N:1 optional)
              └── Quotes (1:N CASCADE)
                  ├── QuoteSignatures (1:N RESTRICT)
                  └── QuoteEmails (1:N CASCADE)

Template
├── TemplateSections (1:N CASCADE)
│   └── TemplateSubcategories (1:N CASCADE)
│       └── TemplateComponentGroups (1:N CASCADE)
└── stamped into → Project (at creation time, one-way copy)

Supplier ──→ Materials (1:N CASCADE)
RateCategory ──→ Rates (1:N CASCADE)
Subcontractor ←──→ Trades (N:N via subcontractor_trades, position 0 = primary)
CompanyProfile (singleton, id='default')
```

---

## Database Schema

### Users & Sessions (Migration 001)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `users` | id, email (UNIQUE), password_hash, name, status (active\|inactive) | No registration endpoint; created via `mage seed:create` |
| `sessions` | id, user_id (FK CASCADE), token_hash (UNIQUE), expires_at, last_activity_at | SHA-256 of plaintext token |

### Clients (Migration 002)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `clients` | id, company_name, contact_name, email, phone, address, notes | Linked to projects via projects.client_id |

### Projects (Migration 003)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `projects` | id, name, client_id (FK SET NULL), status, total | Statuses: Draft, In Review, Active, Completed |
| | materials_markup, labor_markup, equipment_markup, subs_markup, other_markup | Global markup defaults (20%, 25%, 15%, 10%, 10%) |

### Materials & Suppliers (Migration 004)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `suppliers` | id, name (UNIQUE), sort_order | Grouping for materials |
| `materials` | id, name, supplier_id (FK CASCADE), unit_price, unit, supplier_code, price_source | price_source: Supplier or Manual |

### Estimate Structure (Migration 005)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `sections` | id, project_id (FK CASCADE), name, sort_order | Top level |
| `subcategories` | id, section_id (FK CASCADE), name, lump_sum, {type}_markup (nullable override), {type}_markup_enabled | Markup controls per category |
| `component_groups` | id, subcategory_id (FK CASCADE), name, sort_order | Optional grouping |
| `line_items` | id, subcategory_id (FK CASCADE), component_group_id (FK SET NULL), category_type, item_name, quantity, unit, unit_price, material_id (FK SET NULL), subcontractor_id (FK SET NULL), description | category_type: materials\|labor\|equipment\|subs\|other |

### Rates (Migration 006)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `rate_categories` | id, name (UNIQUE) | Seeded: Labor, Equipment Rentals, Subcontractors, Other |
| `rates` | id, name, category_id (FK CASCADE), rate, unit | Used in estimate autocomplete |

### Quotes (Migration 007)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `quotes` | id, project_id (FK CASCADE), version, status, estimate_snapshot (JSONB), totals_snapshot (JSONB), token (UNIQUE), expires_at | Statuses: draft, sent, signed, expired, superseded |
| `quote_signatures` | id, quote_id (FK RESTRICT), signer_name, signer_ip, signed_at | RESTRICT prevents deleting signed quotes |
| `quote_emails` | id, quote_id (FK CASCADE), recipient, provider_id | Postmark message tracking |

### Company Profile (Migration 008)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `company_profile` | id ('default'), name, email, phone, address, city, state, zip, logo_path | Singleton row |

### Subcontractors & Trades (Migration 009-010)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `trades` | id, name (UNIQUE), sort_order | 16 seeded construction trades |
| `subcontractors` | id, name, company, phone, email, is_favorite | Directory of sub companies |
| `subcontractor_trades` | subcontractor_id, trade_id, position | Composite PK; position 0 = primary trade |

### Templates (Migration 012)
| Table | Key Columns | Notes |
|-------|------------|-------|
| `templates` | id, name, description, created_at, updated_at | Reusable project skeletons |
| `template_sections` | id, template_id (FK CASCADE), name, sort_order | Mirrors sections table |
| `template_subcategories` | id, template_section_id (FK CASCADE), name, sort_order | Mirrors subcategories table |
| `template_component_groups` | id, template_subcategory_id (FK CASCADE), name, sort_order | Mirrors component_groups table |

Templates are structure-only (no line items, no markup overrides). When a project is created with a template selected, the entire tree is copied into the project's own sections/subcategories/component_groups within a single transaction. Once stamped, the project is fully independent — editing or deleting the template has no effect on existing projects.

All IDs are TEXT (uuid[:20] or semantic slugs). All timestamps are TIMESTAMPTZ with `DEFAULT now()`.

---

## Request Flow

```
HTTP Request
  │
  ├─ Recover middleware (panic → 500)
  ├─ RequestID middleware (UUID in context + header)
  ├─ Logger middleware (request-scoped slog, logs duration)
  │
  ├─ Public routes: /health, /login, /logout, /static/*, /q/{token}
  │
  └─ Protected routes: protect(SessionMiddleware + RequireAuth)
       │
       ├─ HTML pages: handler → renderer.Render(w, "page.html", data)
       ├─ HTMX partials: handler → renderer.RenderPartial(w, "page.html", "block", data)
       └─ JSON API: handler → json.Encode(w, payload)
```

**HTMX detection**: `HX-Request: true` && `HX-Boosted` not set → partial response. Boosted navigation gets full page.

---

## HTTP Routes

### Public
| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/health` | inline | Health check |
| GET | `/static/*` | FileServer | Static assets |
| GET | `/login` | auth.GetLogin | Login page |
| POST | `/login` | auth.PostLogin | Authenticate |
| GET/POST | `/logout` | auth.Logout | Destroy session |
| GET | `/q/{token}` | app.GetQuotePage | Public quote view |
| POST | `/q/{token}` | app.SubmitSignature | Sign quote |

### Protected (all wrapped with `protect()`)
| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/` | ListProjects | Project dashboard |
| GET | `/projects/new-modal` | NewProjectModal | Create form partial |
| POST | `/projects` | CreateProject | Create project |
| DELETE | `/projects/{id}` | DeleteProject | Delete project |
| PATCH | `/projects/{id}/status` | UpdateProjectStatus | Change status |
| PATCH | `/projects/{id}/client` | UpdateProjectClient | Assign client |
| GET | `/projects/{id}` | GetProjectOverview | Project overview |
| GET | `/projects/{id}/estimate` | EstimateBuilder | Svelte island page |
| GET | `/api/estimate/{id}` | GetEstimate | JSON: full payload |
| POST | `/api/estimate/{id}` | SaveEstimate | JSON: persist estimate |
| GET/POST | `/projects/{id}/quotes` | List/CreateQuote | Quote management |
| POST | `/quotes/{id}/send` | SendQuote | Generate token, email |
| POST | `/quotes/{id}/resend` | ResendQuote | Re-send email |
| GET/POST/DELETE | `/clients/*` | Client CRUD | Client directory |
| GET/POST/DELETE | `/materials/*` | Material CRUD | Materials database |
| POST/DELETE | `/suppliers/*` | Supplier CRUD | Supplier management |
| GET/POST/DELETE | `/rates/*` | Rate CRUD | Rate management |
| POST/DELETE | `/rate-categories/*` | Category CRUD | Rate categories |
| GET/POST/DELETE/PATCH | `/subcontractors/*` | Sub CRUD + favorite | Subcontractor directory |
| GET/POST/DELETE | `/templates/*` | Template CRUD | Template management |
| POST/DELETE | `/templates/{id}/sections/*` | Section CRUD | Template tree nodes |
| POST/DELETE | `/templates/{id}/subcategories/*` | Subcategory CRUD | Template tree nodes |
| POST/DELETE | `/templates/{id}/groups/*` | Group CRUD | Template tree nodes |

---

## Estimate Builder Flow

```
Browser                          Server                         Database
  │                                │                               │
  │  GET /projects/{id}/estimate   │                               │
  │──────────────────────────────→│  Render estimate_builder.html  │
  │←──────────────────────────────│                               │
  │                                │                               │
  │  Svelte mounts, calls          │                               │
  │  GET /api/estimate/{id}        │                               │
  │──────────────────────────────→│  buildEstimatePayload()        │
  │                                │──────────────────────────────→│ Load sections/subcats/
  │                                │←──────────────────────────────│ items/materials/rates/subs
  │←──────────────────────────────│  EstimatePayload JSON         │
  │                                │                               │
  │  User edits line items         │                               │
  │  (Svelte reactive state)       │                               │
  │                                │                               │
  │  Auto-save (2s debounce)       │                               │
  │  POST /api/estimate/{id}       │                               │
  │──────────────────────────────→│  Delete all sections (CASCADE) │
  │                                │──────────────────────────────→│ Recreate from payload
  │                                │  CalculateProjectTotal()      │
  │                                │──────────────────────────────→│ Update project.total
  │←──────────────────────────────│  Return saved state           │
```

**Auto-save states**: clean → dirty → saving → clean (or error with retry)

**Undo**: Snapshots `{globals, sections}` before structural mutations. 20-entry LIFO stack. Ctrl+Z restores.

---

## Markup & Cost Calculation

### 5 Category Types
| Type | Default Markup | Rate Category |
|------|---------------|---------------|
| materials | 20% | N/A (uses materials DB) |
| labor | 25% | Labor |
| equipment | 15% | Equipment Rentals |
| subs | 10% | Subcontractors |
| other | 10% | Other |

### Markup Resolution (per line item)
```
1. If subcategory.markup_enabled[type] = false → 0%
2. If subcategory.markup_overrides[type] != null → override value
3. Else → project global default
```

### Cost Calculation
```
Per line item:
  baseCost   = quantity × unit_price
  markup%    = ResolveMarkup(categoryType, globals, overrides, enabled)
  markupAmt  = baseCost × (markup% / 100)

Per subcategory:
  subtotal   = Σ(baseCost + markupAmt) for all line items + lump_sum

Project total:
  GrandTotal = BaseTotal + MarkupTotal + LumpSumTotal
```

Both Go (`domain/cost.go`) and JavaScript (`lib/markup.js`) implement identical logic.

---

## Quote Lifecycle

```
Estimate (live, editable)
  │
  ├─ Create Quote → snapshot estimate + totals as JSONB, status='draft'
  │                  supersede any active quotes for this project
  │
  ├─ Send Quote → generate 10-char token, set expires_at (+30 days)
  │               status='sent', optionally email via Postmark
  │
  ├─ Client views /q/{token} → read-only quote with cost breakdown
  │
  ├─ Client signs → record QuoteSignature, status='signed'
  │                  project.status → 'In Review'
  │
  └─ Expiry → past expires_at or superseded by newer version
```

Quotes are immutable snapshots. Creating a new quote version supersedes previous ones.

---

## Authentication

- **Passwords**: Argon2id (time=1, memory=64MB, threads=4, keyLen=32, salt=16B)
- **Sessions**: 32 random bytes → base64 token, SHA-256 stored in DB
- **Cookie**: HttpOnly, SameSite=Lax, Secure in production, configurable name/duration
- **No registration**: Users created via `mage seed:create` CLI only
- **No roles/permissions**: Single-tenant, all authenticated users have full access
- **Context**: `UserIDFromContext()`, `UserEmailFromContext()`, `UserNameFromContext()`, `IsAuthenticated()`

---

## Frontend Architecture

### HTMX Pages (all non-estimate pages)
- Go templates: `layouts/base.html` + `partials/sidebar.html` + `pages/*.html`
- HTMX for search, filtering, modals, inline CRUD (partial swaps)
- Alpine.js for client-side toggles, dropdowns, form state
- Tailwind CSS via CDN (no build step)
- `hx-boost` for SPA-like navigation between pages

### Svelte 5 Island (estimate builder only)
- Built with Vite, output to `static/estimate-builder/`
- Mounts on `#estimate-root` div with `data-project-id`
- Components: EstimateBuilder → SectionBlock → SubcategoryBlock → ComponentGroupBlock → LineItemRow
- Autocomplete picker for materials, rates, and subcontractors
- Auto-save + undo managed in Svelte reactive state

### Design System
- Dark mode throughout
- CSS custom properties: `--color-granite`, `--color-ink`, `--color-sunburst`, `--color-sage`, `--color-concrete`, etc.
- Fonts: `--font-ui` (headings/labels), `--font-body` (content)
- Responsive sidebar: 220px desktop, 56px tablet (icons only), off-screen mobile

---

## Package Structure

```
cmd/
├── server/main.go              Entry point, dependency wiring, embedded migrations
└── seed/main.go                CLI user management tool

internal/
├── auth/                       Password hashing, sessions, context helpers, middleware
├── config/                     Environment variable loading
├── database/                   PostgreSQL connection (pgx)
├── domain/                     EstimatePayload types, ResolveMarkup, cost calculation
├── handler/
│   ├── app/                    All application handlers (single Handler struct)
│   │   ├── handler.go          Handler struct, PageData, helpers
│   │   ├── clients.go          Client CRUD
│   │   ├── materials.go        Material + supplier CRUD
│   │   ├── rates.go            Rate + category CRUD
│   │   ├── projects.go         Project CRUD + overview
│   │   ├── subcontractors.go   Subcontractor CRUD + favorites
│   │   ├── templates.go        Template CRUD, tree editor, stamp function
│   │   ├── estimate_api.go     JSON API for estimate load/save
│   │   ├── quotes.go           Quote create/send/resend
│   │   └── public_quote.go     Public quote page + signature
│   └── auth/                   Login/logout handlers
├── middleware/                  Recover, RequestID, Logger
├── repository/                 sqlc-generated code (models + queries)
├── router/                     Route registration + protect() helper
├── service/email/              Postmark integration
└── templates/                  Embedded HTML templates
    ├── layouts/base.html
    ├── partials/sidebar.html
    └── pages/*.html

sqlc/queries/                   SQL query definitions for sqlc
migrations/                     Goose SQL migrations (read by sqlc)
ui/                             Svelte 5 estimate builder source
static/                         Built assets + logos + favicons
```

---

## Build & Development

| Command | Purpose |
|---------|---------|
| `mage dev` | Run development server |
| `mage build` | Compile bin/server + bin/seed |
| `mage test` | Run domain tests |
| `mage sqlc` | Generate repository code from SQL queries |
| `mage ui:build` | Build Svelte bundle to static/ |
| `mage ui:watch` | Watch mode for Svelte development |
| `mage deps` | Download and tidy Go modules |
| `mage clean` | Remove built binaries |
| `EMAIL=x PASSWORD=x NAME=x mage seed:create` | Create user |
| `mage seed:list` | List all users |

---

## Environment Variables

| Variable | Default | Required | Notes |
|----------|---------|----------|-------|
| `DATABASE_URL` | — | Yes | PostgreSQL connection string |
| `ADDR` | `:8080` | No | Server listen address |
| `ENVIRONMENT` | `development` | No | `production` enables secure cookies |
| `SESSION_SECRET` | dev fallback | Prod only | Required when ENVIRONMENT=production |
| `SESSION_DURATION` | `720h` | No | Session lifetime (30 days) |
| `SESSION_COOKIE_NAME` | `skalkaho_session` | No | Cookie name |
| `POSTMARK_API_KEY` | — | No | Enables quote email delivery |
