-- +goose Up
CREATE TABLE material_import_jobs (
    id              TEXT PRIMARY KEY,
    supplier_id     TEXT NOT NULL REFERENCES suppliers(id),
    filename        TEXT NOT NULL,
    temp_file_path  TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'pending'
                        CHECK (status IN ('pending', 'processing', 'complete', 'failed')),
    rows_total      INTEGER,
    rows_created    INTEGER,
    rows_updated    INTEGER,
    rows_skipped    INTEGER,
    rows_flagged    INTEGER,
    rows_failed     INTEGER,
    error_message   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ
);

CREATE TABLE material_import_flags (
    id                    TEXT PRIMARY KEY,
    import_job_id         TEXT NOT NULL REFERENCES material_import_jobs(id) ON DELETE CASCADE,
    supplier_id           TEXT NOT NULL REFERENCES suppliers(id),
    raw_name              TEXT NOT NULL,
    raw_unit              TEXT,
    raw_price             REAL,
    raw_supplier_code     TEXT,
    candidate_material_id TEXT REFERENCES materials(id) ON DELETE SET NULL,
    confidence            REAL,
    claude_reasoning      TEXT,
    resolution            TEXT CHECK (resolution IN ('accepted', 'rejected', 'created')),
    resolved_at           TIMESTAMPTZ,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_import_jobs_status ON material_import_jobs(status)
    WHERE status = 'pending';

-- +goose Down
DROP TABLE IF EXISTS material_import_flags;
DROP TABLE IF EXISTS material_import_jobs;
