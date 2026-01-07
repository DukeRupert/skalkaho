-- +goose Up
ALTER TABLE line_items ADD COLUMN tag TEXT;

-- +goose Down
ALTER TABLE line_items DROP COLUMN tag;
