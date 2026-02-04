-- +goose Up
-- Create job_item_types table for custom line item types per job
CREATE TABLE job_item_types (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    color TEXT NOT NULL DEFAULT 'slate',
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL DEFAULT (NOW()),
    UNIQUE(job_id, slug)
);

CREATE INDEX idx_job_item_types_job ON job_item_types(job_id);

-- Remove CHECK constraint from line_items.type to allow custom types
-- SQLite requires table recreation to remove constraints
CREATE TABLE line_items_new (
    id TEXT PRIMARY KEY,
    category_id TEXT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    quantity REAL NOT NULL,
    unit TEXT NOT NULL,
    unit_price REAL NOT NULL,
    surcharge_percent REAL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

INSERT INTO line_items_new SELECT * FROM line_items;
DROP TABLE line_items;
ALTER TABLE line_items_new RENAME TO line_items;
CREATE INDEX idx_line_items_category ON line_items(category_id);

-- +goose Down
DROP INDEX IF EXISTS idx_line_items_category;

-- Recreate line_items with CHECK constraint (will fail if custom types exist)
CREATE TABLE line_items_old (
    id TEXT PRIMARY KEY,
    category_id TEXT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN ('material', 'labor', 'equipment')),
    name TEXT NOT NULL,
    description TEXT,
    quantity REAL NOT NULL,
    unit TEXT NOT NULL,
    unit_price REAL NOT NULL,
    surcharge_percent REAL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

INSERT INTO line_items_old SELECT * FROM line_items;
DROP TABLE line_items;
ALTER TABLE line_items_old RENAME TO line_items;
CREATE INDEX idx_line_items_category ON line_items(category_id);

DROP INDEX IF EXISTS idx_job_item_types_job;
DROP TABLE IF EXISTS job_item_types;
