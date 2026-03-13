-- +goose Up
CREATE TABLE rate_categories (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Seed default categories
INSERT INTO rate_categories (id, name, sort_order) VALUES
    ('cat_labor', 'Labor', 0),
    ('cat_equipment', 'Equipment Rentals', 1),
    ('cat_subs', 'Subcontractors', 2),
    ('cat_other', 'Other', 3);

CREATE TABLE rates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category_id TEXT NOT NULL REFERENCES rate_categories(id) ON DELETE CASCADE,
    supplier TEXT,
    rate REAL NOT NULL DEFAULT 0,
    unit TEXT NOT NULL DEFAULT 'hour',
    notes TEXT,
    last_updated TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_rates_category_id ON rates(category_id);

-- +goose Down
DROP TABLE IF EXISTS rates;
DROP TABLE IF EXISTS rate_categories;
