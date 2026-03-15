-- +goose Up
ALTER TABLE line_items ADD COLUMN subcontractor_id TEXT REFERENCES subcontractors(id) ON DELETE SET NULL;
CREATE INDEX idx_line_items_subcontractor_id ON line_items(subcontractor_id);

-- +goose Down
DROP INDEX IF EXISTS idx_line_items_subcontractor_id;
ALTER TABLE line_items DROP COLUMN IF EXISTS subcontractor_id;
