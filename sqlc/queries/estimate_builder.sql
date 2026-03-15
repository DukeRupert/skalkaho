-- Sections
-- name: ListSectionsByProject :many
SELECT * FROM sections WHERE project_id = $1 ORDER BY sort_order ASC;

-- name: CreateSection :one
INSERT INTO sections (id, project_id, name, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteSectionsByProject :exec
DELETE FROM sections WHERE project_id = $1;

-- Subcategories (joined through sections)
-- name: ListSubcategoriesBySection :many
SELECT * FROM subcategories WHERE section_id = $1 ORDER BY sort_order ASC;

-- name: CreateSubcategory :one
INSERT INTO subcategories (
    id, section_id, name, sort_order, lump_sum,
    materials_markup, labor_markup, equipment_markup, subs_markup, other_markup,
    materials_markup_enabled, labor_markup_enabled, equipment_markup_enabled,
    subs_markup_enabled, other_markup_enabled
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15
)
RETURNING *;

-- Component Groups
-- name: ListComponentGroupsBySubcategory :many
SELECT * FROM component_groups WHERE subcategory_id = $1 ORDER BY sort_order ASC;

-- name: CreateComponentGroup :one
INSERT INTO component_groups (id, subcategory_id, name, sort_order)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- Line Items
-- name: ListLineItemsBySubcategory :many
SELECT * FROM line_items WHERE subcategory_id = $1 ORDER BY sort_order ASC;

-- name: CreateLineItem :one
INSERT INTO line_items (
    id, subcategory_id, component_group_id, category_type, item_name,
    quantity, unit, unit_price, is_custom, material_id,
    price_override, description, sort_order, subcontractor_id, visual_group
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15
)
RETURNING *;

-- Project markup globals
-- name: GetProjectMarkups :one
SELECT materials_markup, labor_markup, equipment_markup, subs_markup, other_markup
FROM projects WHERE id = $1;

-- name: UpdateProjectMarkups :exec
UPDATE projects
SET materials_markup = $2, labor_markup = $3, equipment_markup = $4,
    subs_markup = $5, other_markup = $6, updated_at = now()
WHERE id = $1;
