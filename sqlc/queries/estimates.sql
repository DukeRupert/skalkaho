-- name: CreateEstimate :one
INSERT INTO estimates (id, job_id, version, status, grand_total, notes)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetEstimate :one
SELECT * FROM estimates WHERE id = ?;

-- name: GetEstimateByJobAndVersion :one
SELECT * FROM estimates WHERE job_id = ? AND version = ?;

-- name: ListEstimatesByJob :many
SELECT * FROM estimates
WHERE job_id = ?
ORDER BY version DESC;

-- name: GetLatestEstimateVersion :one
SELECT CAST(COALESCE(MAX(version), 0) AS INTEGER) as max_version
FROM estimates
WHERE job_id = ?;

-- name: UpdateEstimate :one
UPDATE estimates SET
    status = ?,
    notes = ?,
    sent_at = ?,
    responded_at = ?
WHERE id = ?
RETURNING *;

-- name: DeleteEstimate :exec
DELETE FROM estimates WHERE id = ?;

-- name: CreateEstimateCategory :one
INSERT INTO estimate_categories (id, estimate_id, category_id, parent_category_id, tier, name, description, total, sort_order)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetEstimateCategory :one
SELECT * FROM estimate_categories WHERE id = ?;

-- name: ListEstimateCategoriesByEstimate :many
SELECT * FROM estimate_categories
WHERE estimate_id = ?
ORDER BY tier ASC, sort_order ASC;

-- name: ListEstimateCategoriesTier1 :many
SELECT * FROM estimate_categories
WHERE estimate_id = ? AND tier = 1
ORDER BY sort_order ASC;

-- name: ListEstimateCategoriesByParent :many
SELECT * FROM estimate_categories
WHERE estimate_id = ? AND parent_category_id = ?
ORDER BY sort_order ASC;

-- name: UpdateEstimateCategoryDescription :one
UPDATE estimate_categories SET description = ?
WHERE id = ?
RETURNING *;

-- name: DeleteEstimateCategoriesByEstimate :exec
DELETE FROM estimate_categories WHERE estimate_id = ?;
