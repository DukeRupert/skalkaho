-- name: SearchItemTemplates :many
SELECT * FROM item_templates
WHERE name LIKE '%' || ? || '%'
ORDER BY name
LIMIT 10;

-- name: SearchItemTemplatesByType :many
SELECT * FROM item_templates
WHERE type = ? AND name LIKE '%' || ? || '%'
ORDER BY name
LIMIT 10;

-- name: ListItemTemplates :many
SELECT * FROM item_templates
ORDER BY category, name;

-- name: ListItemTemplatesByCategory :many
SELECT * FROM item_templates
WHERE category = ?
ORDER BY name;

-- name: GetItemTemplate :one
SELECT * FROM item_templates
WHERE id = ?;

-- name: CreateItemTemplate :one
INSERT INTO item_templates (type, category, name, default_unit, default_price)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteItemTemplate :exec
DELETE FROM item_templates
WHERE id = ?;
