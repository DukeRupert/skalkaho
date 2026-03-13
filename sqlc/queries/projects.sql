-- name: CreateProject :one
INSERT INTO projects (id, name, client_id, client_name, description, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProject :one
SELECT * FROM projects WHERE id = $1;

-- name: ListProjects :many
SELECT * FROM projects ORDER BY updated_at DESC;

-- name: ListProjectsByStatus :many
SELECT * FROM projects WHERE status = $1 ORDER BY updated_at DESC;

-- name: SearchProjects :many
SELECT * FROM projects
WHERE name ILIKE '%' || @search_term::text || '%' OR client_name ILIKE '%' || @search_term::text || '%'
ORDER BY updated_at DESC;

-- name: SearchProjectsByStatus :many
SELECT * FROM projects
WHERE (name ILIKE '%' || @search_term::text || '%' OR client_name ILIKE '%' || @search_term::text || '%')
  AND status = @status
ORDER BY updated_at DESC;

-- name: UpdateProject :one
UPDATE projects
SET name = $2, client_id = $3, client_name = $4, description = $5, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateProjectStatus :exec
UPDATE projects SET status = $2, updated_at = now() WHERE id = $1;

-- name: UpdateProjectTotal :exec
UPDATE projects SET total = $2, updated_at = now() WHERE id = $1;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;

-- name: CountProjectsByStatus :one
SELECT
    COUNT(*) FILTER (WHERE status = 'Draft') AS draft_count,
    COUNT(*) FILTER (WHERE status = 'In Review') AS in_review_count,
    COUNT(*) FILTER (WHERE status = 'Active') AS active_count,
    COUNT(*) FILTER (WHERE status = 'Completed') AS completed_count,
    COALESCE(SUM(total) FILTER (WHERE status = 'Draft'), 0)::real AS draft_total,
    COALESCE(SUM(total) FILTER (WHERE status = 'In Review'), 0)::real AS in_review_total,
    COALESCE(SUM(total) FILTER (WHERE status = 'Active'), 0)::real AS active_total,
    COALESCE(SUM(total) FILTER (WHERE status = 'Completed'), 0)::real AS completed_total
FROM projects;
