-- name: CreateCategory :one
INSERT INTO categories (id, org_id, job_id, parent_id, name, surcharge_percent, sort_order)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetCategory :one
SELECT * FROM categories
WHERE id = $1 AND org_id = $2;

-- name: ListCategoriesByJob :many
SELECT * FROM categories
WHERE job_id = $1 AND org_id = $2
ORDER BY sort_order ASC;

-- name: ListTopLevelCategories :many
SELECT * FROM categories
WHERE job_id = $1 AND org_id = $2 AND parent_id IS NULL
ORDER BY sort_order ASC;

-- name: ListChildCategories :many
SELECT * FROM categories
WHERE parent_id = $1 AND org_id = $2
ORDER BY sort_order ASC;

-- name: UpdateCategory :one
UPDATE categories SET
    name = $1,
    surcharge_percent = $2,
    sort_order = $3
WHERE id = $4 AND org_id = $5
RETURNING *;

-- name: UpdateCategoryParent :one
UPDATE categories SET
    parent_id = $1
WHERE id = $2 AND org_id = $3
RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1 AND org_id = $2;

-- name: CountCategoryAncestors :one
WITH RECURSIVE ancestors AS (
    SELECT categories.id, categories.parent_id, 0 as depth
    FROM categories
    WHERE categories.id = $1 AND categories.org_id = $2
    UNION ALL
    SELECT c.id, c.parent_id, a.depth + 1
    FROM categories c
    JOIN ancestors a ON c.id = a.parent_id
    WHERE c.org_id = $2
)
SELECT MAX(depth) as max_depth FROM ancestors;
