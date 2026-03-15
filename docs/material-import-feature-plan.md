# Supplier Material Import — Implementation Plan

## Overview

Allow the owner to import supplier product spreadsheets (Excel) to create and
update materials in the database. Parsing and Claude matching run in a background
worker so the HTTP request returns immediately. The owner monitors progress via
HTMX polling and resolves any low-confidence matches in a review queue.

This feature also establishes the first background worker pattern in the app —
subsequent background jobs (e.g. web scraping) should follow the same harness.

---

## Architecture Summary

```
HTTP Handler                 Background Worker              Claude API
─────────────────            ──────────────────────         ────────────
1. Receive upload       →    3. Poll for pending jobs
2. Write temp file           4. Parse Excel rows
2. Insert job row       →    5. Tier 1/2: deterministic     
2. Return 202 + job ID            upsert via supplier_code
                             6. Tier 3: batch to Claude  →  7. Match + score
                             8. Write results                   return JSON
                             9. Update flags table
                            10. Mark job complete
                            11. Delete temp file

Browser (HTMX polling)
──────────────────────
Poll /imports/{id}/status every 2s until complete
Render summary + flag review queue
```

---

## Match Logic

Rows from the spreadsheet are processed in three tiers. Claude is invoked
for Tier 2 (name normalization only) and Tier 3 (matching + name normalization).

| Tier | Condition | Action |
|------|-----------|--------|
| 1 | `supplier_code` matches existing material for this supplier | Update `unit_price` + `last_updated` unless `price_source = 'Manual'`. Skip if Manual. |
| 2 | `supplier_code` present but no match in DB | Normalize name via Claude → create new material with `price_source = 'Supplier'` |
| 3 | No `supplier_code` in row | Batch to Claude for name/unit/price fuzzy match + name normalization |

Claude returns a confidence score per row (0.0–1.0):
- **≥ 0.98** → auto-update price (unless Manual)
- **0.75–0.97** → insert into flag queue for owner review
- **< 0.75** → treat as new material, create it

---

## Material Naming Convention

All new materials created by the import pipeline (Tier 2 creates, Tier 3 creates,
and flagged items resolved as "Create New") are normalized to a canonical format
optimized for search-first workflows in the estimate builder.

**Format:** `Category, Dimensions/Weight, Details`

The category (what the item IS) always comes first so the user can type a category
name and immediately see all sizes/variants. Dimensions use trade shorthand.
Species, brand, or grade come last.

| Supplier Raw Name | Canonical Name |
|---|---|
| `2x4x8 Douglas Fir` | `Stud, 2x4x8, Douglas Fir` |
| `2x4 Stud 12ft DF` | `Stud, 2x4x12, Douglas Fir` |
| `2x6-16 #2 SPF` | `Lumber, 2x6x16, #2 SPF` |
| `3/4 CDX Plywood 4x8` | `Plywood, 4x8, 3/4 CDX` |
| `Simpson LUS210 Joist Hanger` | `Joist Hanger, 2x10, Simpson LUS210` |
| `Quikrete 80lb Concrete Mix` | `Concrete Mix, 80lb, Quikrete` |
| `#10 x 3 GRK RSS Screw` | `Screw, #10x3, GRK RSS` |
| `16d Framing Nail 3-1/2 Bright` | `Nail, 16d 3-1/2, Bright Framing` |
| `PL Premium Adhesive 28oz` | `Construction Adhesive, 28oz, PL Premium` |

**Implementation:**
- Tier 2: `NormalizeNames()` sends raw names to Claude for canonical formatting
- Tier 3: `BatchMatch()` returns both match results and `canonical_name` per row
- Tier 1 (price updates): existing material names are not changed
- Both Claude calls share the same naming rules and examples in their prompts

---

## Step 1 — Database Migration

**File:** `migrations/XXXX_add_material_imports.sql`

