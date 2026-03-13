-- +goose Up
-- Estimate Builder: 4-level hierarchy (Section -> Subcategory -> ComponentGroup -> LineItem)
-- with 5-type markup engine (materials, labor, equipment, subs, other)

-- Global markup defaults stored on jobs table
ALTER TABLE jobs ADD COLUMN materials_markup REAL NOT NULL DEFAULT 20;
ALTER TABLE jobs ADD COLUMN labor_markup REAL NOT NULL DEFAULT 25;
ALTER TABLE jobs ADD COLUMN equipment_markup REAL NOT NULL DEFAULT 15;
ALTER TABLE jobs ADD COLUMN subs_markup REAL NOT NULL DEFAULT 10;
ALTER TABLE jobs ADD COLUMN other_markup REAL NOT NULL DEFAULT 10;

-- Sections: top level of the estimate hierarchy
CREATE TABLE sections (
    id TEXT PRIMARY KEY,
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    job_id TEXT NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_sections_job_id ON sections(job_id);
CREATE INDEX idx_sections_org_id ON sections(org_id);

-- Subcategories: second level, owns markup overrides and enable toggles
CREATE TABLE subcategories (
    id TEXT PRIMARY KEY,
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    section_id TEXT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    lump_sum REAL NOT NULL DEFAULT 0,
    -- Markup overrides (null = inherit global)
    materials_markup REAL,
    labor_markup REAL,
    equipment_markup REAL,
    subs_markup REAL,
    other_markup REAL,
    -- Markup enable toggles
    materials_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    labor_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    equipment_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    subs_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    other_markup_enabled BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX idx_subcategories_section_id ON subcategories(section_id);
CREATE INDEX idx_subcategories_org_id ON subcategories(org_id);

-- Component groups: third level, optional grouping within subcategories
CREATE TABLE component_groups (
    id TEXT PRIMARY KEY,
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    subcategory_id TEXT NOT NULL REFERENCES subcategories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_component_groups_subcategory_id ON component_groups(subcategory_id);
CREATE INDEX idx_component_groups_org_id ON component_groups(org_id);

-- Estimate line items: fourth level, the actual cost entries
CREATE TABLE estimate_line_items (
    id TEXT PRIMARY KEY,
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    subcategory_id TEXT NOT NULL REFERENCES subcategories(id) ON DELETE CASCADE,
    component_group_id TEXT REFERENCES component_groups(id) ON DELETE SET NULL,
    category_type TEXT NOT NULL CHECK (category_type IN ('materials', 'labor', 'equipment', 'subs', 'other')),
    item_name TEXT NOT NULL,
    quantity REAL NOT NULL DEFAULT 1,
    unit TEXT NOT NULL DEFAULT 'ea',
    unit_price REAL NOT NULL DEFAULT 0,
    is_custom BOOLEAN NOT NULL DEFAULT true,
    material_id INTEGER,
    price_override BOOLEAN NOT NULL DEFAULT false,
    description TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_estimate_line_items_subcategory_id ON estimate_line_items(subcategory_id);
CREATE INDEX idx_estimate_line_items_component_group_id ON estimate_line_items(component_group_id);
CREATE INDEX idx_estimate_line_items_org_id ON estimate_line_items(org_id);

-- +goose Down
DROP INDEX IF EXISTS idx_estimate_line_items_org_id;
DROP INDEX IF EXISTS idx_estimate_line_items_component_group_id;
DROP INDEX IF EXISTS idx_estimate_line_items_subcategory_id;
DROP TABLE IF EXISTS estimate_line_items;

DROP INDEX IF EXISTS idx_component_groups_org_id;
DROP INDEX IF EXISTS idx_component_groups_subcategory_id;
DROP TABLE IF EXISTS component_groups;

DROP INDEX IF EXISTS idx_subcategories_org_id;
DROP INDEX IF EXISTS idx_subcategories_section_id;
DROP TABLE IF EXISTS subcategories;

DROP INDEX IF EXISTS idx_sections_org_id;
DROP INDEX IF EXISTS idx_sections_job_id;
DROP TABLE IF EXISTS sections;

ALTER TABLE jobs DROP COLUMN IF EXISTS other_markup;
ALTER TABLE jobs DROP COLUMN IF EXISTS subs_markup;
ALTER TABLE jobs DROP COLUMN IF EXISTS equipment_markup;
ALTER TABLE jobs DROP COLUMN IF EXISTS labor_markup;
ALTER TABLE jobs DROP COLUMN IF EXISTS materials_markup;
