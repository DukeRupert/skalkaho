-- +goose Up
-- Add per-type surcharges to jobs table
ALTER TABLE jobs ADD COLUMN material_surcharge_percent REAL;
ALTER TABLE jobs ADD COLUMN labor_surcharge_percent REAL;
ALTER TABLE jobs ADD COLUMN equipment_surcharge_percent REAL;

-- Add surcharge_percent to custom item types
ALTER TABLE job_item_types ADD COLUMN surcharge_percent REAL;

-- Migrate existing data: copy current surcharge_percent to all type columns
UPDATE jobs SET
    material_surcharge_percent = surcharge_percent,
    labor_surcharge_percent = surcharge_percent,
    equipment_surcharge_percent = surcharge_percent;

-- +goose Down
ALTER TABLE job_item_types DROP COLUMN surcharge_percent;
ALTER TABLE jobs DROP COLUMN equipment_surcharge_percent;
ALTER TABLE jobs DROP COLUMN labor_surcharge_percent;
ALTER TABLE jobs DROP COLUMN material_surcharge_percent;
