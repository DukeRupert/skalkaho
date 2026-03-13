# Skalkaho — Feature & Design Specification

> Derived from ProBuilder prototype · March 2026  
> Firefly Software · Confidential

---

## Table of Contents

1. [Application Overview](#1-application-overview)
2. [Data Model](#2-data-model)
3. [Markup / Surcharge Engine](#3-markup--surcharge-engine)
4. [Feature Catalog](#4-feature-catalog)
5. [Key Interaction Patterns](#5-key-interaction-patterns)
6. [Quote & Client Signature](#6-quote--client-signature)
7. [Go + htmx Tech Stack Evaluation](#7-go--htmx-tech-stack-evaluation)
8. [Scope Summary & Phase Recommendation](#8-scope-summary--phase-recommendation)
9. [Estimate Builder — Component Architecture](#9-estimate-builder--component-architecture)

---

## 1. Application Overview

ProBuilder is a single-page construction estimating application. The client runs a small framing/general contracting operation in the Bitterroot Valley and uses the tool to build job quotes from a private database of materials, labor rates, and equipment rentals. The application lives entirely in the browser in its prototype form — no backend, no persistence beyond page load.

### Navigation Model

A fixed 220px sidebar provides global navigation and context-aware project sub-navigation. The main content area renders one page at a time. There are seven distinct pages:

| Page | Description |
|---|---|
| Projects Dashboard | Project list, stats, new project modal |
| Clients | CRM-lite contact database |
| Materials Database | Supplier-organized price list |
| Labor & Equipment Rates | Categorized rate cards |
| Job Overview | Per-project summary, cost breakdown, document actions |
| Estimate Builder | The core quoting interface |
| Quote Page (public) | Client-facing quote view + e-signature |

The sidebar shows a "Current Project" section when a project is active, with sub-links to Overview and Estimate Builder, and a back link to the projects list. Recent projects are listed below main nav when no project is active.

### Visual Design Language

Dark theme throughout. GitHub-style color tokens: near-black backgrounds (`#0e1117`, `#161b22`, `#1c2129`), accent blue (`#4a7cff`), semantic green/amber/red/purple for status. Monospace font for all numeric values. Sticky topbars on every page. Modal dialogs with backdrop blur for create/edit/delete flows.

---

## 2. Data Model

The prototype defines five primary entities plus three new tables for the quote and signature system. All prototype entities are held in React state — persistence is a backend concern.

### 2.1 Project

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| name | string | job name, e.g. "Missoula Pole Barn" |
| client_id | FK → clients | |
| client_name | string | denormalized for display |
| status | enum | Draft \| In Review \| Active \| Completed |
| total | decimal | computed from estimate, stored on save |
| description | text | optional free text |
| created_at | timestamptz | |
| updated_at | timestamptz | |

### 2.2 Client

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| company_name | string | required, primary display name |
| contact_name | string | individual contact person |
| email | string | |
| phone | string | |
| address | string | single free-text field |
| notes | text | contractor notes about the client |
| created_at | timestamptz | |

### 2.3 Material

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| name | string | e.g. "2x4x10" |
| supplier | string | references a named supplier |
| unit_price | decimal | cost per unit before markup |
| unit | string | ea, sheet, roll, sqft, etc. |
| supplier_code | string | vendor SKU, optional |
| price_source | enum | Supplier \| Manual |
| last_updated | date | for staleness tracking |

### 2.4 Rate

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| name | string | e.g. "Logan per day $20", "Skytrack Week" |
| category | string | Labor \| Equipment Rentals \| Subcontractors \| Other (user-configurable) |
| supplier | string | vendor name, optional |
| rate | decimal | dollar amount per unit |
| unit | string | hour, day, week, sqft, ea, etc. |
| notes | text | notes on rate derivation |
| last_updated | date | for staleness tracking |

### 2.5 Estimate Structure

The estimate is a 4-level hierarchy stored as separate flat tables joined by foreign keys.

**sections**

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| project_id | FK → projects | |
| name | string | e.g. "Framing", "Roofing" |
| sort_order | integer | |

**subcategories**

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| section_id | FK → sections | |
| name | string | |
| sort_order | integer | |
| lump_sum | decimal | post-markup flat addition, default 0 |
| materials_markup | decimal | null = inherit global |
| labor_markup | decimal | null = inherit global |
| equipment_markup | decimal | null = inherit global |
| subs_markup | decimal | null = inherit global |
| other_markup | decimal | null = inherit global |
| materials_markup_enabled | boolean | default true |
| labor_markup_enabled | boolean | default true |
| equipment_markup_enabled | boolean | default true |
| subs_markup_enabled | boolean | default true |
| other_markup_enabled | boolean | default true |

**component_groups**

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| subcategory_id | FK → subcategories | |
| name | string | e.g. "Wall Studs", "Roof Purlins" |
| sort_order | integer | |

**line_items**

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| subcategory_id | FK → subcategories | |
| component_group_id | FK → component_groups | nullable |
| category_type | enum | Materials \| Labor \| Equipment \| Subs \| Other |
| item_name | string | |
| quantity | decimal | |
| unit | string | |
| unit_price | decimal | |
| is_custom | boolean | false if pulled from materials/rates DB |
| material_id | FK → materials | nullable, set when pulled from DB |
| price_override | boolean | true if user manually typed a price |
| description | text | optional note |
| sort_order | integer | |

### 2.6 Quote & Signature Tables

**quotes**

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| project_id | FK → projects | |
| version | integer | 1, 2, 3... scoped per project |
| status | enum | draft \| sent \| signed \| expired \| superseded |
| estimate_snapshot | jsonb | full estimate at time of send — immutable after send |
| totals_snapshot | jsonb | `{materials, labor, equipment, subs, other, total}` |
| token | text | unique, random, URL-safe (nanoid 10 chars) |
| expires_at | timestamptz | sent_at + 3 days |
| sent_at | timestamptz | null until first send |
| created_at | timestamptz | |
| created_by | FK → users | |

**quote_signatures**

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| quote_id | FK → quotes | |
| signer_name | text | typed name from client |
| signer_ip | text | from request headers |
| signed_at | timestamptz | |

**quote_emails**

| Field | Type | Notes |
|---|---|---|
| id | string (nanoid) | primary key |
| quote_id | FK → quotes | |
| recipient | text | email address sent to |
| sent_at | timestamptz | |
| provider_id | text | Postmark message ID for delivery tracking |

---

## 3. Markup / Surcharge Engine

The markup engine is the core business logic of the application. It must be understood precisely before any backend implementation.

### 3.1 Global Defaults

Five global markup percentages are configured per estimate:

| Type | Default |
|---|---|
| Materials | 20% |
| Labor | 25% |
| Equipment | 15% |
| Subs | 10% |
| Other | 10% |

These are the fallback for every line item unless overridden at the subcategory level.

### 3.2 Subcategory Overrides

Each subcategory exposes a markup row showing all five types. Any field can be overridden independently. Setting the global default resets all subcategory overrides for that type to null (reverts to inheriting global).

### 3.3 Per-Type Enable Toggles

Each markup type at the subcategory level has an on/off toggle. When off, markup for that type is 0% within that subcategory regardless of global or override value. Use case: fixed-price sub quote where you don't want to mark up the sub labor.

### 3.4 Markup Inheritance Resolution

Resolution order for a given line item:

1. Find the line item's subcategory
2. Check if the relevant type's `markup_enabled` toggle is false → use 0%
3. Check if the subcategory has a non-null override for that type → use it
4. Fall back to the global default for that type

> **Key insight:** markup is applied per line item based on its `category_type`, not the section it belongs to. A Materials line item in a section named "Labor" uses the materials markup rate.

### 3.5 Lump Sum Adjustment

Each subcategory has a `lump_sum` field (decimal, default 0). This amount is added directly to the subcategory total after markup calculations. Use case: fixed subcontractor quote, or an allowance for items not yet individually priced.

### 3.6 Calculation Formulas

| Value | Formula |
|---|---|
| Line item base | `quantity × unit_price` |
| Effective markup % | per resolution order in 3.4 |
| After-markup unit price | `unit_price × (1 + markup / 100)` |
| Line item total | `quantity × after_markup_unit_price` |
| Subcategory total | `Σ line item totals + lump_sum` |
| Section total | `Σ subcategory totals within section` |
| Project total | `Σ all section totals` |

---

## 4. Feature Catalog

**Status legend:**  
`✅ In Prototype` — fully implemented in JSX  
`⚠️ Aspirational` — UI present, not wired  
`○ Not Shown` — implied but absent from prototype  
`🆕 New` — not in prototype, added to spec  

### 4.1 Projects Module

| Feature | Status | Notes |
|---|---|---|
| Project list with search | ✅ In Prototype | client-side filter by name + client name |
| Status tab filter | ✅ In Prototype | All / Draft / In Review / Active / Completed |
| Sortable columns | ✅ In Prototype | all 6 columns, asc/desc toggle |
| Stats row (4 KPIs) | ✅ In Prototype | Draft count, In Review $, Active $, Completed $ |
| Create project modal | ✅ In Prototype | requires name; optionally links a client |
| Delete project with confirm | ✅ In Prototype | modal confirmation before delete |
| Edit project metadata | ○ Not Shown | no edit button visible in prototype |
| Pagination | ○ Not Shown | single user, volume not a concern |

### 4.2 Clients Module

| Feature | Status | Notes |
|---|---|---|
| Client list with search | ✅ In Prototype | searches all fields including notes |
| Sortable columns | ✅ In Prototype | 6 columns |
| Stats row (3 KPIs) | ✅ In Prototype | total clients, total projects, avg projects/client |
| Expandable row detail | ✅ In Prototype | click row to show address + notes inline |
| Add client modal | ✅ In Prototype | company, contact, email, phone, address, notes |
| Edit client modal | ✅ In Prototype | pre-fills form with existing data |
| Delete client with confirm | ✅ In Prototype | modal confirmation |
| Link client → projects | ○ Not Shown | no project list shown from client view |

### 4.3 Materials Database

| Feature | Status | Notes |
|---|---|---|
| Material list with search | ✅ In Prototype | name, supplier code, supplier searchable |
| Supplier tabs (filterable) | ✅ In Prototype | Massa, WBC, Triad, Pro Tech, Home Depot + add/remove |
| Price source filter | ✅ In Prototype | All / Supplier (auto-priced) / Manual |
| Sortable columns | ✅ In Prototype | 7 columns including price and date |
| Add/Edit material modal | ✅ In Prototype | name, supplier, unit, price, code, source |
| Delete material with confirm | ✅ In Prototype | notes it won't affect existing estimates |
| Add/Remove supplier tabs | ✅ In Prototype | inline input in tab bar |
| Delete supplier + all items | ✅ In Prototype | modal confirm with item count shown |
| Upload supplier pricing (CSV/XLSX) | ⚠️ Aspirational | UI present, upload handler not wired |
| last_updated tracking | ✅ In Prototype | auto-set on save |

### 4.4 Labor & Equipment Rates

| Feature | Status | Notes |
|---|---|---|
| Rate list with search | ✅ In Prototype | name, category, supplier, notes searchable |
| Category tabs (filterable) | ✅ In Prototype | Labor, Equipment Rentals, Subcontractors, Other + add/remove |
| Sortable columns | ✅ In Prototype | 7 columns |
| Category count badges | ✅ In Prototype | labor / equipment / subs counts in header |
| Add/Edit rate modal | ✅ In Prototype | name, category, supplier, rate, unit, notes |
| Delete rate with confirm | ✅ In Prototype | notes it won't affect existing estimates |
| Add/Remove category tabs | ✅ In Prototype | inline input in tab bar |
| Delete category + all rates | ✅ In Prototype | modal confirm with count |

### 4.5 Estimate Builder

| Feature | Status | Notes |
|---|---|---|
| 4-level hierarchy (Section → Sub → Group → Item) | ✅ In Prototype | all levels collapsible |
| Add / rename / delete sections | ✅ In Prototype | inline rename via input |
| Add / rename / delete subcategories | ✅ In Prototype | inline rename via input |
| Add / rename / delete component groups | ✅ In Prototype | inline input, add via "Add component group" |
| Add / delete line items | ✅ In Prototype | per subcategory or per group |
| Autocomplete from materials/rates DB | ✅ In Prototype | keyboard nav, category-filtered, populates price + unit |
| Line item type selector | ✅ In Prototype | Materials \| Labor \| Equipment \| Subs \| Other |
| Manual price override | ✅ In Prototype | typing price sets price_override flag |
| Per-line note / description | ✅ In Prototype | toggle shows/hides note row per line |
| Ungrouped items section | ✅ In Prototype | auto-rendered for items with no group |
| Global markup toolbar (5 types) | ✅ In Prototype | Materials, Labor, Equipment, Subs, Other % |
| Per-subcategory markup overrides (5 types) | ✅ In Prototype | overrides global; null = inherit |
| Per-type enable/disable toggles | ✅ In Prototype | toggle per type per subcategory |
| Lump sum adjustment per subcategory | ✅ In Prototype | added post-markup to subcategory total |
| Real-time totals (section + project) | ✅ In Prototype | recalculates on every change |
| Fixed footer summary bar | ✅ In Prototype | Materials before/after, Labor, Equipment, Subs, Other, Total |
| Undo (20-step history) | ✅ In Prototype | snapshots before structural changes |
| Save to project | ✅ In Prototype | updates project total, navigates to overview |
| Estimate persistence between sessions | ○ Not Shown | prototype loses data on reload — backend concern |
| Estimate versioning / history | ○ Not Shown | handled via quote versions, not estimate versions |

### 4.6 Job Overview

| Feature | Status | Notes |
|---|---|---|
| Client + project info card | ✅ In Prototype | name, contact, email, phone, description, dates |
| Status selector modal | ✅ In Prototype | click status badge to change, shows descriptions |
| Cost summary KPIs | ✅ In Prototype | total + per-category cards (only non-zero shown) |
| Sections breakdown table | ✅ In Prototype | per-section cost by category |
| Activity log | ⚠️ Aspirational | static sample data in prototype — wired in Phase 2 |
| Edit Estimate button | ✅ In Prototype | navigates to builder |
| Send Quote to Client | 🆕 New | generates token, sets expiry, triggers email or link copy |
| Quote version history | 🆕 New | list of v1, v2... with status badges |
| Convert to Invoice | ⚠️ Aspirational | deferred to Phase 2 |
| PDF documents (4 types) | ○ Not Shown | deferred — browser print/save handles client needs |

### 4.7 Quote Page (Public, No Auth)

This page is new — not in the prototype. Accessible at `/q/{token}`.

| Feature | Status | Notes |
|---|---|---|
| Token validation | 🆕 New | 404 if not found; "expired" state if past expires_at |
| Already-signed state | 🆕 New | shows confirmation if quote_signatures record exists |
| Quote content display | 🆕 New | project name, client name, section totals, grand total |
| Contractor branding | 🆕 New | contractor name prominent on the page |
| Acceptance form | 🆕 New | typed name field + "I Accept This Quote" button |
| Signature record on submit | 🆕 New | writes quote_signatures row with name, IP, timestamp |
| Post-sign confirmation | 🆕 New | "Thank you" state shown after signing |
| Project status auto-advance | 🆕 New | project transitions to "In Review" on signature |
| Print/Save as PDF | 🆕 New | native browser print — no library required |

---

## 5. Key Interaction Patterns

Several interaction patterns have direct implications for how the backend and frontend communicate.

### 5.1 Autocomplete Input

The autocomplete in the Estimate Builder is the most complex UI component. It requires:

- Keystroke-level filtering against a combined materials + rates list
- Category-type filtering (only show Materials items when type = Materials, etc.)
- Keyboard navigation: ArrowUp / ArrowDown, Enter to select, Escape to dismiss
- On selection: auto-populate name, unit_price, unit, category_type, material_id on the line item
- Focus handoff: after selecting, focus moves to the qty field
- Dropdown positioned absolutely within a table cell, z-index above table borders

> **htmx note:** Achievable with `hx-trigger="input"` sending to a server search endpoint returning a dropdown partial. For a single-user deployment on the same VPS the round-trip latency is acceptable. Alpine.js caching the full list client-side is also viable and avoids any server round-trip.

### 5.2 Inline Editing

Section names, subcategory names, and component group names are edited inline via input fields styled to be invisible until focused. There is no separate "edit mode" — the text and the edit control are the same element. On blur or Enter, an `hx-patch` sends the updated name to the server.

### 5.3 Undo System

Undo is implemented by snapshotting all four estimate collections (sections, subcategories, groups, line_items) before any structural mutation (add, delete, rename). The snapshot stack is capped at 20 entries. Undo replaces all four collections simultaneously — full-state rollback, not a granular command pattern.

> **htmx note:** Client-side undo via Alpine.js holding a JSON snapshot stack is the pragmatic choice. Undo history is lost on page reload, which is acceptable. Server-side undo would require storing snapshot rows in Postgres on every mutation.

### 5.4 Real-Time Calculation

Every markup %, quantity, and unit price change triggers a full recalculation of all totals. The calculation is pure and fast — no async involved.

> **htmx note:** This is the strongest argument for client-side calculation logic. Server round-trips per keystroke on a dense form will feel slow even on a local VPS. Recommended approach: Alpine.js owns the live calculation state. htmx handles structural mutations (add/delete section/subcategory/group/item) and final save. Go remains the source of truth for persisted data.

---

## 6. Quote & Client Signature

### 6.1 Workflow

```
Contractor finalizes estimate
        ↓
"Send Quote to Client" action
        ↓
System snapshots estimate + totals → creates quotes row (status: sent)
        ↓
Contractor copies link  OR  enters client email → app sends via Postmark
        ↓
Client opens /q/{token}
        ↓
Client types name → clicks "I Accept This Quote"
        ↓
System writes quote_signatures row (name, IP, timestamp)
Project status → "In Review"
Contractor sees signed status in Job Overview
```

### 6.2 Quote Lifecycle

```
draft → sent → signed
                ↓ (3 days pass, unsigned)
             expired  ←── resend resets clock (same row, new sent_at + expires_at)

sent → superseded  (when a newer version is sent on the same project)
```

Resend is a simple update — no new row, no version bump:

```sql
UPDATE quotes
SET sent_at   = now(),
    expires_at = now() + interval '3 days',
    status     = 'sent'
WHERE id = $1
```

Version only increments when the estimate content itself has changed.

### 6.3 Version Behavior

- Version numbers are scoped per project (project 7 can have v1, v2, v3 independently)
- When the contractor sends a new version, the previous `sent` quote transitions to `superseded`
- Superseded quote records and their signatures are preserved — never deleted
- `estimate_snapshot` and `totals_snapshot` are immutable once status transitions from `draft` to `sent`

### 6.4 Quote Page Content

The public quote page at `/q/{token}` shows:

- Contractor name (prominent)
- Project name and client name
- Section breakdown: section name + total per section
- Grand total (Materials, Labor, Equipment, Subs, Other, Total)
- Quote version number and sent date
- Expiry notice if unsigned
- Signature form (typed name + accept button)

Markup percentages are never shown to the client. Line-item detail is not shown in Phase 1 — totals only.

### 6.5 Email Delivery

Two options offered in the Send Quote flow:

**Copy Link** — generates the token URL, copies to clipboard. Contractor pastes into their own email. No provider dependency.

**Send via Email** — requires Postmark API key configured in app settings. App sends a plain transactional email with the quote link. Each send appends a row to `quote_emails` for delivery tracking.

Both options are available simultaneously — contractor can copy the link and also trigger the email.

### 6.6 Signature Record

The signature is not legally equivalent to DocuSign. For a small contractor context it provides:

- Snapshot of the quote content at the time of signing (jsonb)
- Signer's typed name
- Timestamp (UTC)
- IP address
- Immutable record (no updates or deletes on `quote_signatures`)

This is sufficient to establish "client X accepted quote vN on date Y" in any dispute.

---

## 7. Go + htmx Tech Stack Evaluation

### 7.1 Feature Area Assessment

| Feature Area | Fit | Rationale |
|---|---|---|
| CRUD pages (Clients, Materials, Rates) | ✅ Excellent | Standard form/table/modal patterns — htmx + templ handles this cleanly with hx-get / hx-post / hx-swap |
| Projects list + filtering + sorting | ✅ Excellent | Server-side filter/sort on GET params. No client state needed. htmx shines here |
| Project status workflow | ✅ Excellent | Simple enum PATCH. Status badge re-renders via htmx |
| Markup calculation engine | ✅ Excellent | Pure Go math. Testable, deterministic. Business logic lives in Go structs |
| Estimate structure (4-level hierarchy) | ✅ Good | Manageable with htmx + Alpine for local collapse state. Structural mutations are htmx POST/DELETE calls returning re-rendered partials |
| Real-time total recalculation | ⚠️ Moderate | Alpine.js maintains total state client-side and sends on save. Avoids server round-trips per keystroke |
| Autocomplete input | ⚠️ Moderate | Achievable with htmx hx-trigger="input" + server search returning dropdown partial. Alpine.js caching the full list is also viable |
| Inline editing (names) | ✅ Good | Alpine.js x-data with local edit state. On blur/enter, hx-patch sends to server. ~10 lines of Alpine |
| Undo (20-step history) | ⚠️ Moderate | Alpine.js holds JSON snapshot stack client-side. Lost on reload — acceptable tradeoff |
| Quote token generation | ✅ Excellent | crypto/rand is stdlib Go. Nanoid-style tokens are trivial |
| Quote expiry + resend logic | ✅ Excellent | Timestamptz arithmetic in Postgres. Cron or check-on-read for expiry transitions |
| Public quote page (/q/{token}) | ✅ Excellent | Simple server-rendered page, no auth, token lookup, conditional states |
| E-signature capture | ✅ Excellent | Single POST handler — write signature row, update project status. Straightforward |
| Email delivery (Postmark) | ✅ Excellent | Clean Go HTTP client. Postmark has a simple REST API. No SDK required |
| Supplier pricing CSV/XLSX upload | ✅ Excellent | encoding/csv is stdlib. excelize handles XLSX. Clean Go territory |
| Activity log | ✅ Excellent | Append-only event table in Postgres. Triggered by middleware or explicit service layer calls |
| Print/Save as PDF | ✅ Excellent | Native browser print — zero server involvement |
| Dark theme / design system | ✅ Neutral | Purely a CSS concern. Tailwind works identically regardless of backend |

### 7.2 The One Genuine Challenge

The Estimate Builder is the only page that pushes against htmx's natural grain. It combines:

- Dense tabular data entry with per-row interactivity
- Cascading calculations that must update immediately on any input change
- Complex keyboard-driven autocomplete
- Client-side undo that would be awkward to persist on every keystroke

**Recommended approach:** Alpine.js owns live calculation state and undo history on the client. htmx handles structural mutations (add/delete section/subcategory/group/item) and the final save. Go is the source of truth for all persisted data. This keeps the architecture coherent without fighting the browser for every keystroke.

### 7.3 Verdict

Go + htmx is a sound choice for this application. 15 of 18 feature areas are an excellent or good fit. The three moderate-fit areas all have well-understood solutions using Alpine.js for client-side state. The quote and signature system is particularly well-suited to Go — token generation, timestamptz logic, immutable records, and a simple Postmark HTTP client are all idiomatic Go. The prototype being React does not change this conclusion.

---

## 8. Scope Summary & Phase Recommendation

### Phase 1 — Core Product (Billable v1)

- Auth (single user, simple session)
- Projects: full CRUD, status workflow, dashboard stats
- Clients: full CRUD
- Materials DB: full CRUD, supplier tabs, last_updated tracking
- Labor & Equipment Rates: full CRUD, category tabs
- Estimate Builder: full 4-level hierarchy, markup engine, autocomplete, save to project
- Job Overview: client info, cost summary, sections breakdown, status change, quote version list
- Quote system: send quote (copy link + Postmark email), public `/q/{token}` page, e-signature capture, expiry + resend
- Browser print/save as PDF on quote page

### Phase 2 — Extended (Post-Launch)

- Activity log (real, wired to mutations)
- Supplier pricing CSV/XLSX upload
- Estimate versioning / revision history
- Convert to Invoice
- Internal Estimate PDF (server-generated)
- Crew Material List PDF
- Supplier Order List PDF (requires cross-section deduplication logic)
- Richer quote page (line-item detail toggle)

### Out of Scope

- Multi-user / permissions (single deployment per client)
- Pagination (single user, volume not a concern)
- Third-party integrations (QuickBooks, etc.)
- Legally binding e-signature (DocuSign-level)

---

---

## 9. Estimate Builder — Component Architecture

The Estimate Builder is the only page in Skalkaho that does not follow the standard Go + htmx + templ pattern. Its combination of cascading live calculations, deep mutable nested state, keyboard-driven autocomplete, and client-side undo makes it a genuine SPA component rather than a server-rendered form. Attempting to build it in htmx + Alpine would produce fragile, hard-to-maintain code.

### 9.1 Technology Decision

**Svelte 5** is selected for the Estimate Builder component. Reasons specific to this component:

- The 4-level nested state tree (Section → Subcategory → Group → Line Item) requires frequent deep mutations. Svelte 5 tracks direct mutations reactively — no immutable update boilerplate required.
- Svelte compiles to a small vanilla JS bundle (~10KB framework overhead vs ~45KB for React + ReactDOM).
- Svelte is already the stated preference for complex interactive UIs within Firefly Software's stack.

For comparison, updating a single line item's quantity in React requires spreading the entire tree. In Svelte 5 it is a direct assignment: `item.quantity = newQty`. This difference compounds across every interaction in the builder.

### 9.2 Integration Pattern — Option B (Embedded Island)

The Estimate Builder is a **mounted island** inside an otherwise server-rendered page. Go owns routing, layout, and the page chrome. Svelte owns everything inside the mount point.

```
GET /projects/{id}/estimate
  └── Go renders full page via templ
       ├── Sidebar (server-rendered, htmx)
       ├── Topbar (server-rendered, htmx)
       └── <div id="estimate-root" data-project-id="{{.Project.ID}}"></div>
            └── Svelte component mounts here
                 └── fetches GET /api/estimate/{project_id}
                      └── renders entire builder UI
```

The Svelte bundle is compiled by Vite and output into Go's static assets directory. The templ layout loads it as a standard `<script type="module">` tag only on the estimate builder page.

This pattern keeps Go in charge of navigation, sidebar state, and the "Current Project" context. Switching between projects is a normal server-rendered page navigation — no SPA routing concerns.

### 9.3 JSON API

Two endpoints form the entire API surface for the Estimate Builder. No other endpoints are required.

**`GET /api/estimate/{project_id}`**

Returns the full estimate payload in a single response. The `materials_db` and `rates_db` arrays are included so the autocomplete has everything it needs client-side without additional requests. For a single-user deployment with a few hundred materials this payload will be well under 100KB.

**`POST /api/estimate/{project_id}`**

Accepts the full estimate payload as the request body. Go validates the payload, then persists all tables (sections, subcategories, component_groups, line_items) inside a single Postgres transaction. If any part fails the transaction rolls back, the in-memory component state is unchanged, and an error is returned. On success, the saved payload is returned and the component updates its reference state.

### 9.4 Payload Shape

The JSON contract between Go and Svelte:

```json
{
  "project": {
    "id": "abc123",
    "name": "Missoula Pole Barn",
    "status": "Draft"
  },
  "globals": {
    "materials_markup": 20,
    "labor_markup": 25,
    "equipment_markup": 15,
    "subs_markup": 10,
    "other_markup": 10
  },
  "sections": [
    {
      "id": "s1",
      "name": "Framing",
      "sort_order": 0,
      "subcategories": [
        {
          "id": "sub1",
          "name": "Wall Framing",
          "sort_order": 0,
          "lump_sum": 0,
          "markup_overrides": {
            "materials": null,
            "labor": null,
            "equipment": null,
            "subs": null,
            "other": null
          },
          "markup_enabled": {
            "materials": true,
            "labor": true,
            "equipment": true,
            "subs": true,
            "other": true
          },
          "component_groups": [],
          "line_items": []
        }
      ]
    }
  ],
  "materials_db": [],
  "rates_db": []
}
```

The payload is nested (sections own subcategories, subcategories own groups and line items) to mirror the UI hierarchy directly. Go flattens this back into separate tables on POST.

### 9.5 ID Strategy

All entities — including newly created ones — carry a `nanoid`-generated string ID from the moment of creation. The client generates IDs using the `nanoid` JS library (~130 bytes). Go treats all IDs as strings and performs upserts on save, so it does not distinguish between client-generated and server-generated IDs.

This keeps ID handling consistent throughout the component. There are no "temp ID" special cases, no conditional logic for new-vs-saved items, and no ID remapping after the first save.

### 9.6 Auto-Save Behavior

The component auto-saves using a debounced POST. The flow:

```
User makes any change
  → isDirty = true
  → reset 2-second debounce timer
  → timer fires → POST /api/estimate/{project_id}
  → on success: isDirty = false, record saved_at timestamp
  → on failure: isDirty remains true, show error indicator, do not clear timer
```

**Save status indicator** (visible in the builder topbar):

| State | Display |
|---|---|
| Clean | `Saved just now` (or timestamp) |
| Dirty, timer pending | `Unsaved changes` |
| In flight | `Saving...` |
| Error | `Save failed — retrying` |

### 9.7 Dirty State & Navigation Warnings

Two navigation scenarios require a dirty-state check:

**Browser-level navigation** (tab close, address bar, browser back): handled by `window.beforeunload`. If `isDirty` is true, the browser shows its native "Leave site?" dialog.

**In-app navigation** (sidebar links): htmx navigation bypasses `beforeunload`. Sidebar links on the estimate builder page are intercepted — if `isDirty` is true, a confirm modal is shown before the navigation proceeds. The Svelte component exposes `isDirty` as a signal readable by the parent templ page to enable this intercept.

### 9.8 Undo System

20-step undo using a client-side snapshot stack. Before any structural mutation (add, delete, rename, reorder of any entity at any level), the full current state is serialized to JSON and pushed onto the stack. Undo replaces the entire current state with the previous snapshot — full rollback, not a granular command pattern.

The stack is capped at 20 entries. Undo history is lost on page reload. This is an acceptable tradeoff — the auto-save ensures persisted state is always recent, and undo is a within-session convenience rather than a recovery mechanism.

### 9.9 Build Order

Each step is independently testable and represents a logical commit boundary:

| Step | Description |
|---|---|
| 1 | Go API endpoints — GET and POST handlers, payload structs, sqlc queries. Tested with curl, no Svelte yet. |
| 2 | Svelte scaffold — Vite + Svelte 5, bundle output to Go static assets, templ shell mounts component, renders project name from fetched payload. |
| 3 | Read-only estimate display — full nested hierarchy rendered from API response. No editing. Confirms data shape end to end. |
| 4 | Markup engine — pure calculation functions in Svelte, unit tested in isolation. No UI changes. |
| 5 | Line item editing — quantity, unit price, type selector. Auto-save fires. Totals update live. |
| 6 | Add / delete line items — with autocomplete from materials_db / rates_db. |
| 7 | Add / delete / rename sections, subcategories, component groups. |
| 8 | Markup controls — global toolbar, per-subcategory overrides, enable toggles, lump sum. |
| 9 | Undo system — 20-step snapshot stack, Ctrl+Z binding. |
| 10 | Dirty state & save indicator — beforeunload, sidebar intercept, save status UI. |

---

*End of Specification*
