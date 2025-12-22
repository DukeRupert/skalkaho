-- name: CreateCategory :one
INSERT INTO categories (id, job_id, parent_id, name, surcharge_percent, sort_order)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = ?;

-- name: ListCategoriesByJob :many
SELECT * FROM categories
WHERE job_id = ?
ORDER BY sort_order ASC;

-- name: ListTopLevelCategories :many
SELECT * FROM categories
WHERE job_id = ? AND parent_id IS NULL
ORDER BY sort_order ASC;

-- name: ListChildCategories :many
SELECT * FROM categories
WHERE parent_id = ?
ORDER BY sort_order ASC;

-- name: UpdateCategory :one
UPDATE categories SET
    name = ?,
    surcharge_percent = ?,
    sort_order = ?
WHERE id = ?
RETURNING *;

-- name: UpdateCategoryParent :one
UPDATE categories SET
    parent_id = ?
WHERE id = ?
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = ?;

-- name: CountCategoryAncestors :one
WITH RECURSIVE ancestors AS (
    SELECT categories.id, categories.parent_id, 0 as depth
    FROM categories
    WHERE categories.id = ?
    UNION ALL
    SELECT c.id, c.parent_id, a.depth + 1
    FROM categories c
    JOIN ancestors a ON c.id = a.parent_id
)
SELECT MAX(depth) as max_depth FROM ancestors;
