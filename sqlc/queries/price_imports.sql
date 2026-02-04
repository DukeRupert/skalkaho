-- name: CreatePriceImport :one
INSERT INTO price_imports (id, filename, status, total_rows)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetPriceImport :one
SELECT * FROM price_imports WHERE id = $1;

-- name: ListPriceImports :many
SELECT * FROM price_imports
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdatePriceImportStatus :one
UPDATE price_imports
SET status = $1, matched_rows = $2, error_message = $3, total_rows = $4
WHERE id = $5
RETURNING *;

-- name: MarkPriceImportApplied :one
UPDATE price_imports
SET status = 'applied', applied_at = datetime('now')
WHERE id = $1
RETURNING *;

-- name: CreatePriceImportMatch :one
INSERT INTO price_import_matches (
    import_id, row_number, source_name, source_unit, source_price,
    matched_template_id, confidence, match_reason, status
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: ListMatchesByImport :many
SELECT
    m.*,
    t.name as template_name,
    t.default_unit as template_unit,
    t.default_price as template_price
FROM price_import_matches m
LEFT JOIN item_templates t ON m.matched_template_id = t.id
WHERE m.import_id = $1
ORDER BY m.confidence DESC, m.row_number;

-- name: UpdateMatchStatus :one
UPDATE price_import_matches SET status = $1 WHERE id = $2 RETURNING *;

-- name: BulkAutoApproveMatches :exec
UPDATE price_import_matches
SET status = 'auto_approved'
WHERE import_id = $1 AND confidence >= $2 AND status = 'pending';

-- name: ListApprovedMatches :many
SELECT
    m.*,
    t.name as template_name
FROM price_import_matches m
JOIN item_templates t ON m.matched_template_id = t.id
WHERE m.import_id = $1 AND m.status IN ('approved', 'auto_approved');

-- name: CountMatchesByStatus :many
SELECT status, COUNT(*) as count
FROM price_import_matches
WHERE import_id = $1
GROUP BY status;

-- name: UpdateMatchWithName :one
UPDATE price_import_matches
SET status = $1, new_name = $2
WHERE id = $3
RETURNING *;

-- name: ListUnmatchedItems :many
SELECT * FROM price_import_matches
WHERE import_id = $1 AND matched_template_id IS NULL AND status = 'pending'
ORDER BY row_number;

-- name: MarkMatchAsCreated :one
UPDATE price_import_matches
SET status = 'created', matched_template_id = $1
WHERE id = $2
RETURNING *;
