-- +goose Up
ALTER TABLE line_items ADD COLUMN visual_group TEXT;

-- +goose Down
ALTER TABLE line_items DROP COLUMN IF EXISTS visual_group;
