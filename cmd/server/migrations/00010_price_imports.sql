-- +goose Up
CREATE TABLE price_imports (
    id TEXT PRIMARY KEY,
    filename TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('pending', 'processing', 'ready', 'applied', 'failed')),
    total_rows INTEGER NOT NULL DEFAULT 0,
    matched_rows INTEGER NOT NULL DEFAULT 0,
    error_message TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    applied_at TEXT
);

CREATE TABLE price_import_matches (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    import_id TEXT NOT NULL REFERENCES price_imports(id) ON DELETE CASCADE,
    row_number INTEGER NOT NULL,
    source_name TEXT NOT NULL,
    source_unit TEXT,
    source_price REAL NOT NULL,
    matched_template_id INTEGER REFERENCES item_templates(id),
    confidence REAL NOT NULL DEFAULT 0,
    match_reason TEXT,
    status TEXT NOT NULL CHECK (status IN ('pending', 'approved', 'rejected', 'auto_approved')) DEFAULT 'pending',
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_price_import_matches_import ON price_import_matches(import_id);

-- +goose Down
DROP INDEX IF EXISTS idx_price_import_matches_import;
DROP TABLE IF EXISTS price_import_matches;
DROP TABLE IF EXISTS price_imports;
