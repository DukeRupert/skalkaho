-- +goose Up
CREATE TABLE sections (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_sections_project_id ON sections(project_id);

CREATE TABLE subcategories (
    id TEXT PRIMARY KEY,
    section_id TEXT NOT NULL REFERENCES sections(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    lump_sum REAL NOT NULL DEFAULT 0,
    materials_markup REAL,
    labor_markup REAL,
    equipment_markup REAL,
    subs_markup REAL,
    other_markup REAL,
    materials_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    labor_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    equipment_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    subs_markup_enabled BOOLEAN NOT NULL DEFAULT true,
    other_markup_enabled BOOLEAN NOT NULL DEFAULT true
);

CREATE INDEX idx_subcategories_section_id ON subcategories(section_id);

CREATE TABLE component_groups (
    id TEXT PRIMARY KEY,
    subcategory_id TEXT NOT NULL REFERENCES subcategories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_component_groups_subcategory_id ON component_groups(subcategory_id);

CREATE TABLE line_items (
    id TEXT PRIMARY KEY,
    subcategory_id TEXT NOT NULL REFERENCES subcategories(id) ON DELETE CASCADE,
    component_group_id TEXT REFERENCES component_groups(id) ON DELETE SET NULL,
    category_type TEXT NOT NULL CHECK (category_type IN ('materials', 'labor', 'equipment', 'subs', 'other')),
    item_name TEXT NOT NULL,
    quantity REAL NOT NULL DEFAULT 1,
    unit TEXT NOT NULL DEFAULT 'ea',
    unit_price REAL NOT NULL DEFAULT 0,
    is_custom BOOLEAN NOT NULL DEFAULT true,
    material_id TEXT REFERENCES materials(id) ON DELETE SET NULL,
    price_override BOOLEAN NOT NULL DEFAULT false,
    description TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_line_items_subcategory_id ON line_items(subcategory_id);
CREATE INDEX idx_line_items_component_group_id ON line_items(component_group_id);

-- +goose Down
DROP TABLE IF EXISTS line_items;
DROP TABLE IF EXISTS component_groups;
DROP TABLE IF EXISTS subcategories;
DROP TABLE IF EXISTS sections;
