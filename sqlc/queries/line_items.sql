-- name: CreateLineItem :one
INSERT INTO line_items (id, category_id, type, name, description, quantity, unit, unit_price, surcharge_percent, sort_order)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetLineItem :one
SELECT * FROM line_items
WHERE id = ?;

-- name: ListLineItemsByCategory :many
SELECT * FROM line_items
WHERE category_id = ?
ORDER BY sort_order ASC;

-- name: ListLineItemsByJob :many
SELECT li.* FROM line_items li
JOIN categories c ON li.category_id = c.id
WHERE c.job_id = ?
ORDER BY li.sort_order ASC;

-- name: UpdateLineItem :one
UPDATE line_items SET
    type = ?,
    name = ?,
    description = ?,
    quantity = ?,
    unit = ?,
    unit_price = ?,
    surcharge_percent = ?,
    sort_order = ?
WHERE id = ?
RETURNING *;

-- name: DeleteLineItem :exec
DELETE FROM line_items
WHERE id = ?;