```sql
CREATE TABLE material_import_jobs (
    id              TEXT PRIMARY KEY,           -- nanoid
    supplier_id     TEXT NOT NULL REFERENCES suppliers(id),
    filename        TEXT NOT NULL,              -- original filename for display
    temp_file_path  TEXT NOT NULL,              -- absolute path, deleted on completion
    status          TEXT NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'processing', 'complete', 'failed')),
    rows_total      INTEGER,                    -- populated after parse
    rows_created    INTEGER,
    rows_updated    INTEGER,
    rows_skipped    INTEGER,                    -- Manual price_source, no change
    rows_flagged    INTEGER,
    rows_failed     INTEGER,
    error_message   TEXT,                       -- populated if status = 'failed'
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ
);

CREATE TABLE material_import_flags (
    id                   TEXT PRIMARY KEY,      -- nanoid
    import_job_id        TEXT NOT NULL REFERENCES material_import_jobs(id) ON DELETE CASCADE,
    supplier_id          TEXT NOT NULL REFERENCES suppliers(id),
    -- Raw values from spreadsheet
    raw_name             TEXT NOT NULL,
    raw_unit             TEXT,
    raw_price            REAL,
    raw_supplier_code    TEXT,
    -- Claude's assessment
    candidate_material_id TEXT REFERENCES materials(id) ON DELETE SET NULL,
    confidence           REAL,                  -- 0.0–1.0
    claude_reasoning     TEXT,
    -- Owner resolution
    resolution           TEXT CHECK (resolution IN ('accepted', 'rejected', 'created')),
    resolved_at          TIMESTAMPTZ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for the worker's polling query
CREATE INDEX idx_import_jobs_status ON material_import_jobs(status)
    WHERE status = 'pending';
```

**Test:** Insert a job row with status `pending`, confirm index is used in
`EXPLAIN` for the poll query. Insert a flag row, confirm cascade delete when
job is deleted.

---

## Step 2 — sqlc Queries

**File:** `db/query/material_imports.sql`

### Jobs

```sql
-- name: CreateImportJob :one
INSERT INTO material_import_jobs (id, supplier_id, filename, temp_file_path)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: ClaimNextImportJob :one
-- Postgres queue pattern — SKIP LOCKED prevents double-processing
UPDATE material_import_jobs
SET status = 'processing', started_at = NOW()
WHERE id = (
    SELECT id FROM material_import_jobs
    WHERE status = 'pending'
    ORDER BY created_at ASC
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: GetImportJob :one
SELECT * FROM material_import_jobs WHERE id = ?;

-- name: UpdateImportJobProgress :one
UPDATE material_import_jobs
SET rows_total   = ?,
    rows_created = ?,
    rows_updated = ?,
    rows_skipped = ?,
    rows_flagged = ?,
    rows_failed  = ?
WHERE id = ?
RETURNING *;

-- name: CompleteImportJob :one
UPDATE material_import_jobs
SET status = 'complete', completed_at = NOW()
WHERE id = ?
RETURNING *;

-- name: FailImportJob :one
UPDATE material_import_jobs
SET status = 'failed', completed_at = NOW(), error_message = ?
WHERE id = ?
RETURNING *;

-- name: ListImportJobs :many
SELECT j.*, s.name AS supplier_name
FROM material_import_jobs j
JOIN suppliers s ON s.id = j.supplier_id
ORDER BY j.created_at DESC;
```

### Flags

```sql
-- name: CreateImportFlag :one
INSERT INTO material_import_flags (
    id, import_job_id, supplier_id,
    raw_name, raw_unit, raw_price, raw_supplier_code,
    candidate_material_id, confidence, claude_reasoning
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListImportFlags :many
SELECT f.*, m.name AS candidate_name, m.unit_price AS candidate_price
FROM material_import_flags f
LEFT JOIN materials m ON m.id = f.candidate_material_id
WHERE f.import_job_id = ?
ORDER BY f.confidence DESC;

-- name: ListPendingFlags :many
SELECT f.*, m.name AS candidate_name, m.unit_price AS candidate_price
FROM material_import_flags f
LEFT JOIN materials m ON m.id = f.candidate_material_id
WHERE f.supplier_id = ? AND f.resolution IS NULL
ORDER BY f.confidence DESC;

-- name: ResolveFlag :one
UPDATE material_import_flags
SET resolution = ?, resolved_at = NOW()
WHERE id = ?
RETURNING *;
```

Run `sqlc generate`, confirm no errors.

**Test:** Confirm `ClaimNextImportJob` uses `FOR UPDATE SKIP LOCKED` in the
generated query. Spot-check nullable fields come through correctly.

---

## Step 3 — Background Worker Harness

**File:** `worker/worker.go`

