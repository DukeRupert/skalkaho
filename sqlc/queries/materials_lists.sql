-- Materials Lists

-- name: CreateMaterialsList :one
INSERT INTO materials_lists (id, quote_id, section_id, label, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMaterialsList :one
SELECT
    ml.*,
    s.name AS section_name
FROM materials_lists ml
LEFT JOIN sections s ON s.id = ml.section_id
WHERE ml.id = $1;

-- name: ListMaterialsListsByQuote :many
SELECT
    ml.*,
    s.name AS section_name
FROM materials_lists ml
LEFT JOIN sections s ON s.id = ml.section_id
WHERE ml.quote_id = $1
ORDER BY ml.section_id NULLS FIRST, ml.generated_at DESC;

-- name: UpdateMaterialsListNotes :one
UPDATE materials_lists
SET notes = $2
WHERE id = $1
RETURNING *;

-- name: DeleteMaterialsList :exec
DELETE FROM materials_lists WHERE id = $1;

-- Live item queries (called at render time, not stored)

-- name: GetMaterialsListItems :many
-- Full quote list — all sections, subcategories, line items
-- Organized by section → subcategory for the foreman view
SELECT
    s.name          AS section_name,
    s.sort_order    AS section_order,
    sc.name         AS subcategory_name,
    sc.sort_order   AS subcategory_order,
    li.item_name    AS item_name,
    li.quantity,
    li.unit,
    li.sort_order   AS item_order
FROM line_items li
JOIN subcategories sc ON sc.id = li.subcategory_id
JOIN sections s       ON s.id  = sc.section_id
WHERE s.project_id = $1
  AND li.quantity > 0
ORDER BY
    s.sort_order  ASC,
    sc.sort_order ASC,
    li.sort_order ASC;

-- name: GetMaterialsListItemsBySection :many
-- Single section list — filtered to one section
SELECT
    sc.name         AS subcategory_name,
    sc.sort_order   AS subcategory_order,
    li.item_name    AS item_name,
    li.quantity,
    li.unit,
    li.sort_order   AS item_order
FROM line_items li
JOIN subcategories sc ON sc.id = li.subcategory_id
WHERE sc.section_id = $1
  AND li.quantity > 0
ORDER BY
    sc.sort_order ASC,
    li.sort_order ASC;

-- Get a single section by ID
-- name: GetSection :one
SELECT * FROM sections WHERE id = $1;

-- Get sections for a project (for the dropdown)
-- name: ListSectionsForProject :many
SELECT id, name, sort_order
FROM sections
WHERE project_id = $1
ORDER BY sort_order ASC;
