-- name: SearchItemTemplates :many
SELECT * FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL) AND name LIKE '%' || $2 || '%'
ORDER BY name
LIMIT 10;

-- name: SearchItemTemplatesByType :many
SELECT * FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL) AND type = $2 AND name LIKE '%' || $3 || '%'
ORDER BY name
LIMIT 10;

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
INSERT INTO item_templates (org_id, type, category, name, default_unit, default_price)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteItemTemplate :exec
DELETE FROM item_templates
WHERE id = $1 AND org_id = $2;

-- name: UpdateItemTemplate :one
UPDATE item_templates
SET type = $1, category = $2, name = $3, default_unit = $4, default_price = $5
WHERE id = $6 AND org_id = $7
RETURNING *;

-- name: UpdateItemTemplatePrice :exec
UPDATE item_templates SET default_price = $1 WHERE id = $2 AND org_id = $3;

-- name: UpdateItemTemplatePriceAndName :exec
UPDATE item_templates SET default_price = $1, name = $2 WHERE id = $3 AND org_id = $4;

-- name: ListItemTemplateCategories :many
SELECT DISTINCT category FROM item_templates
WHERE (org_id = $1 OR org_id IS NULL)
ORDER BY category;
