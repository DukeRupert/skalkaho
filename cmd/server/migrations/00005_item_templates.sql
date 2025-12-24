-- +goose Up
CREATE TABLE item_templates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL CHECK (type IN ('material', 'labor', 'equipment')),
    category TEXT NOT NULL,
    name TEXT NOT NULL,
    default_unit TEXT NOT NULL,
    default_price REAL NOT NULL DEFAULT 0
);

CREATE INDEX idx_item_templates_name ON item_templates(name);
CREATE INDEX idx_item_templates_category ON item_templates(category);

-- Also add equipment as valid line item type
-- +goose StatementBegin
CREATE TABLE line_items_new (
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

INSERT INTO line_items_new SELECT * FROM line_items;
DROP TABLE line_items;
ALTER TABLE line_items_new RENAME TO line_items;
CREATE INDEX idx_line_items_category ON line_items(category_id);
-- +goose StatementEnd

-- +goose Down
DROP INDEX IF EXISTS idx_item_templates_name;
DROP INDEX IF EXISTS idx_item_templates_category;
DROP TABLE IF EXISTS item_templates;

-- Revert line_items to original constraint
CREATE TABLE line_items_old (
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

INSERT INTO line_items_old SELECT * FROM line_items WHERE type IN ('material', 'labor');
DROP TABLE line_items;
ALTER TABLE line_items_old RENAME TO line_items;
CREATE INDEX idx_line_items_category ON line_items(category_id);
