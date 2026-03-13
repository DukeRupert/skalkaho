-- +goose Up
CREATE TABLE suppliers (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE materials (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    supplier_id TEXT NOT NULL REFERENCES suppliers(id) ON DELETE CASCADE,
    unit_price REAL NOT NULL DEFAULT 0,
    unit TEXT NOT NULL DEFAULT 'ea',
    supplier_code TEXT,
    price_source TEXT NOT NULL DEFAULT 'Manual' CHECK (price_source IN ('Supplier', 'Manual')),
    last_updated TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_materials_supplier_id ON materials(supplier_id);

-- +goose Down
DROP TABLE IF EXISTS materials;
DROP TABLE IF EXISTS suppliers;
