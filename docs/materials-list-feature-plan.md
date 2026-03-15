# Job Site Materials List — Implementation Plan

## Overview

Generate persistent, printable materials lists from an accepted quote for use
on the job site. Lists are organized by category → subcategory and contain
quantity and unit only — no prices, no supplier information. Three generation
modes mirror the supply list feature:

- **Entire quote** — one list covering all categories and subcategories
- **Per category** — one list per top-level category in a single action
- **Single category** — owner selects a specific category

Lists are **live** — they reference the quote's current line items rather than
snapshotting values. If the quote changes, the list reflects the change on next
view. This prioritizes job-site accuracy over record stability.

A nullable `supply_list_id` foreign key is included for future supply
association features but defaults to NULL and is not used in this implementation.

---

## How This Differs From Supply Lists

| | Supply List | Materials List |
|---|---|---|
| Audience | Owner / procurement | Foreman / job site |
| Organized by | Supplier | Category → Subcategory |
| Shows prices | Yes | No |
| Shows supplier | Yes | No |
| Data model | Snapshot | Live (query at render time) |
| Generation scope | By supplier | By category |

---

## Data Model

Because the list is live, `materials_lists` stores only metadata — the scope
of the list (quote, optional category) and display/tracking fields. Items are
queried from the quote at render time, filtered by the list's scope.

```
materials_lists    (metadata + scope)
    ↓ query at render time
quote line_items → materials    (live data, filtered by category/subcategory)
```

---

## Step 1 — Database Migration

**File:** `migrations/XXXX_add_materials_lists.sql`

```sql
CREATE TABLE materials_lists (
    id              TEXT PRIMARY KEY,           -- nanoid
    quote_id        TEXT NOT NULL REFERENCES quotes(id) ON DELETE RESTRICT,
    section_id      TEXT REFERENCES sections(id) ON DELETE RESTRICT,
                                                -- NULL = entire quote list
    supply_list_id  TEXT REFERENCES supply_lists(id) ON DELETE SET NULL,
                                                -- reserved for future use, always NULL for now
    label           TEXT NOT NULL,              -- "All Categories", or section name
    notes           TEXT,
    generated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_materials_lists_quote ON materials_lists(quote_id);
CREATE INDEX idx_materials_lists_section ON materials_lists(quote_id, section_id);
```

No items table — items are queried live from the quote at render time.

**Test:** Create a materials list scoped to a quote, confirm `ON DELETE RESTRICT`
blocks quote deletion. Create one scoped to a section, confirm `ON DELETE
RESTRICT` blocks section deletion. Set `supply_list_id`, delete the supply list,
confirm it goes to NULL (not a restrict — supply list deletion should not block
materials list existence).

---

## Step 2 — sqlc Queries

**File:** `db/query/materials_lists.sql`

### Materials List Metadata

```sql
-- name: CreateMaterialsList :one
INSERT INTO materials_lists (id, quote_id, section_id, label, notes)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetMaterialsList :one
SELECT
    ml.*,
    s.name AS section_name
FROM materials_lists ml
LEFT JOIN sections s ON s.id = ml.section_id
WHERE ml.id = ?;

-- name: ListMaterialsListsByQuote :many
SELECT
    ml.*,
    s.name AS section_name
FROM materials_lists ml
LEFT JOIN sections s ON s.id = ml.section_id
WHERE ml.quote_id = ?
ORDER BY ml.section_id NULLS FIRST, ml.generated_at DESC;

-- name: UpdateMaterialsListNotes :one
UPDATE materials_lists
SET notes = ?
WHERE id = ?
RETURNING *;

-- name: DeleteMaterialsList :exec
DELETE FROM materials_lists WHERE id = ?;
```

### Live Item Queries

These queries power the print view — called at render time, not stored.

```sql
-- name: GetMaterialsListItems :many
-- Full quote list — all sections, subcategories, line items with materials
-- Organized by section → subcategory for the foreman view
SELECT
    s.name          AS section_name,
    s.sort_order    AS section_order,
    sc.name         AS subcategory_name,
    sc.sort_order   AS subcategory_order,
    li.name         AS item_name,
    li.quantity,
    li.unit,
    li.sort_order   AS item_order
FROM line_items li
JOIN component_groups cg ON cg.id = li.component_group_id
JOIN subcategories sc    ON sc.id = cg.subcategory_id
JOIN sections s          ON s.id  = sc.section_id
WHERE s.quote_id = ?
  AND li.quantity > 0
ORDER BY
    s.sort_order  ASC,
    sc.sort_order ASC,
    li.sort_order ASC;

-- name: GetMaterialsListItemsBySection :many
-- Single section list — filtered to one section
SELECT
    sc.name         AS subcategory_name,
    sc.sort_order   AS subcategory_order,
    li.name         AS item_name,
    li.quantity,
    li.unit,
    li.sort_order   AS item_order
FROM line_items li
JOIN component_groups cg ON cg.id = li.component_group_id
JOIN subcategories sc    ON sc.id = cg.subcategory_id
WHERE sc.section_id = ?
  AND li.quantity > 0
ORDER BY
    sc.sort_order ASC,
    li.sort_order ASC;
```

