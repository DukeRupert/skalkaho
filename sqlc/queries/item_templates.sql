-- name: SearchItemTemplates :many
SELECT * FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL) AND (name ILIKE '%' || $2 || '%' OR category ILIKE '%' || $2 || '%')
ORDER BY name
LIMIT 20;

-- name: SearchItemTemplatesByType :many
SELECT * FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL) AND type = $2 AND (name ILIKE '%' || $3 || '%' OR category ILIKE '%' || $3 || '%')
ORDER BY name
LIMIT 20;

-- name: ListItemTemplates :many
SELECT * FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL)
ORDER BY category, name;

-- name: ListItemTemplatesByCategory :many
SELECT * FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL) AND category = $2
ORDER BY name;

-- name: GetItemTemplate :one
SELECT * FROM item_templates
WHERE id = $1 AND org_id = $2;

-- name: CreateItemTemplate :one
INSERT INTO item_templates (org_id, type, category, name, default_unit, default_price, subcategory)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: DeleteItemTemplate :exec
DELETE FROM item_templates
WHERE id = $1 AND org_id = $2;

-- name: UpdateItemTemplate :one
UPDATE item_templates
SET type = $1, category = $2, name = $3, default_unit = $4, default_price = $5, subcategory = $6
WHERE id = $7 AND org_id = $8
RETURNING *;

-- name: UpdateItemTemplatePrice :exec
UPDATE item_templates SET default_price = $1 WHERE id = $2 AND org_id = $3;

-- name: UpdateItemTemplatePriceAndName :exec
UPDATE item_templates SET default_price = $1, name = $2 WHERE id = $3 AND org_id = $4;

-- name: ListItemTemplateCategories :many
SELECT DISTINCT category FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL)
ORDER BY category;

-- name: ListItemTemplateSubcategories :many
SELECT subcategory, COUNT(*)::bigint as item_count FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL) AND category = $2
GROUP BY subcategory
ORDER BY subcategory;

-- name: BulkUpdateItemTemplateSubcategory :exec
UPDATE item_templates SET subcategory = $1
WHERE id = ANY(@ids::bigint[]) AND org_id = $2;

-- name: ListItemTemplatesBySubcategory :many
SELECT * FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL) AND category = $2 AND subcategory = $3
ORDER BY name;
