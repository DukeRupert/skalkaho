# Supply List Feature — Implementation Plan

## Overview

Allow the owner to generate persistent supply lists from an accepted quote.
Lists are snapshots of line item quantities and prices at generation time —
subsequent quote edits do not affect existing lists. Three generation modes:

- **All suppliers** — one list covering every supplier on the quote
- **Per supplier (all)** — generates one list per supplier in a single action
- **Per supplier (individual)** — owner selects a specific supplier

Each list tracks an estimated delivery date and an actual delivery date.
Output is a printable/PDF page within Skalkaho.

Ad hoc line items with no supplier reference are grouped under "No Supplier"
and appear on the all-suppliers list only.

---

## Data Model

```
supply_lists
  └── supply_list_items   (snapshot of line items at generation time)
```

A `supply_list` is scoped to one quote and optionally one supplier. The
all-suppliers list has `supplier_id = NULL`. Per-supplier lists have a
specific `supplier_id`.

---

## Step 1 — Database Migration

**File:** `migrations/XXXX_add_supply_lists.sql`

```sql
CREATE TABLE supply_lists (
    id                  TEXT PRIMARY KEY,       -- nanoid
    quote_id            TEXT NOT NULL REFERENCES quotes(id) ON DELETE RESTRICT,
    supplier_id         TEXT REFERENCES suppliers(id) ON DELETE RESTRICT,
                                                -- NULL = all-suppliers list
    label               TEXT NOT NULL,          -- display name, e.g. "All Suppliers" or supplier name
    estimated_delivery  DATE,                   -- expected delivery date
    actual_delivery     DATE,                   -- filled in when goods arrive
    generated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes               TEXT
);

-- A quote should have at most one all-suppliers list and one list per supplier
-- at any given time. Not enforced by unique constraint (owner may regenerate),
-- but the UI should warn if one already exists.
CREATE INDEX idx_supply_lists_quote_id ON supply_lists(quote_id);
CREATE INDEX idx_supply_lists_supplier ON supply_lists(quote_id, supplier_id);

CREATE TABLE supply_list_items (
    id               TEXT PRIMARY KEY,          -- nanoid
    supply_list_id   TEXT NOT NULL REFERENCES supply_lists(id) ON DELETE CASCADE,
    -- Snapshot fields — copied from line item + material at generation time
    material_id      TEXT,                      -- nullable, for reference only (ad hoc items have none)
    supplier_id      TEXT,                      -- nullable (ad hoc items)
    supplier_name    TEXT,                      -- snapshotted, not a FK
    material_name    TEXT NOT NULL,
    supplier_code    TEXT,                      -- snapshotted supplier SKU
    quantity         REAL NOT NULL,
    unit             TEXT NOT NULL,
    unit_price       REAL NOT NULL DEFAULT 0,
    extended_price   REAL NOT NULL DEFAULT 0,   -- quantity * unit_price, computed at snapshot time
    sort_order       INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_supply_list_items_list ON supply_list_items(supply_list_id);
```

**Test:** Create a supply list with items, confirm CASCADE delete removes items
when list is deleted. Confirm `ON DELETE RESTRICT` on `quotes` prevents deleting
a quote that has supply lists. Delete a supplier from the suppliers table —
confirm `ON DELETE RESTRICT` blocks it (protecting snapshot integrity).

---

## Step 2 — sqlc Queries

**File:** `db/query/supply_lists.sql`

### Supply Lists

```sql
-- name: CreateSupplyList :one
INSERT INTO supply_lists (id, quote_id, supplier_id, label, estimated_delivery, notes)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetSupplyList :one
SELECT * FROM supply_lists WHERE id = ?;

-- name: ListSupplyListsByQuote :many
SELECT
    sl.*,
    COUNT(sli.id) AS item_count,
    SUM(sli.extended_price) AS total_value
FROM supply_lists sl
LEFT JOIN supply_list_items sli ON sli.supply_list_id = sl.id
WHERE sl.quote_id = ?
GROUP BY sl.id
ORDER BY sl.supplier_id NULLS FIRST, sl.generated_at DESC;

-- name: UpdateSupplyListDates :one
UPDATE supply_lists
SET estimated_delivery = ?,
    actual_delivery    = ?
WHERE id = ?
RETURNING *;

-- name: UpdateSupplyListNotes :one
UPDATE supply_lists
SET notes = ?
WHERE id = ?
RETURNING *;

-- name: DeleteSupplyList :exec
DELETE FROM supply_lists WHERE id = ?;
```

### Supply List Items

```sql
-- name: CreateSupplyListItem :one
INSERT INTO supply_list_items (
    id, supply_list_id, material_id, supplier_id, supplier_name,
    material_name, supplier_code, quantity, unit, unit_price, extended_price, sort_order
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListSupplyListItems :many
SELECT * FROM supply_list_items
WHERE supply_list_id = ?
ORDER BY supplier_name NULLS LAST, sort_order ASC, material_name ASC;

-- name: ListSupplyListItemsGrouped :many
-- Returns items with supplier grouping for the all-suppliers print view
SELECT * FROM supply_list_items
WHERE supply_list_id = ?
ORDER BY
    CASE WHEN supplier_name IS NULL THEN 1 ELSE 0 END,  -- No Supplier last
    supplier_name ASC,
    material_name ASC;
```

