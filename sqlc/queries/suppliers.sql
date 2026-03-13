-- name: CreateSupplier :one
INSERT INTO suppliers (id, name, sort_order)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListSuppliers :many
SELECT * FROM suppliers ORDER BY sort_order ASC, name ASC;

-- name: GetSupplier :one
SELECT * FROM suppliers WHERE id = $1;

-- name: DeleteSupplier :exec
DELETE FROM suppliers WHERE id = $1;

-- name: CountMaterialsBySupplier :one
SELECT COUNT(*) AS total FROM materials WHERE supplier_id = $1;
