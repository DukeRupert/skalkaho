-- name: CreateJobItemType :one
INSERT INTO job_item_types (id, job_id, name, slug, color, sort_order, surcharge_percent)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetJobItemType :one
SELECT * FROM job_item_types WHERE id = ?;

-- name: GetJobItemTypeBySlug :one
SELECT * FROM job_item_types
WHERE job_id = ? AND slug = ?;

-- name: ListJobItemTypes :many
SELECT * FROM job_item_types
WHERE job_id = ?
ORDER BY sort_order ASC, name ASC;

-- name: UpdateJobItemType :one
UPDATE job_item_types SET
    name = ?,
    slug = ?,
    color = ?,
    sort_order = ?,
    surcharge_percent = ?
WHERE id = ?
RETURNING *;

-- name: DeleteJobItemType :exec
DELETE FROM job_item_types WHERE id = ?;

-- name: CountLineItemsByType :one
SELECT COUNT(*) FROM line_items li
JOIN categories c ON li.category_id = c.id
WHERE c.job_id = ? AND li.type = ?;
