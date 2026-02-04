-- name: CreateEstimate :one
INSERT INTO estimates (id, job_id, version, status, grand_total, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetEstimate :one
SELECT * FROM estimates WHERE id = $1;

-- name: GetEstimateByJobAndVersion :one
SELECT * FROM estimates WHERE job_id = $1 AND version = $2;

-- name: ListEstimatesByJob :many
SELECT * FROM estimates
WHERE job_id = $1
ORDER BY version DESC;

-- name: GetLatestEstimateVersion :one
SELECT CAST(COALESCE(MAX(version), 0) AS INTEGER) as max_version
FROM estimates
WHERE job_id = $1;

-- name: UpdateEstimate :one
UPDATE estimates SET
    status = $1,
    notes = $2,
    sent_at = $3,
    responded_at = $4
WHERE id = $5
RETURNING *;

-- name: UpdateEstimateStatus :one
UPDATE estimates SET status = $1
WHERE id = $2
RETURNING *;

-- name: MarkEstimateSent :one
UPDATE estimates SET status = 'sent', sent_at = datetime('now')
WHERE id = $1
RETURNING *;

-- name: MarkEstimateAccepted :one
UPDATE estimates SET status = 'accepted', responded_at = datetime('now')
WHERE id = $1
RETURNING *;

-- name: DeleteEstimate :exec
DELETE FROM estimates WHERE id = $1;

-- name: CreateEstimateCategory :one
INSERT INTO estimate_categories (id, estimate_id, category_id, parent_category_id, tier, name, description, total, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetEstimateCategory :one
SELECT * FROM estimate_categories WHERE id = $1;

-- name: ListEstimateCategoriesByEstimate :many
SELECT * FROM estimate_categories
WHERE estimate_id = $1
ORDER BY tier ASC, sort_order ASC;

-- name: ListEstimateCategoriesTier1 :many
SELECT * FROM estimate_categories
WHERE estimate_id = $1 AND tier = 1
ORDER BY sort_order ASC;

-- name: ListEstimateCategoriesByParent :many
SELECT * FROM estimate_categories
WHERE estimate_id = $1 AND parent_category_id = $2
ORDER BY sort_order ASC;

-- name: UpdateEstimateCategoryDescription :one
UPDATE estimate_categories SET description = $1
WHERE id = $2
RETURNING *;

-- name: DeleteEstimateCategoriesByEstimate :exec
DELETE FROM estimate_categories WHERE estimate_id = $1;
