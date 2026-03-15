-- +goose Up
CREATE TABLE templates (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE template_sections (
    id          TEXT PRIMARY KEY,
    template_id TEXT NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    sort_order  INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_template_sections_template_id ON template_sections(template_id);

CREATE TABLE template_subcategories (
    id                  TEXT PRIMARY KEY,
    template_section_id TEXT NOT NULL REFERENCES template_sections(id) ON DELETE CASCADE,
    name                TEXT NOT NULL,
    sort_order          INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_template_subcategories_section_id ON template_subcategories(template_section_id);

CREATE TABLE template_component_groups (
    id                       TEXT PRIMARY KEY,
    template_subcategory_id  TEXT NOT NULL REFERENCES template_subcategories(id) ON DELETE CASCADE,
    name                     TEXT NOT NULL,
    sort_order               INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_template_component_groups_subcategory_id ON template_component_groups(template_subcategory_id);

-- +goose Down
DROP TABLE IF EXISTS template_component_groups;
DROP TABLE IF EXISTS template_subcategories;
DROP TABLE IF EXISTS template_sections;
DROP TABLE IF EXISTS templates;
