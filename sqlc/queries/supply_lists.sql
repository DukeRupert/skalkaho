-- Supply Lists

-- name: CreateSupplyList :one
INSERT INTO supply_lists (id, quote_id, supplier_id, label, estimated_delivery, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetSupplyList :one
SELECT * FROM supply_lists WHERE id = $1;

-- name: ListSupplyListsByQuote :many
SELECT
    sl.*,
    COUNT(sli.id) AS item_count,
    COALESCE(SUM(sli.extended_price), 0)::REAL AS total_value
FROM supply_lists sl
LEFT JOIN supply_list_items sli ON sli.supply_list_id = sl.id
WHERE sl.quote_id = $1
GROUP BY sl.id
ORDER BY sl.supplier_id NULLS FIRST, sl.generated_at DESC;

-- name: UpdateSupplyListDates :one
UPDATE supply_lists
SET estimated_delivery = $2,
    actual_delivery    = $3
WHERE id = $1
RETURNING *;

-- name: UpdateSupplyListNotes :one
UPDATE supply_lists
SET notes = $2
WHERE id = $1
RETURNING *;

-- name: DeleteSupplyList :exec
DELETE FROM supply_lists WHERE id = $1;

-- Supply List Items

-- name: CreateSupplyListItem :one
INSERT INTO supply_list_items (
    id, supply_list_id, material_id, supplier_id, supplier_name,
    material_name, supplier_code, quantity, unit, unit_price, extended_price, sort_order
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: ListSupplyListItems :many
SELECT * FROM supply_list_items
WHERE supply_list_id = $1
ORDER BY supplier_name NULLS LAST, sort_order ASC, material_name ASC;

-- name: ListSupplyListItemsGrouped :many
SELECT * FROM supply_list_items
WHERE supply_list_id = $1
ORDER BY
    CASE WHEN supplier_name IS NULL THEN 1 ELSE 0 END,
    supplier_name ASC,
    material_name ASC;

-- Snapshot query: get material line items with supplier info for a project
-- name: ListMaterialLineItemsForProject :many
SELECT
    li.id AS line_item_id,
    li.item_name,
    li.quantity,
    li.unit,
    li.unit_price,
    li.material_id,
    li.sort_order,
    m.name AS material_name,
    m.supplier_code,
    m.supplier_id,
    s.name AS supplier_name
FROM line_items li
JOIN subcategories sc ON sc.id = li.subcategory_id
JOIN sections sec ON sec.id = sc.section_id
JOIN materials m ON m.id = li.material_id
JOIN suppliers s ON s.id = m.supplier_id
WHERE sec.project_id = $1
  AND li.category_type = 'materials'
  AND li.material_id IS NOT NULL
ORDER BY s.name ASC, li.sort_order ASC, m.name ASC;

-- Snapshot query: get ad hoc line items (no material reference) for a project
-- name: ListAdHocLineItemsForProject :many
SELECT
    li.id AS line_item_id,
    li.item_name,
    li.quantity,
    li.unit,
    li.unit_price,
    li.sort_order
FROM line_items li
JOIN subcategories sc ON sc.id = li.subcategory_id
JOIN sections sec ON sec.id = sc.section_id
WHERE sec.project_id = $1
  AND li.category_type = 'materials'
  AND li.material_id IS NULL
ORDER BY li.sort_order ASC, li.item_name ASC;

-- Get distinct suppliers on a project's material line items
-- name: ListSuppliersForProject :many
SELECT DISTINCT s.id, s.name
FROM line_items li
JOIN subcategories sc ON sc.id = li.subcategory_id
JOIN sections sec ON sec.id = sc.section_id
JOIN materials m ON m.id = li.material_id
JOIN suppliers s ON s.id = m.supplier_id
WHERE sec.project_id = $1
  AND li.category_type = 'materials'
  AND li.material_id IS NOT NULL
ORDER BY s.name ASC;