Establishes the background worker pattern for the app. A single goroutine polls
the jobs table on an interval, claims one job at a time, and dispatches to a
handler function.

```go
type Worker struct {
    db       *db.Queries
    rawDB    *sql.DB         // for transactions
    interval time.Duration
    handlers map[string]JobHandler
}

type JobHandler func(ctx context.Context, job db.MaterialImportJob) error

func New(db *db.Queries, rawDB *sql.DB, interval time.Duration) *Worker

// Register adds a handler for a job type. For now only "material_import" exists.
func (w *Worker) Register(jobType string, h JobHandler)

// Start begins polling in a goroutine. Respects ctx cancellation for clean shutdown.
func (w *Worker) Start(ctx context.Context)
```

Poll loop behavior:
1. Call `ClaimNextImportJob` — if no rows, sleep `interval` and retry
2. Dispatch to registered handler
3. On handler error → call `FailImportJob` with error message
4. On panic → recover, call `FailImportJob`, continue polling (don't crash the app)

Wire the worker into `main.go`:

```go
w := worker.New(queries, rawDB, 5*time.Second)
w.Register("material_import", importer.Handle)

ctx, cancel := context.WithCancel(context.Background())
defer cancel()
go w.Start(ctx)
```

Shutdown: use `signal.NotifyContext` so the worker drains cleanly on SIGTERM.

**Test:** Insert a `pending` job row manually, confirm the worker picks it up
within one poll interval, confirm status transitions to `processing` then
`complete`. Kill and restart the app mid-job, confirm the stuck `processing`
job is not re-claimed (it stays processing — a separate recovery step can be
added later if needed).

---

## Step 4 — Excel Parser

**File:** `worker/importer/parse.go`

Responsible only for reading the Excel file and returning a slice of raw rows.
Decoupled from matching logic so the input source can change later (web scraping
returns the same row struct).

```go
type RawRow struct {
    Name         string
    Unit         string
    Price        float64
    SupplierCode string  // empty string if not present
    RowIndex     int     // for error reporting
}

func ParseExcel(path string) ([]RawRow, error)
```

Uses `excelize` or `xlsx` library to read the file. Skips header rows,
skips blank rows, normalizes whitespace. Does not validate business logic —
that happens in the matcher.

### Unit normalization

Suppliers use inconsistent unit strings. `ParseExcel` applies a normalization
map before returning rows so downstream code works with consistent values:

```go
var unitNormalizer = map[string]string{
    "each":        "ea",
    "linear foot": "lf",
    "linear feet": "lf",
    "lineal foot": "lf",
    "lineal feet": "lf",
    "square foot": "sf",
    "square feet": "sf",
    "board foot":  "bf",
    "board feet":  "bf",
    "cubic yard":  "cy",
    "cubic yards": "cy",
    "hour":        "hr",
    "hours":       "hr",
    "gallon":      "gal",
    "gallons":     "gal",
    "pound":       "lb",
    "pounds":      "lb",
}

func normalizeUnit(raw string) string {
    key := strings.ToLower(strings.TrimSpace(raw))
    if normalized, ok := unitNormalizer[key]; ok {
        return normalized
    }
    return key  // pass through as-is — Claude uses it as a weak signal
}
```

Unrecognized units are lowercased and passed through unchanged. Claude treats
the unit as a supporting signal for Tier 3 matching, so an unrecognized unit
degrades confidence slightly rather than causing an error. The normalizer map
should be expanded as real supplier spreadsheets are encountered.

**Test:** Parse a sample supplier spreadsheet. Confirm header row is skipped.
Confirm rows with missing name are skipped. Confirm `SupplierCode` is empty
string (not a zero value that looks like a code) when the column is absent.
Confirm "each" → "ea", "linear feet" → "lf", and an unknown unit passes
through lowercased.

---

## Step 5 — Claude Matcher

**File:** `worker/importer/matcher.go`

Handles Tier 3 matching — rows without a `supplier_code`. Takes all Tier 3
rows and the full list of existing materials for this supplier, sends one
batched API call to Claude, returns scored match results.

### Prompt design

System prompt:
```
You are a materials matching assistant for a construction quoting application.
You will be given a list of existing materials and a list of unmatched rows
from a supplier spreadsheet. For each unmatched row, identify the best matching
existing material and return a confidence score from 0.0 to 1.0.

Rules:
- Use name similarity as the primary signal
- Use unit as a supporting signal (ea, lf, sf, etc.)
- Use price proximity as a weak supporting signal only
- Return null for candidate_id if no reasonable match exists (confidence < 0.50)
- Return only valid JSON, no explanation outside the JSON structure

Response format:
{
  "matches": [
    {
      "row_index": 0,
      "candidate_id": "abc123" | null,
      "confidence": 0.94,
      "reasoning": "Name is nearly identical, units match"
    }
  ]
}
```

User message contains:
```json
{
  "existing_materials": [
    { "id": "abc123", "name": "2x4x8 Douglas Fir", "unit": "ea", "price": 8.50 }
  ],
  "unmatched_rows": [
    { "row_index": 4, "name": "2x4 Stud 8ft DF", "unit": "each", "price": 8.75 }
  ]
}
```

```go
type MatchResult struct {
    RowIndex    int
    CandidateID string   // empty if no match
    Confidence  float64
    Reasoning   string
}

func BatchMatch(ctx context.Context, rows []RawRow, existing []db.Material) ([]MatchResult, error)
```

**Test:** Send a batch with obvious matches, near matches, and clear non-matches.
Confirm Claude returns valid JSON. Confirm response parsing handles null
`candidate_id` correctly. Confirm behavior when all rows are unmatched
(empty existing materials list).

---

## Step 6 — Import Job Handler

**File:** `worker/importer/importer.go`

Orchestrates the full import pipeline for a single job. This is what the worker
dispatches to.

```go
func Handle(ctx context.Context, job db.MaterialImportJob) error
```

Pipeline:

1. **Parse** — call `ParseExcel(job.TempFilePath)`, get `[]RawRow`
2. **Update total** — call `UpdateImportJobProgress` with `rows_total`
3. **Load existing** — call `ListMaterialsBySupplier` for this supplier
4. **Split tiers**
   - Tier 1: row has `supplier_code` AND matches an existing material
   - Tier 2: row has `supplier_code` AND no match
   - Tier 3: row has no `supplier_code`
5. **Process Tier 1** — for each match:
   - If `price_source = 'Manual'` → increment `rows_skipped`, continue
   - Else → update `unit_price` + `last_updated`, increment `rows_updated`
6. **Process Tier 2** — create new material, `price_source = 'Supplier'`,
   increment `rows_created`
7. **Process Tier 3** — call `BatchMatch`, then for each result:
   - Confidence ≥ 0.98 → update price (unless Manual → skip), increment
     `rows_updated` or `rows_skipped`
   - Confidence 0.75–0.97 → insert `material_import_flag`, increment
     `rows_flagged`
   - Confidence < 0.75 → create new material, increment `rows_created`
8. **Update progress** — call `UpdateImportJobProgress` with final counts
9. **Complete job** — call `CompleteImportJob`
10. **Delete temp file** — `os.Remove(job.TempFilePath)`

All DB writes in steps 5–7 run inside a single transaction. If anything fails,
the transaction rolls back and `FailImportJob` is called with the error.

**Test:** Run a full import with a known spreadsheet. Confirm counts match
expected create/update/skip/flag breakdown. Confirm temp file is deleted on
success and on failure. Confirm a Manual-priced material is not updated.

---

## Step 7 — HTTP Handler

**File:** `handler/import/import.go`

Thin handler — validates input, writes the file, enqueues the job, returns
immediately.

| Method | Route | Action |
|--------|-------|--------|
| GET | `/imports` | List all import jobs |
| POST | `/imports` | Accept upload, enqueue job, redirect to job status page |
| GET | `/imports/{id}` | Job status page (HTMX polling target) |
| GET | `/imports/{id}/status` | Polling fragment — returns summary or spinner |
| GET | `/imports/{id}/flags` | Flag review queue for this import |
| POST | `/imports/{id}/flags/{fid}/resolve` | Owner resolves a flag |

### Upload handler (`POST /imports`)

1. Validate `supplier_id` exists
2. Read multipart file upload
3. Write to `os.TempDir()/skalkaho-import-{nanoid}.xlsx`
4. Insert `material_import_job` with status `pending`
5. Redirect to `/imports/{id}`

### Status polling (`GET /imports/{id}/status`)

Called by HTMX every 2 seconds. Returns:
- If `pending` or `processing` → spinner fragment with current counts if available
- If `complete` → summary fragment (rows created/updated/skipped/flagged) +
  flag review section if `rows_flagged > 0`
- If `failed` → error message fragment

HTMX polling stops automatically when the response includes
`HX-Trigger: importComplete` header, which the template listens for to swap
in the final state and stop polling.

### Flag resolution (`POST /imports/{id}/flags/{fid}/resolve`)

Accepts `resolution=accepted|rejected|created` in the form body.
- `accepted` → update the candidate material's price, mark flag resolved
- `rejected` → mark flag resolved, no DB change to material
- `created` → create new material from raw row values, mark flag resolved

**Test:** Upload a file, confirm job appears in DB with status `pending` within
one request cycle. Poll status endpoint, confirm it returns spinner then
eventually summary. Resolve a flag as `accepted`, confirm material price updated.

---

## Step 8 — Templates

**Directory:** `templates/imports/`

### `index.html`
List of past import jobs: supplier name, filename, date, status, row counts.
Link to each job's status page.

### `show.html`
Job status page. Contains the polling div:

```html
<div id="import-status"
     hx-get="/imports/{{ .Job.ID }}/status"
     hx-trigger="every 2s"
     hx-swap="outerHTML">
  <!-- initial spinner -->
</div>
```

### `_status_pending.html` (partial)
Spinner + "Processing your spreadsheet…" + any available row counts.

### `_status_complete.html` (partial)
Summary card: created / updated / skipped / flagged counts.
If flagged > 0, renders the flag review table inline.

### `_status_failed.html` (partial)
Error message + option to retry (link back to upload form for this supplier).

### `_flag_row.html` (partial)
Single flag row for the review queue:
- Raw name, unit, price from spreadsheet
- Candidate material name + current price (if matched)
- Claude's confidence score + reasoning
- Accept / Reject / Create New buttons (HTMX, swap row out on resolution)

### `form.html`
Upload form: supplier dropdown + file input. Posted to `POST /imports`.

**Test:** Upload a file with flagged rows, confirm review queue renders. Resolve
each flag type (accept/reject/create), confirm row swaps out and material DB
state is correct.

---

## Step 9 — Router Wiring

**File:** `router/router.go`

```go
imp := importhandler.NewHandler(queries, rawDB)

r.HandleFunc("GET /imports",                           imp.Index)
r.HandleFunc("POST /imports",                          imp.Create)
r.HandleFunc("GET /imports/{id}",                      imp.Show)
r.HandleFunc("GET /imports/{id}/status",               imp.Status)
r.HandleFunc("GET /imports/{id}/flags",                imp.Flags)
r.HandleFunc("POST /imports/{id}/flags/{fid}/resolve", imp.ResolveFlag)
```

Add nav link. Consider placing under a "Data" or "Catalog" section alongside
materials and suppliers.

**Test:** Full end-to-end — upload spreadsheet, watch polling resolve to
complete, verify DB state, resolve all flags, confirm import history shows
correct counts.

---

## Commit Sequence

| # | Commit message |
|---|----------------|
| 1 | `db: add material_import_jobs and material_import_flags migration` |
| 2 | `db: add sqlc queries for import jobs and flags` |
| 3 | `worker: add background worker harness with poll loop and clean shutdown` |
| 4 | `worker/importer: add Excel parser (ParseExcel)` |
| 5 | `worker/importer: add Claude batch matcher (BatchMatch)` |
| 6 | `worker/importer: add import job handler orchestrating full pipeline` |
| 7 | `handler: add import HTTP handler, upload, status polling, flag resolution` |
| 8 | `templates: add import list, status, and flag review templates` |
| 9 | `router: wire import routes, add nav link` |

---

## Future Considerations

- **Web scraping input:** Replace `ParseExcel` with a `ParseScrape` function
  returning the same `[]RawRow` struct. The rest of the pipeline is unchanged.
- **Stuck job recovery:** A job stuck in `processing` after a server restart
  can be recovered by a startup check that resets `processing → pending` for
  jobs with no `completed_at` older than a threshold.
- **Multiple workers:** `FOR UPDATE SKIP LOCKED` already supports this — spin
  up additional worker goroutines or processes safely.
