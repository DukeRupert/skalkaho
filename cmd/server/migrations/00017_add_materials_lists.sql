-- +goose Up
CREATE TABLE materials_lists (
    id              TEXT PRIMARY KEY,
    quote_id        TEXT NOT NULL REFERENCES quotes(id) ON DELETE RESTRICT,
    section_id      TEXT REFERENCES sections(id) ON DELETE RESTRICT,
    supply_list_id  TEXT REFERENCES supply_lists(id) ON DELETE SET NULL,
    label           TEXT NOT NULL,
    notes           TEXT,
    generated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_materials_lists_quote ON materials_lists(quote_id);
CREATE INDEX idx_materials_lists_section ON materials_lists(quote_id, section_id);

-- +goose Down
DROP TABLE IF EXISTS materials_lists;