**Note on ad hoc items:** Line items with no material reference are included —
the foreman needs to know about all materials regardless of whether they're
linked to a supplier record. The query pulls from `line_items` directly, not
via `materials`, so ad hoc items are naturally included.

Run `sqlc generate`, confirm no errors.

**Test:** Run `GetMaterialsListItems` against a quote with items across multiple
sections, subcategories, and component groups. Confirm ordering is correct.
Confirm items with `quantity = 0` are excluded. Confirm ad hoc items (no
material reference) are included.

---

## Step 3 — Generator

**File:** `handler/materialslist/generate.go`

Unlike the supply list generator, this does not insert item rows — it only
creates the `materials_list` metadata record. Items are queried live at render
time. The generator's job is to determine scope, build the label, and insert
the record.

```go
// GenerateAll creates one materials list covering the entire quote.
func GenerateAll(ctx context.Context, q *db.Queries, quoteID string) (*db.MaterialsList, error)

// GeneratePerSection creates one materials list per top-level section on the quote.
// Returns created lists in section sort_order.
func GeneratePerSection(ctx context.Context, q *db.Queries, quoteID string) ([]*db.MaterialsList, error)

// GenerateForSection creates a materials list for a single specified section.
func GenerateForSection(ctx context.Context, q *db.Queries, quoteID, sectionID string) (*db.MaterialsList, error)
```

### Label generation

- `GenerateAll` → label = `"All Categories"`
- `GeneratePerSection` → label = section name (e.g. `"Framing"`)
- `GenerateForSection` → label = section name

No transaction needed — each call is a single INSERT. If `GeneratePerSection`
fails partway through, the created lists remain (partial generation is
acceptable — the owner can delete and regenerate).

**Test:** Generate all three types against a quote with three sections. Confirm
correct labels. Confirm `section_id` is NULL for the all-categories list.
Confirm `GeneratePerSection` creates exactly one list per section.

---

## Step 4 — HTTP Handler

**File:** `handler/materialslist/materialslist.go`

| Method | Route | Action |
|--------|-------|--------|
| GET | `/quotes/{id}/materials-lists` | Index — all materials lists for a quote |
| POST | `/quotes/{id}/materials-lists/generate/all` | Generate entire quote list |
| POST | `/quotes/{id}/materials-lists/generate/per-section` | Generate one list per section |
| POST | `/quotes/{id}/materials-lists/generate/section/{sid}` | Generate for one section |
| GET | `/quotes/{id}/materials-lists/{lid}` | View / print a materials list |
| PATCH | `/quotes/{id}/materials-lists/{lid}/notes` | Update notes |
| DELETE | `/quotes/{id}/materials-lists/{lid}` | Delete a materials list |

### Show handler

The show handler does two things:
1. Load the `materials_list` metadata record
2. Query live items using `GetMaterialsListItems` or
   `GetMaterialsListItemsBySection` depending on whether `section_id` is set
3. Pass both to the template for rendering

Because items are live, no caching or snapshot logic is needed — every page
load reflects current quote state.

### Guard: quote must be accepted

Same as supply lists — block generation against a draft quote with a clear
error message. Check quote status at the handler level before calling the
generator.

**Test:** Generate all three list types, confirm correct item sets are returned
in the show view. Edit a line item quantity on the quote, reload the materials
list, confirm the new quantity appears. Attempt generation on a draft quote,
confirm it is blocked.

---

## Step 5 — Templates

**Directory:** `templates/materials-lists/`

### `index.html` — Materials list index for a quote

```
Materials Lists — Project: Riverside Remodel

[Generate Entire List]  [Generate Per Category]

──────────────────────────────────────────────────────────
All Categories     Generated Jan 15, 2026     [View] [Delete]

Framing            Generated Jan 15, 2026     [View] [Delete]
Electrical         Generated Jan 15, 2026     [View] [Delete]
Plumbing           Generated Jan 15, 2026     [View] [Delete]
──────────────────────────────────────────────────────────

[Generate for Specific Category ▾]   ← dropdown of sections on this quote
```