Run `sqlc generate`, confirm no errors.

**Test:** Confirm `ListSupplyListItemsGrouped` correctly sorts ad hoc items
(no supplier) to the end.

---

## Step 3 — Generator

**File:** `handler/supplylist/generate.go`

The core logic. Reads the quote's line items, joins to materials and suppliers,
and writes a snapshot into `supply_lists` + `supply_list_items`. All inserts
run in a single transaction — if anything fails, nothing is written.

```go
// GenerateAll creates one supply list covering all suppliers on the quote.
func GenerateAll(ctx context.Context, q *db.Queries, tx *sql.Tx, quoteID string) (*db.SupplyList, error)

// GeneratePerSupplier creates one supply list per supplier found on the quote.
// Returns the created lists in supplier name order.
func GeneratePerSupplier(ctx context.Context, q *db.Queries, tx *sql.Tx, quoteID string) ([]*db.SupplyList, error)

// GenerateForSupplier creates a supply list for a single specified supplier.
func GenerateForSupplier(ctx context.Context, q *db.Queries, tx *sql.Tx, quoteID, supplierID string) (*db.SupplyList, error)
```

### Snapshot logic (shared by all three)

1. Query all line items for the quote that have a material reference, join to
   `materials` and `suppliers` to get current name, supplier, code, unit, price
2. Query all ad hoc line items (no material reference) separately — these have
   no supplier and go on the all-suppliers list only
3. For `GenerateAll`: include everything, group by supplier in sort order
4. For `GeneratePerSupplier`: group by supplier, exclude ad hoc items, one
   `supply_list` per group
5. For `GenerateForSupplier`: filter to the specified supplier only
6. Compute `extended_price = quantity * unit_price` at snapshot time
7. Copy supplier name as a string (not FK) — snapshot is independent of future
   supplier record changes

### Handling ad hoc items

Ad hoc line items (no material reference) appear only on the all-suppliers list,
grouped under "No Supplier" at the bottom. They are excluded from per-supplier
lists since there is no supplier to assign them to. The print view should make
this grouping visually clear.

**Test:** Generate all three list types from a quote with:
- Items from two suppliers
- One ad hoc item with no supplier
Confirm all-suppliers list contains all items including ad hoc. Confirm
per-supplier lists exclude ad hoc items. Confirm per-supplier list for
Supplier A contains only Supplier A items. Confirm snapshot values don't
change when material prices are updated after generation.

---

## Step 4 — HTTP Handler

**File:** `handler/supplylist/supplylist.go`

| Method | Route | Action |
|--------|-------|--------|
| GET | `/quotes/{id}/supply-lists` | List all supply lists for a quote |
| POST | `/quotes/{id}/supply-lists/generate/all` | Generate all-suppliers list |
| POST | `/quotes/{id}/supply-lists/generate/per-supplier` | Generate one list per supplier |
| POST | `/quotes/{id}/supply-lists/generate/supplier/{sid}` | Generate for one supplier |
| GET | `/quotes/{id}/supply-lists/{lid}` | View / print a supply list |
| PATCH | `/quotes/{id}/supply-lists/{lid}/dates` | Update estimated/actual delivery dates |
| PATCH | `/quotes/{id}/supply-lists/{lid}/notes` | Update notes |
| DELETE | `/quotes/{id}/supply-lists/{lid}` | Delete a supply list |

### Generate handlers

All three generation routes:
1. Wrap the generator call in a transaction
2. On success → redirect to the supply list index for the quote (or to the
   newly created list for single-supplier generation)
