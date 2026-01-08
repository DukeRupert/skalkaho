-- name: CreateJob :one
INSERT INTO jobs (id, name, customer_name, surcharge_percent, material_surcharge_percent, labor_surcharge_percent, equipment_surcharge_percent, surcharge_mode, status, expires_at, client_id)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
    material_surcharge_percent = ?,
    labor_surcharge_percent = ?,
    equipment_surcharge_percent = ?,
    surcharge_mode = ?,
    status = ?,
    expires_at = ?,
    client_id = ?
WHERE id = ?
RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = ?;

-- name: ListJobsWithEstimateStatus :many
SELECT
    j.*,
    (SELECT COUNT(*) FROM estimates WHERE job_id = j.id) as estimate_count,
    (SELECT status FROM estimates WHERE job_id = j.id ORDER BY version DESC LIMIT 1) as latest_estimate_status,
    (SELECT sr.status FROM signature_requests sr
     INNER JOIN estimates e ON sr.estimate_id = e.id
     WHERE e.job_id = j.id
     ORDER BY sr.created_at DESC LIMIT 1) as latest_signature_status
FROM jobs j
WHERE (@status = '' OR j.status = @status)
ORDER BY j.created_at DESC
LIMIT @limit OFFSET @offset;
