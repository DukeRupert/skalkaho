-- +goose Up
-- Rename item templates: put descriptive name before size/dimensions
-- Makes autocomplete search more intuitive (search "stud" not "2x4x92")

-- Studs (pre-cut wall-height lengths)
UPDATE item_templates SET name = 'Stud ' || name
WHERE name IN ('2x4x92 5/8"', '2x4x104 5/8"', '2x4x116"', '2x6x92 5/8"', '2x6x104 5/8"', '2x6x116"')
AND org_id IS NULL;

-- Plain dimensional lumber: add "Lumber" prefix
UPDATE item_templates SET name = 'Lumber ' || name
WHERE name IN (
    '2x4x8', '2x4x10', '2x4x12', '2x4x14', '2x4x16', '2x4x20',
    '2x6x8', '2x6x10', '2x6x12', '2x6x14', '2x6x16', '2x6x20',
    '2x8x8', '2x8x10', '2x8x12', '2x8x14', '2x8x16', '2x8x20',
    '2x10x8', '2x10x10', '2x10x12', '2x10x14', '2x10x16', '2x10x20',
    '2x12x8', '2x12x10', '2x12x12', '2x12x14', '2x12x16', '2x12x20'
)
AND org_id IS NULL;

-- Pressure Treated: "NxNxN PT" → "PT NxNxN"
UPDATE item_templates SET name = 'PT ' || SUBSTR(name, 1, LENGTH(name) - 3)
WHERE name LIKE '% PT' AND category = 'Lumber' AND org_id IS NULL;

-- Douglas Fir: "NxNxN DF" → "DF NxNxN"
UPDATE item_templates SET name = 'DF ' || SUBSTR(name, 1, LENGTH(name) - 3)
WHERE name LIKE '% DF' AND category = 'Lumber' AND org_id IS NULL;

-- T&G Pine: "NxNxN T&G Pine" → "T&G Pine NxNxN"
UPDATE item_templates SET name = 'T&G Pine ' || SUBSTR(name, 1, LENGTH(name) - 10)
WHERE name LIKE '% T&G Pine' AND category = 'Lumber' AND org_id IS NULL;

-- LVL: "dims LVL" → "LVL dims"
UPDATE item_templates SET name = 'LVL ' || SUBSTR(name, 1, LENGTH(name) - 4)
WHERE name LIKE '% LVL' AND category = 'Lumber' AND org_id IS NULL;

-- CDX Plywood: "N/M CDX" → "CDX N/M"
UPDATE item_templates SET name = 'CDX ' || SUBSTR(name, 1, LENGTH(name) - 4)
WHERE name LIKE '% CDX' AND category = 'Lumber' AND org_id IS NULL;

-- Hangers: "NxN Hangers" → "Hangers NxN"
UPDATE item_templates SET name = 'Hangers ' || SUBSTR(name, 1, LENGTH(name) - 8)
WHERE name LIKE '% Hangers' AND name NOT LIKE '% BCI Hangers'
AND category = 'Lumber' AND org_id IS NULL;

-- BCI Hangers: "N N/N BCI Hangers" → "BCI Hangers N N/N"
UPDATE item_templates SET name = 'BCI Hangers ' || SUBSTR(name, 1, LENGTH(name) - 12)
WHERE name LIKE '% BCI Hangers' AND category = 'Lumber' AND org_id IS NULL;

-- Knife Post Base: "NxN Knife Post Base" → "Knife Post Base NxN"
UPDATE item_templates SET name = 'Knife Post Base ' || SUBSTR(name, 1, LENGTH(name) - 16)
WHERE name LIKE '% Knife Post Base' AND category = 'Lumber' AND org_id IS NULL;

-- Individual lumber items
UPDATE item_templates SET name = 'Truss Screw 6" 50/Box' WHERE name = '6" Truss Screw 50/Box' AND org_id IS NULL;
UPDATE item_templates SET name = 'Titan anchor bolt 1/2"x5"' WHERE name = '1/2"x5" Titan anchor bolt' AND org_id IS NULL;
UPDATE item_templates SET name = 'Advantech T&G 3/4"' WHERE name = '3/4" Advantech T&G' AND org_id IS NULL;
UPDATE item_templates SET name = 'Rim Joist 11 7/8 1''' WHERE name = '11 7/8 Rim Joist 1''' AND org_id IS NULL;
UPDATE item_templates SET name = 'House Wrap 9''x150''' WHERE name = '9''x150'' House Wrap' AND org_id IS NULL;
UPDATE item_templates SET name = 'Tyvek 9''x100''' WHERE name = '9''x100'' Tyvek' AND org_id IS NULL;
UPDATE item_templates SET name = 'Sill Seal 5 1/2"x50''' WHERE name = '5 1/2"x50'' Sill Seal' AND org_id IS NULL;
UPDATE item_templates SET name = 'Concrete Sackrete 80lb' WHERE name = '80lb Concrete Sackrete' AND org_id IS NULL;

