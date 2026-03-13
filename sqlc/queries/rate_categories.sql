-- name: CreateRateCategory :one
INSERT INTO rate_categories (id, name, sort_order)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListRateCategories :many
SELECT * FROM rate_categories ORDER BY sort_order ASC, name ASC;

-- name: GetRateCategory :one
SELECT * FROM rate_categories WHERE id = $1;

-- name: DeleteRateCategory :exec
DELETE FROM rate_categories WHERE id = $1;

-- name: CountRatesByCategory :one
SELECT COUNT(*) AS total FROM rates WHERE category_id = $1;
