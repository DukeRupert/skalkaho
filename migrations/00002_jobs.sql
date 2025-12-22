-- +goose Up
CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    customer_name TEXT,
    surcharge_percent REAL NOT NULL DEFAULT 0,
    surcharge_mode TEXT NOT NULL DEFAULT 'stacking'
        CHECK (surcharge_mode IN ('stacking', 'override')),
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- +goose Down
DROP TABLE IF EXISTS jobs;
