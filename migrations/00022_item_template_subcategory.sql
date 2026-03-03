-- +goose Up
-- Add subcategory column for product family grouping (e.g., "Stud 2x4", "Lumber 2x6")
-- Enables focused batch entry: pick a product family, see only size variants.

ALTER TABLE item_templates ADD COLUMN subcategory TEXT;
CREATE INDEX idx_item_templates_subcategory ON item_templates(subcategory);

-- Populate subcategory from name prefix patterns (Lumber category)

-- Studs
UPDATE item_templates SET subcategory = 'Stud 2x4'
WHERE name LIKE 'Stud 2x4%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'Stud 2x6'
WHERE name LIKE 'Stud 2x6%' AND category = 'Lumber';

-- Dimensional lumber
UPDATE item_templates SET subcategory = 'Lumber 2x4'
WHERE name LIKE 'Lumber 2x4%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'Lumber 2x6'
WHERE name LIKE 'Lumber 2x6%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'Lumber 2x8'
WHERE name LIKE 'Lumber 2x8%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'Lumber 2x10'
WHERE name LIKE 'Lumber 2x10%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'Lumber 2x12'
WHERE name LIKE 'Lumber 2x12%' AND category = 'Lumber';

-- Pressure Treated
UPDATE item_templates SET subcategory = 'PT 2x4'
WHERE name LIKE 'PT 2x4%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'PT 2x6'
WHERE name LIKE 'PT 2x6%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'PT 2x8'
WHERE name LIKE 'PT 2x8%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'PT 2x10'
WHERE name LIKE 'PT 2x10%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'PT 2x12'
WHERE name LIKE 'PT 2x12%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'PT 4x4'
WHERE name LIKE 'PT 4x4%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'PT 6x6'
WHERE name LIKE 'PT 6x6%' AND category = 'Lumber';

-- Douglas Fir
UPDATE item_templates SET subcategory = 'DF 4x4'
WHERE name LIKE 'DF 4x4%' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'DF 6x6'
WHERE name LIKE 'DF 6x6%' AND category = 'Lumber';

-- Glulam
UPDATE item_templates SET subcategory = 'Glulam'
WHERE name LIKE 'GL %' AND category = 'Lumber';

-- LVL
UPDATE item_templates SET subcategory = 'LVL'
WHERE name LIKE 'LVL %' AND category = 'Lumber';

-- CDX
UPDATE item_templates SET subcategory = 'CDX'
WHERE name LIKE 'CDX %' AND category = 'Lumber';

-- T&G Pine
UPDATE item_templates SET subcategory = 'T&G Pine'
WHERE name LIKE 'T&G Pine %' AND category = 'Lumber';

-- Pine 1x boards
UPDATE item_templates SET subcategory = 'Pine 1x'
WHERE name LIKE '1x%Pine' AND category = 'Lumber';

-- Hangers
UPDATE item_templates SET subcategory = 'Hangers'
WHERE name LIKE 'Hangers %' AND category = 'Lumber';

UPDATE item_templates SET subcategory = 'BCI Hangers'
WHERE name LIKE 'BCI Hangers %' AND category = 'Lumber';

-- Knife Post Bases
UPDATE item_templates SET subcategory = 'Knife Post Base'
WHERE name LIKE 'Knife Post Base %' AND category = 'Lumber';

-- Plywood and OSB
UPDATE item_templates SET subcategory = 'Plywood'
WHERE name LIKE '%Plywood%' AND category = 'Lumber' AND subcategory IS NULL;

UPDATE item_templates SET subcategory = 'OSB'
WHERE name LIKE '%OSB%' AND category = 'Lumber' AND subcategory IS NULL;

-- Insulation category

-- Full Treat 3 Ply 2x6
UPDATE item_templates SET subcategory = 'Full Treat 3 Ply 2x6'
WHERE name LIKE 'Full Treat 3 Ply 2x6%' AND category = 'Insulation';

-- Treated Bottom 3 Ply 2x6
UPDATE item_templates SET subcategory = 'Treated Bottom 3 Ply 2x6'
WHERE name LIKE 'Treated Bottom 3 Ply 2x6%' AND category = 'Insulation';

-- Treated Bottom 3 Ply 2x8
UPDATE item_templates SET subcategory = 'Treated Bottom 3 Ply 2x8'
WHERE name LIKE 'Treated Bottom 3 Ply 2x8%' AND category = 'Insulation';

-- Treated Bottom 4 Ply 2x6
UPDATE item_templates SET subcategory = 'Treated Bottom 4 Ply 2x6'
WHERE name LIKE 'Treated Bottom 4 Ply 2x6%' AND category = 'Insulation';

-- Treated Bottom 4 Ply 2x8
UPDATE item_templates SET subcategory = 'Treated Bottom 4 Ply 2x8'
WHERE name LIKE 'Treated Bottom 4 Ply 2x8%' AND category = 'Insulation';

-- Fiberglass Batts
UPDATE item_templates SET subcategory = 'Fiberglass Batts'
WHERE name LIKE 'R% Fiberglass Batts%' AND category = 'Insulation';

-- Triad Screws
UPDATE item_templates SET subcategory = 'Triad Screws'
WHERE name LIKE 'Triad #%' AND category = 'Insulation';

-- Triad Trim/Components (remaining Triad items)
UPDATE item_templates SET subcategory = 'Triad Trim'
WHERE name LIKE 'Triad %' AND category = 'Insulation' AND subcategory IS NULL;

-- Roofing products
UPDATE item_templates SET subcategory = 'Standing Seam'
WHERE name LIKE '%Standing Seam%' AND category = 'Insulation';

UPDATE item_templates SET subcategory = 'GAF Roofing'
WHERE name LIKE 'GAF %' AND category = 'Insulation';

-- Catch-all: singletons get subcategory = name
UPDATE item_templates SET subcategory = name WHERE subcategory IS NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_item_templates_subcategory;
ALTER TABLE item_templates DROP COLUMN IF EXISTS subcategory;
