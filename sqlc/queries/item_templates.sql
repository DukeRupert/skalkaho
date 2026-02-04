-- name: SearchItemTemplates :many
SELECT * FROM item_templates
WHERE name LIKE '%' || $1 || '%'
ORDER BY name
LIMIT 10;

-- name: SearchItemTemplatesByType :many
SELECT * FROM item_templates
WHERE type = $1 AND name LIKE '%' || $2 || '%'
ORDER BY name
LIMIT 10;

-- name: ListItemTemplates :many
SELECT * FROM item_templates
ORDER BY category, name;

-- name: ListItemTemplatesByCategory :many
SELECT * FROM item_templates
WHERE category = $1
ORDER BY name;

-- name: GetItemTemplate :one
SELECT * FROM item_templates
WHERE id = $1;

-- name: CreateItemTemplate :one
INSERT INTO item_templates (type, category, name, default_unit, default_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteItemTemplate :exec
DELETE FROM item_templates
WHERE id = $1;

-- name: UpdateItemTemplate :one
UPDATE item_templates
SET type = $1, category = $2, name = $3, default_unit = $4, default_price = $5
WHERE id = $6
RETURNING *;

-- name: UpdateItemTemplatePrice :exec
UPDATE item_templates SET default_price = $1 WHERE id = $2;

-- name: UpdateItemTemplatePriceAndName :exec
UPDATE item_templates SET default_price = $1, name = $2 WHERE id = $3;