-- Ply beams: move treatment type to front
UPDATE item_templates SET name = 'Full Treat ' || REPLACE(name, ' Full Treat', '')
WHERE name LIKE '% Full Treat' AND category = 'Insulation' AND org_id IS NULL;

UPDATE item_templates SET name = 'Treated Bottom ' || REPLACE(name, ' Treated Bottom', '')
WHERE name LIKE '% Treated Bottom' AND category = 'Insulation' AND org_id IS NULL;

-- Metal roofing/siding with gauge prefix
UPDATE item_templates SET name = 'Roof Standing Seam 24ga with trim (Elite Metal Supplies)' WHERE name = '24 ga Roof Standing Seam with trim (Elite Metal Supplies)' AND org_id IS NULL;
UPDATE item_templates SET name = 'Roof Standing Seam 26ga with trim (Elite Metal Supplies)' WHERE name = '26 ga Roof Standing Seam with trim (Elite Metal Supplies)' AND org_id IS NULL;
UPDATE item_templates SET name = 'Steel Siding 26ga with trim (Elite Metal Supplies)' WHERE name = '26ga Steel Siding with trim sqft (Elite Metal Supplies)' AND org_id IS NULL;
UPDATE item_templates SET name = 'Steel Board & Batt 10" per sqft' WHERE name = '10" Steel Board & Batt per sqft' AND org_id IS NULL;

-- Concrete/rebar/misc with size prefix
UPDATE item_templates SET name = 'Sonotube 12" per foot' WHERE name = '12" Sonotube per foot' AND org_id IS NULL;
UPDATE item_templates SET name = 'PIER W/ 4X4 BRACKET 12X12X8' WHERE name = '12X12X8 PIER W/ 4X4 BRACKET' AND org_id IS NULL;
UPDATE item_templates SET name = '#4 LITE BAR 7/16"X20''' WHERE name = '7/16"X20'' #4 LITE BAR' AND org_id IS NULL;
UPDATE item_templates SET name = 'Steel Rebar 1/2"x20''' WHERE name = '1/2"x20'' Steel Rebar' AND org_id IS NULL;
UPDATE item_templates SET name = 'STEP FLASHING BLACK 4"X4"X14"' WHERE name = '4"X4"X14" STEP FLASHING BLACK' AND org_id IS NULL;
UPDATE item_templates SET name = 'PRO-TAC JOIST TAPE 1.625"X65''' WHERE name = '1.625"X65'' PRO-TAC JOIST TAPE' AND org_id IS NULL;

-- +goose Down
-- Reverse all renames

-- Studs: remove "Stud " prefix
UPDATE item_templates SET name = SUBSTR(name, 6)
WHERE name LIKE 'Stud %' AND category = 'Lumber' AND org_id IS NULL;

-- Plain lumber: remove "Lumber " prefix
UPDATE item_templates SET name = SUBSTR(name, 8)
WHERE name LIKE 'Lumber %' AND category = 'Lumber' AND org_id IS NULL;

-- PT: move back to suffix
UPDATE item_templates SET name = SUBSTR(name, 4) || ' PT'
WHERE name LIKE 'PT %' AND category = 'Lumber' AND org_id IS NULL;

-- DF: move back to suffix
UPDATE item_templates SET name = SUBSTR(name, 4) || ' DF'
WHERE name LIKE 'DF %' AND category = 'Lumber' AND org_id IS NULL;

-- T&G Pine: move back to suffix
UPDATE item_templates SET name = SUBSTR(name, 10) || ' T&G Pine'
WHERE name LIKE 'T&G Pine %' AND category = 'Lumber' AND org_id IS NULL;

-- LVL: move back to suffix
UPDATE item_templates SET name = SUBSTR(name, 5) || ' LVL'
WHERE name LIKE 'LVL %' AND category = 'Lumber' AND org_id IS NULL;

-- CDX: move back to suffix
UPDATE item_templates SET name = SUBSTR(name, 5) || ' CDX'
WHERE name LIKE 'CDX %' AND category = 'Lumber' AND org_id IS NULL;