3. If a list for this supplier already exists on this quote, show a
   confirmation prompt before regenerating ("A supply list for this supplier
   already exists. Regenerate?")

### Date update

`PATCH /dates` accepts `estimated_delivery` and `actual_delivery` as form
fields. Either can be blank (clearing the date). Returns updated date fragment
via HTMX swap. Setting `actual_delivery` signals the order was received.

**Test:** Generate all three list types, confirm redirect behavior. Update
dates, confirm HTMX swap. Delete a list, confirm items are cascade deleted.

---

## Step 5 — Templates

**Directory:** `templates/supply-lists/`

### `index.html` — Supply list index for a quote

Shows all supply lists for a quote:

```
Supply Lists — Project: Riverside Remodel

[Generate All Suppliers]  [Generate Per Supplier]

─────────────────────────────────────────────────────
All Suppliers          23 items   $14,820   Generated Jan 15   [View] [Delete]
  Est. delivery: Feb 1, 2026
  Actual delivery: —

Acme Lumber            12 items   $8,400    Generated Jan 15   [View] [Delete]
  Est. delivery: Jan 28, 2026
  Actual delivery: Jan 30, 2026 ✓

Henderson Electric      8 items   $4,200    Generated Jan 15   [View] [Delete]
  Est. delivery: Feb 3, 2026
  Actual delivery: —
─────────────────────────────────────────────────────

[Generate for Specific Supplier ▾]   ← dropdown of suppliers on this quote
```

### `show.html` — Single supply list (view + print)

The printable view. Renders cleanly for `@media print` — hide nav, action
buttons, and chrome.

**Header section:**
```
SUPPLY LIST
Project:    Riverside Remodel
Supplier:   Acme Lumber            (or "All Suppliers")
Generated:  January 15, 2026
Est. Delivery: February 1, 2026
Act. Delivery: —
Notes:      [editable inline via HTMX, hidden on print]
```

**All-suppliers view** — grouped by supplier with subtotals:

```
ACME LUMBER
─────────────────────────────────────────────────────────────────
SKU          Item                    Qty    Unit    Unit $    Total
LBR-2x4-8   2x4x8 Douglas Fir       120    ea      $8.50   $1,020.00
LBR-PLY-34  3/4" Plywood            40     sheet   $62.00  $2,480.00
                                              Subtotal:   $3,500.00

HENDERSON ELECTRIC
─────────────────────────────────────────────────────────────────
SKU          Item                    Qty    Unit    Unit $    Total
ELC-14AWG   14 AWG Wire (250ft)      3      roll    $89.00    $267.00
                                              Subtotal:     $267.00

NO SUPPLIER (Ad Hoc Items)
─────────────────────────────────────────────────────────────────
             Miscellaneous hardware   1      lot     $0.00      $0.00

                                        GRAND TOTAL:   $3,767.00
```

**Per-supplier view** — same layout but single supplier block, no "No Supplier"
section.

### `_dates_form.html` (partial)

HTMX-swappable date editing form:

```
Est. Delivery: [date input]   Act. Delivery: [date input]   [Save]
```

**Print styles:** Add a `print.css` or `@media print` block that:
- Hides navigation, action buttons, date edit forms
- Removes page chrome
- Ensures table borders print clearly
- Adds the Skalkaho logo/company name to the printed header

**Test:** View a supply list, confirm items are grouped and sorted correctly.
Print the page, confirm nav and edit controls are hidden. Update dates inline,
confirm HTMX swap without page reload.

---

## Step 6 — Router Wiring

**File:** `router/router.go`

```go
sl := supplylist.NewHandler(queries, rawDB)

r.HandleFunc("GET  /quotes/{id}/supply-lists",                          sl.Index)
r.HandleFunc("POST /quotes/{id}/supply-lists/generate/all",             sl.GenerateAll)
r.HandleFunc("POST /quotes/{id}/supply-lists/generate/per-supplier",    sl.GeneratePerSupplier)
r.HandleFunc("POST /quotes/{id}/supply-lists/generate/supplier/{sid}",  sl.GenerateForSupplier)
r.HandleFunc("GET  /quotes/{id}/supply-lists/{lid}",                    sl.Show)
r.HandleFunc("PATCH /quotes/{id}/supply-lists/{lid}/dates",             sl.UpdateDates)
r.HandleFunc("PATCH /quotes/{id}/supply-lists/{lid}/notes",             sl.UpdateNotes)
r.HandleFunc("DELETE /quotes/{id}/supply-lists/{lid}",                  sl.Delete)
```

Add a "Supply Lists" tab or link on the quote detail page — visible only when
the quote status is accepted/signed. Generating supply lists against a draft
quote should be blocked at the handler level with a clear error.

**Test:** Full end-to-end — accept a quote, generate all three list types,
update dates on each, print one, delete one. Attempt to generate a supply list
against a draft quote, confirm it is blocked.

---

## Commit Sequence

| # | Commit message |
|---|----------------|
| 1 | `db: add supply_lists and supply_list_items migration` |
| 2 | `db: add sqlc queries for supply lists and items` |
| 3 | `handler/supplylist: add snapshot generator (GenerateAll, GeneratePerSupplier, GenerateForSupplier)` |
| 4 | `handler/supplylist: add HTTP handler and route actions` |
| 5 | `templates: add supply list index, show/print view, and date edit partial` |
| 6 | `router: wire supply list routes, add link from quote detail page` |

---

## Future Considerations

- **PDF export** — if a proper PDF file (vs. browser print) is needed later,
  the existing PDF skill/tooling can render the `show.html` template to a file.
  The print-clean template design makes this straightforward to add.
- **Email to supplier** — supply list could be emailed directly to a supplier
  contact. Would require supplier contact records (email field on suppliers
  table) and an email sending integration.
- **Quantity adjustments** — if the owner wants to adjust quantities on a
  supply list after generation (e.g. ordering extra for waste), a per-item
  `ordered_quantity` field could be added alongside the snapshotted `quantity`.
