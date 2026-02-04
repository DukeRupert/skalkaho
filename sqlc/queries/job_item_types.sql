-- name: CreateJobItemType :one
INSERT INTO job_item_types (id, org_id, job_id, name, slug, color, sort_order, surcharge_percent)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetJobItemType :one
SELECT * FROM job_item_types WHERE id = $1 AND org_id = $2;

-- name: GetJobItemTypeBySlug :one
SELECT * FROM job_item_types
WHERE job_id = $1 AND org_id = $2 AND slug = $3;

-- name: ListJobItemTypes :many
SELECT * FROM job_item_types
WHERE job_id = $1 AND org_id = $2
ORDER BY sort_order ASC, name ASC;

-- name: UpdateJobItemType :one
UPDATE job_item_types SET
    name = $1,
    slug = $2,
    color = $3,
    sort_order = $4,
    surcharge_percent = $5
WHERE id = $6 AND org_id = $7
RETURNING *;

-- name: DeleteJobItemType :exec
DELETE FROM job_item_types WHERE id = $1 AND org_id = $2;

-- name: CountLineItemsByType :one
SELECT COUNT(*) FROM line_items li
JOIN categories c ON li.category_id = c.id
WHERE c.job_id = $1 AND c.org_id = $2 AND li.type = $3 AND li.org_id = $2;
