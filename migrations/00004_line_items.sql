-- +goose Up
CREATE TABLE line_items (
    id TEXT PRIMARY KEY,
    category_id TEXT NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    type TEXT NOT NULL CHECK (type IN ('material', 'labor')),
    name TEXT NOT NULL,
    description TEXT,
    quantity REAL NOT NULL,
    unit TEXT NOT NULL,
    unit_price REAL NOT NULL,
    surcharge_percent REAL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_line_items_category ON line_items(category_id);

-- +goose Down
DROP INDEX IF EXISTS idx_line_items_category;
DROP TABLE IF EXISTS line_items;
