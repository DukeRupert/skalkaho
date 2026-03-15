-- +goose Up

-- Lookup table for construction trades
CREATE TABLE trades (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    sort_order INTEGER NOT NULL DEFAULT 0
);

-- Seed common trades
INSERT INTO trades (id, name, sort_order) VALUES
    ('trade_concrete',    'Concrete',           1),
    ('trade_drywall',     'Drywall',            2),
    ('trade_electrical',  'Electrical',         3),
    ('trade_excavation',  'Excavation',         4),
    ('trade_flooring',    'Flooring',           5),
    ('trade_framing',     'Framing',            6),
    ('trade_hvac',        'HVAC',               7),
    ('trade_insulation',  'Insulation',         8),
    ('trade_masonry',     'Masonry',            9),
    ('trade_painting',    'Painting',          10),
    ('trade_plumbing',    'Plumbing',          11),
    ('trade_roofing',     'Roofing',           12),
    ('trade_siding',      'Siding',            13),
    ('trade_steel',       'Steel / Structural', 14),
    ('trade_tile',        'Tile',              15),
    ('trade_windows',     'Windows & Doors',   16);

-- Subcontractor directory
CREATE TABLE subcontractors (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    company     TEXT,
    phone       TEXT,
    email       TEXT,
    address     TEXT,
    notes       TEXT,
    is_favorite BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Join table: subcontractor <-> trades (position 0 = primary trade)
CREATE TABLE subcontractor_trades (
    subcontractor_id TEXT NOT NULL REFERENCES subcontractors(id) ON DELETE CASCADE,
    trade_id         TEXT NOT NULL REFERENCES trades(id) ON DELETE CASCADE,
    position         INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (subcontractor_id, trade_id)
);

CREATE INDEX idx_subcontractor_trades_trade ON subcontractor_trades(trade_id);

-- +goose Down
DROP TABLE IF EXISTS subcontractor_trades;
DROP TABLE IF EXISTS subcontractors;
DROP TABLE IF EXISTS trades;
