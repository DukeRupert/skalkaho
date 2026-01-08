-- +goose Up
-- Estimates: client-facing snapshots of quotes
CREATE TABLE estimates (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    version INTEGER NOT NULL DEFAULT 1,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'accepted', 'rejected')),
    grand_total REAL NOT NULL DEFAULT 0,
    notes TEXT,
    sent_at TEXT,
    responded_at TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Estimate categories: tier 1 and tier 2 snapshots with descriptions
CREATE TABLE estimate_categories (
    id TEXT PRIMARY KEY,
    estimate_id TEXT NOT NULL REFERENCES estimates(id) ON DELETE CASCADE,
    category_id TEXT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    parent_category_id TEXT REFERENCES categories(id) ON DELETE CASCADE,
    tier INTEGER NOT NULL CHECK (tier IN (1, 2)),
    name TEXT NOT NULL,
    description TEXT,
    total REAL NOT NULL DEFAULT 0,
    sort_order INTEGER NOT NULL DEFAULT 0
);

-- Unique constraint: one version number per job
CREATE UNIQUE INDEX idx_estimates_job_version ON estimates(job_id, version);
CREATE INDEX idx_estimates_job ON estimates(job_id);
CREATE INDEX idx_estimates_status ON estimates(status);
CREATE INDEX idx_estimate_categories_estimate ON estimate_categories(estimate_id);

-- +goose Down
DROP INDEX IF EXISTS idx_estimate_categories_estimate;
DROP INDEX IF EXISTS idx_estimates_status;
DROP INDEX IF EXISTS idx_estimates_job;
DROP INDEX IF EXISTS idx_estimates_job_version;
DROP TABLE IF EXISTS estimate_categories;
DROP TABLE IF EXISTS estimates;
