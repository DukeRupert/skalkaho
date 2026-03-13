-- Estimate Builder queries for the 4-level hierarchy

-- =====================
-- Sections
-- =====================

-- name: ListSectionsByJob :many
SELECT * FROM sections
WHERE job_id = $1 AND org_id = $2
ORDER BY sort_order ASC;

-- name: CreateSection :one
INSERT INTO sections (id, org_id, job_id, name, sort_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteSectionsByJob :exec
DELETE FROM sections
WHERE job_id = $1 AND org_id = $2;

-- =====================
-- Subcategories
-- =====================

-- name: ListSubcategoriesByJob :many
SELECT sc.* FROM subcategories sc
JOIN sections s ON sc.section_id = s.id
WHERE s.job_id = $1 AND sc.org_id = $2
ORDER BY sc.sort_order ASC;

-- name: CreateSubcategory :one
INSERT INTO subcategories (
    id, org_id, section_id, name, sort_order, lump_sum,
    materials_markup, labor_markup, equipment_markup, subs_markup, other_markup,
    materials_markup_enabled, labor_markup_enabled, equipment_markup_enabled,
    subs_markup_enabled, other_markup_enabled
) VALUES (
    $1, $2, $3, $4, $5, $6,
    $7, $8, $9, $10, $11,
    $12, $13, $14, $15, $16
)
RETURNING *;

-- =====================
-- Component Groups
-- =====================

-- name: ListComponentGroupsByJob :many
SELECT cg.* FROM component_groups cg
JOIN subcategories sc ON cg.subcategory_id = sc.id
JOIN sections s ON sc.section_id = s.id
WHERE s.job_id = $1 AND cg.org_id = $2
ORDER BY cg.sort_order ASC;

-- name: CreateComponentGroup :one
INSERT INTO component_groups (id, org_id, subcategory_id, name, sort_order)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- =====================
-- Estimate Line Items
-- =====================

-- name: ListEstimateLineItemsByJob :many
SELECT eli.* FROM estimate_line_items eli
JOIN subcategories sc ON eli.subcategory_id = sc.id
JOIN sections s ON sc.section_id = s.id
WHERE s.job_id = $1 AND eli.org_id = $2
ORDER BY eli.sort_order ASC;

-- name: CreateEstimateLineItem :one
INSERT INTO estimate_line_items (
    id, org_id, subcategory_id, component_group_id, category_type,
    item_name, quantity, unit, unit_price,
    is_custom, material_id, price_override, description, sort_order
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    $10, $11, $12, $13, $14
)
RETURNING *;

-- =====================
-- Job markup globals
-- =====================

-- name: GetJobMarkupGlobals :one
SELECT id, materials_markup, labor_markup, equipment_markup, subs_markup, other_markup
FROM jobs
WHERE id = $1 AND org_id = $2;

-- name: UpdateJobMarkupGlobals :exec
UPDATE jobs SET
    materials_markup = $1,
    labor_markup = $2,
    equipment_markup = $3,
    subs_markup = $4,
    other_markup = $5
WHERE id = $6 AND org_id = $7;