Note: unlike supply lists there are no dates on the index — materials lists
have no delivery tracking. Notes are visible only on the detail view.

### `show.html` — Single materials list (view + print)

**Header:**
```
MATERIALS LIST
Project:   Riverside Remodel
Category:  Framing            (or "All Categories")
Generated: January 15, 2026
Note: Quantities reflect current quote — verify before ordering.
```

**Entire quote view** — grouped by section then subcategory:

```
FRAMING
────────────────────────────────────────
  Wall Framing
  ──────────────────────────────────────
  Item                        Qty    Unit
  2x4x8 Douglas Fir           120    ea
  2x6x10 Douglas Fir           40    ea
  3/4" Plywood Sheathing       24    sheet

  Roof Framing
  ──────────────────────────────────────
  Item                        Qty    Unit
  2x10x16 Ridge Board           2    ea
  Rafter Ties                  36    ea

ELECTRICAL
────────────────────────────────────────
  Rough-In
  ──────────────────────────────────────
  Item                        Qty    Unit
  14 AWG Wire                   3    roll
  Junction Boxes               12    ea
```

**Single section view** — same layout but section header is omitted (it's
already the title), subcategories are the top-level grouping:

```
MATERIALS LIST — FRAMING

  Wall Framing
  ──────────────────────────────────────
  Item                        Qty    Unit
  2x4x8 Douglas Fir           120    ea
  ...
```

### Print styles

Same `@media print` approach as supply lists:
- Hide nav, action buttons, notes edit form
- Keep section/subcategory hierarchy visually clear
- Bold section headers, normal weight subcategory headers
- Clean table borders

Add a footer note on the printed page:
```
* Quantities reflect current quote state at time of printing.
* Verify with project manager before placing orders.
```

**Test:** View an entire quote list, confirm section → subcategory grouping.
View a single section list, confirm subcategory-only grouping (no redundant
section header). Print both, confirm chrome is hidden.

---

## Step 6 — Router Wiring

**File:** `router/router.go`

```go
ml := materialslist.NewHandler(queries)

r.HandleFunc("GET  /quotes/{id}/materials-lists",                           ml.Index)
r.HandleFunc("POST /quotes/{id}/materials-lists/generate/all",              ml.GenerateAll)
r.HandleFunc("POST /quotes/{id}/materials-lists/generate/per-section",      ml.GeneratePerSection)
r.HandleFunc("POST /quotes/{id}/materials-lists/generate/section/{sid}",    ml.GenerateForSection)
r.HandleFunc("GET  /quotes/{id}/materials-lists/{lid}",                     ml.Show)
r.HandleFunc("PATCH /quotes/{id}/materials-lists/{lid}/notes",              ml.UpdateNotes)
r.HandleFunc("DELETE /quotes/{id}/materials-lists/{lid}",                   ml.Delete)
```

Add a "Materials Lists" tab on the quote detail page alongside "Supply Lists".
Both tabs visible only when quote is accepted. Consider grouping them under a
"Job Documents" section on the quote page so the pattern is clear as more
document types are added.

**Test:** Full end-to-end — accept a quote, generate all three list types, view
each, update a quote line item, reload a list and confirm updated quantity,
print a list, delete a list.

---

## Commit Sequence

| # | Commit message |
|---|----------------|
| 1 | `db: add materials_lists migration` |
| 2 | `db: add sqlc queries for materials lists and live item queries` |
| 3 | `handler/materialslist: add generator (GenerateAll, GeneratePerSection, GenerateForSection)` |
| 4 | `handler/materialslist: add HTTP handler and route actions` |
| 5 | `templates: add materials list index, show/print view` |
| 6 | `router: wire materials list routes, add Job Documents section to quote detail` |

---

## Future Considerations

- **Supply list association** — the `supply_list_id` FK is already in place.
  When this is implemented, a materials list could show which supply list
  covers each item, giving the foreman visibility into what has been ordered
  and what hasn't.
- **Checked-off items** — a foreman might want to mark items as received or
  staged on the job site. Would require a `materials_list_checks` table keyed
  by `(materials_list_id, line_item_id)`. The live model makes this slightly
  tricky — if a line item is deleted from the quote, its check record becomes
  orphaned. Worth thinking through carefully if this is requested.
- **Per-subcategory generation** — current design generates at the section
  level. If a section has many subcategories and different crews handle each,
  per-subcategory lists could be added following the same pattern.
