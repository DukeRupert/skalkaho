-- name: CreateJob :one
INSERT INTO jobs (id, org_id, name, customer_name, surcharge_percent, material_surcharge_percent, labor_surcharge_percent, equipment_surcharge_percent, surcharge_mode, status, expires_at, client_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1 AND org_id = $2;

-- name: ListJobs :many
SELECT * FROM jobs
WHERE org_id = $1
ORDER BY created_at DESC;

-- name: ListJobsPaginated :many
SELECT * FROM jobs
WHERE org_id = sqlc.arg('org_id')
  AND (sqlc.arg('status') = '' OR status = sqlc.arg('status'))
ORDER BY created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListJobsPaginatedByName :many
SELECT * FROM jobs
WHERE org_id = sqlc.arg('org_id')
  AND (sqlc.arg('status') = '' OR status = sqlc.arg('status'))
ORDER BY name ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListJobsPaginatedByNameDesc :many
SELECT * FROM jobs
WHERE org_id = sqlc.arg('org_id')
  AND (sqlc.arg('status') = '' OR status = sqlc.arg('status'))
ORDER BY name DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListJobsPaginatedOldest :many
SELECT * FROM jobs
WHERE org_id = sqlc.arg('org_id')
  AND (sqlc.arg('status') = '' OR status = sqlc.arg('status'))
ORDER BY created_at ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: CountJobs :one
SELECT COUNT(*) FROM jobs
WHERE org_id = sqlc.arg('org_id')
  AND (sqlc.arg('status') = '' OR status = sqlc.arg('status'));

-- name: UpdateJobStatus :one
UPDATE jobs SET status = $1 WHERE id = $2 AND org_id = $3 RETURNING *;

-- name: UpdateJob :one
UPDATE jobs SET
    name = $1,
    customer_name = $2,
    surcharge_percent = $3,
    material_surcharge_percent = $4,
    labor_surcharge_percent = $5,
    equipment_surcharge_percent = $6,
    surcharge_mode = $7,
    status = $8,
    expires_at = $9,
    client_id = $10
WHERE id = $11 AND org_id = $12
RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = $1 AND org_id = $2;

-- name: ListJobsWithEstimateStatus :many
SELECT
    j.*,
    (SELECT COUNT(*) FROM estimates WHERE job_id = j.id AND org_id = j.org_id) as estimate_count,
    (SELECT status FROM estimates WHERE job_id = j.id AND org_id = j.org_id ORDER BY version DESC LIMIT 1) as latest_estimate_status,
    (SELECT sr.status FROM signature_requests sr
     INNER JOIN estimates e ON sr.estimate_id = e.id
     WHERE e.job_id = j.id AND sr.org_id = j.org_id
     ORDER BY sr.created_at DESC LIMIT 1) as latest_signature_status
FROM jobs j
WHERE j.org_id = sqlc.arg('org_id')
  AND (sqlc.arg('status') = '' OR j.status = sqlc.arg('status'))
ORDER BY j.created_at DESC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ListJobsByClientWithEstimateStatus :many
SELECT
    j.*,
    (SELECT COUNT(*) FROM estimates WHERE job_id = j.id AND org_id = j.org_id) as estimate_count,
    (SELECT status FROM estimates WHERE job_id = j.id AND org_id = j.org_id ORDER BY version DESC LIMIT 1) as latest_estimate_status,
    (SELECT sr.status FROM signature_requests sr
     INNER JOIN estimates e ON sr.estimate_id = e.id
     WHERE e.job_id = j.id AND sr.org_id = j.org_id
     ORDER BY sr.created_at DESC LIMIT 1) as latest_signature_status
FROM jobs j
WHERE j.client_id = $1 AND j.org_id = $2
ORDER BY j.created_at DESC;
