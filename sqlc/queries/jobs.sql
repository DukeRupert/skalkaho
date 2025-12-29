-- name: CreateJob :one
INSERT INTO jobs (id, name, customer_name, surcharge_percent, surcharge_mode, status, expires_at)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetJob :one
SELECT * FROM jobs
WHERE id = ?;

-- name: ListJobs :many
SELECT * FROM jobs
ORDER BY created_at DESC;

-- name: ListJobsPaginated :many
SELECT * FROM jobs
WHERE (@status = '' OR status = @status)
ORDER BY created_at DESC
LIMIT @limit OFFSET @offset;

-- name: ListJobsPaginatedByName :many
SELECT * FROM jobs
WHERE (@status = '' OR status = @status)
ORDER BY name ASC
LIMIT @limit OFFSET @offset;

-- name: ListJobsPaginatedByNameDesc :many
SELECT * FROM jobs
WHERE (@status = '' OR status = @status)
ORDER BY name DESC
LIMIT @limit OFFSET @offset;

-- name: ListJobsPaginatedOldest :many
SELECT * FROM jobs
WHERE (@status = '' OR status = @status)
ORDER BY created_at ASC
LIMIT @limit OFFSET @offset;

-- name: CountJobs :one
SELECT COUNT(*) FROM jobs
WHERE (@status = '' OR status = @status);

-- name: UpdateJobStatus :one
UPDATE jobs SET status = ? WHERE id = ? RETURNING *;

-- name: UpdateJob :one
UPDATE jobs SET
    name = ?,
    customer_name = ?,
    surcharge_percent = ?,
    surcharge_mode = ?,
    status = ?,
    expires_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = ?;
