-- name: CreateJob :one
INSERT INTO jobs (id, name, customer_name, surcharge_percent, surcharge_mode)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: GetJob :one
SELECT * FROM jobs
WHERE id = ?;

-- name: ListJobs :many
SELECT * FROM jobs
ORDER BY created_at DESC;

-- name: UpdateJob :one
UPDATE jobs SET
    name = ?,
    customer_name = ?,
    surcharge_percent = ?,
    surcharge_mode = ?
WHERE id = ?
RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = ?;
