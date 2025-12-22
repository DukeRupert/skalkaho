-- +goose Up
CREATE TABLE settings (
    id TEXT PRIMARY KEY DEFAULT 'default',
    default_surcharge_mode TEXT NOT NULL DEFAULT 'stacking'
        CHECK (default_surcharge_mode IN ('stacking', 'override')),
    default_surcharge_percent REAL NOT NULL DEFAULT 0
);

INSERT INTO settings (id) VALUES ('default');

-- +goose Down
DROP TABLE IF EXISTS settings;
