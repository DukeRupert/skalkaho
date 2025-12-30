-- name: CreatePriceImport :one
INSERT INTO price_imports (id, filename, status, total_rows)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetPriceImport :one
SELECT * FROM price_imports WHERE id = ?;

-- name: ListPriceImports :many
SELECT * FROM price_imports
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: UpdatePriceImportStatus :one
UPDATE price_imports
SET status = ?, matched_rows = ?, error_message = ?
WHERE id = ?
RETURNING *;

-- name: MarkPriceImportApplied :one
UPDATE price_imports
SET status = 'applied', applied_at = datetime('now')
WHERE id = ?
RETURNING *;

-- name: CreatePriceImportMatch :one
INSERT INTO price_import_matches (
    import_id, row_number, source_name, source_unit, source_price,
    matched_template_id, confidence, match_reason, status
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: ListMatchesByImport :many
SELECT
    m.*,
    t.name as template_name,
    t.default_unit as template_unit,
    t.default_price as template_price
FROM price_import_matches m
LEFT JOIN item_templates t ON m.matched_template_id = t.id
WHERE m.import_id = ?
ORDER BY m.confidence DESC, m.row_number;

-- name: UpdateMatchStatus :one
UPDATE price_import_matches SET status = ? WHERE id = ? RETURNING *;

-- name: BulkAutoApproveMatches :exec
UPDATE price_import_matches
SET status = 'auto_approved'
WHERE import_id = ? AND confidence >= ? AND status = 'pending';

-- name: ListApprovedMatches :many
SELECT
    m.*,
    t.name as template_name
FROM price_import_matches m
JOIN item_templates t ON m.matched_template_id = t.id
WHERE m.import_id = ? AND m.status IN ('approved', 'auto_approved');

-- name: CountMatchesByStatus :many
SELECT status, COUNT(*) as count
FROM price_import_matches
WHERE import_id = ?
GROUP BY status;
