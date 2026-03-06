-- name: CreateLineItem :one
INSERT INTO line_items (id, org_id, category_id, type, name, description, quantity, unit, unit_price, surcharge_percent, sort_order, tag)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetLineItem :one
SELECT * FROM line_items
WHERE id = $1 AND org_id = $2;

-- name: ListLineItemsByCategory :many
SELECT * FROM line_items
WHERE category_id = $1 AND org_id = $2
ORDER BY sort_order ASC;

-- name: ListLineItemsByJob :many
SELECT li.* FROM line_items li
JOIN categories c ON li.category_id = c.id
WHERE c.job_id = $1 AND li.org_id = $2 AND c.org_id = $2
ORDER BY li.sort_order ASC;

-- name: UpdateLineItem :one
UPDATE line_items SET
    type = $1,
    name = $2,
    description = $3,
    quantity = $4,
    unit = $5,
    unit_price = $6,
    surcharge_percent = $7,
    sort_order = $8,
    tag = $9
WHERE id = $10 AND org_id = $11
RETURNING *;

-- name: UpdateLineItemSortOrder :exec
UPDATE line_items SET sort_order = $1
WHERE id = $2 AND org_id = $3;

-- name: DeleteLineItem :exec
DELETE FROM line_items
WHERE id = $1 AND org_id = $2;
