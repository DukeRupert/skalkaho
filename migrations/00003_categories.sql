-- +goose Up
CREATE TABLE categories (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    parent_id TEXT REFERENCES categories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    surcharge_percent REAL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_categories_job ON categories(job_id);
CREATE INDEX idx_categories_parent ON categories(parent_id);

-- +goose Down
DROP INDEX IF EXISTS idx_categories_parent;
DROP INDEX IF EXISTS idx_categories_job;
DROP TABLE IF EXISTS categories;
