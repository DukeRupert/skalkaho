-- +goose Up

CREATE TABLE supply_lists (
    id                  TEXT PRIMARY KEY,
    quote_id            TEXT NOT NULL REFERENCES quotes(id) ON DELETE RESTRICT,
    supplier_id         TEXT REFERENCES suppliers(id) ON DELETE RESTRICT,
    label               TEXT NOT NULL,
    estimated_delivery  DATE,
    actual_delivery     DATE,
    generated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    notes               TEXT
);

CREATE INDEX idx_supply_lists_quote_id ON supply_lists(quote_id);
CREATE INDEX idx_supply_lists_supplier ON supply_lists(quote_id, supplier_id);

CREATE TABLE supply_list_items (
    id               TEXT PRIMARY KEY,
    supply_list_id   TEXT NOT NULL REFERENCES supply_lists(id) ON DELETE CASCADE,
    material_id      TEXT,
    supplier_id      TEXT,
    supplier_name    TEXT,
    material_name    TEXT NOT NULL,
    supplier_code    TEXT,
    quantity         REAL NOT NULL,
    unit             TEXT NOT NULL,
    unit_price       REAL NOT NULL DEFAULT 0,
    extended_price   REAL NOT NULL DEFAULT 0,
    sort_order       INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_supply_list_items_list ON supply_list_items(supply_list_id);

-- +goose Down

DROP TABLE IF EXISTS supply_list_items;
DROP TABLE IF EXISTS supply_lists;