-- Hangers: move size back to front
UPDATE item_templates SET name = SUBSTR(name, 9) || ' Hangers'
WHERE name LIKE 'Hangers %' AND category = 'Lumber' AND org_id IS NULL;

-- BCI Hangers: move size back to front
UPDATE item_templates SET name = SUBSTR(name, 13) || ' BCI Hangers'
WHERE name LIKE 'BCI Hangers %' AND category = 'Lumber' AND org_id IS NULL;

-- Knife Post Base: move size back to front
UPDATE item_templates SET name = SUBSTR(name, 17) || ' Knife Post Base'
WHERE name LIKE 'Knife Post Base %' AND category = 'Lumber' AND org_id IS NULL;

-- Individual lumber items
UPDATE item_templates SET name = '6" Truss Screw 50/Box' WHERE name = 'Truss Screw 6" 50/Box' AND org_id IS NULL;
UPDATE item_templates SET name = '1/2"x5" Titan anchor bolt' WHERE name = 'Titan anchor bolt 1/2"x5"' AND org_id IS NULL;
UPDATE item_templates SET name = '3/4" Advantech T&G' WHERE name = 'Advantech T&G 3/4"' AND org_id IS NULL;
UPDATE item_templates SET name = '11 7/8 Rim Joist 1''' WHERE name = 'Rim Joist 11 7/8 1''' AND org_id IS NULL;
UPDATE item_templates SET name = '9''x150'' House Wrap' WHERE name = 'House Wrap 9''x150''' AND org_id IS NULL;
UPDATE item_templates SET name = '9''x100'' Tyvek' WHERE name = 'Tyvek 9''x100''' AND org_id IS NULL;
UPDATE item_templates SET name = '5 1/2"x50'' Sill Seal' WHERE name = 'Sill Seal 5 1/2"x50''' AND org_id IS NULL;
UPDATE item_templates SET name = '80lb Concrete Sackrete' WHERE name = 'Concrete Sackrete 80lb' AND org_id IS NULL;

-- Ply beams: move treatment back to end
UPDATE item_templates SET name = REPLACE(name, 'Full Treat ', '') || ' Full Treat'
WHERE name LIKE 'Full Treat %' AND category = 'Insulation' AND org_id IS NULL;

UPDATE item_templates SET name = REPLACE(name, 'Treated Bottom ', '') || ' Treated Bottom'
WHERE name LIKE 'Treated Bottom %' AND category = 'Insulation' AND org_id IS NULL;

-- Metal roofing/siding
UPDATE item_templates SET name = '24 ga Roof Standing Seam with trim (Elite Metal Supplies)' WHERE name = 'Roof Standing Seam 24ga with trim (Elite Metal Supplies)' AND org_id IS NULL;
UPDATE item_templates SET name = '26 ga Roof Standing Seam with trim (Elite Metal Supplies)' WHERE name = 'Roof Standing Seam 26ga with trim (Elite Metal Supplies)' AND org_id IS NULL;
UPDATE item_templates SET name = '26ga Steel Siding with trim sqft (Elite Metal Supplies)' WHERE name = 'Steel Siding 26ga with trim (Elite Metal Supplies)' AND org_id IS NULL;
UPDATE item_templates SET name = '10" Steel Board & Batt per sqft' WHERE name = 'Steel Board & Batt 10" per sqft' AND org_id IS NULL;

-- Concrete/rebar/misc
UPDATE item_templates SET name = '12" Sonotube per foot' WHERE name = 'Sonotube 12" per foot' AND org_id IS NULL;
UPDATE item_templates SET name = '12X12X8 PIER W/ 4X4 BRACKET' WHERE name = 'PIER W/ 4X4 BRACKET 12X12X8' AND org_id IS NULL;
UPDATE item_templates SET name = '7/16"X20'' #4 LITE BAR' WHERE name = '#4 LITE BAR 7/16"X20''' AND org_id IS NULL;
UPDATE item_templates SET name = '1/2"x20'' Steel Rebar' WHERE name = 'Steel Rebar 1/2"x20''' AND org_id IS NULL;
UPDATE item_templates SET name = '4"X4"X14" STEP FLASHING BLACK' WHERE name = 'STEP FLASHING BLACK 4"X4"X14"' AND org_id IS NULL;
UPDATE item_templates SET name = '1.625"X65'' PRO-TAC JOIST TAPE' WHERE name = 'PRO-TAC JOIST TAPE 1.625"X65''' AND org_id IS NULL;
