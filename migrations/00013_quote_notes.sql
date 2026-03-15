-- +goose Up
ALTER TABLE quotes ADD COLUMN notes TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE quotes DROP COLUMN notes;
