-- name: CreateImportJob :one
INSERT INTO material_import_jobs (id, supplier_id, filename, temp_file_path)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ClaimNextImportJob :one
UPDATE material_import_jobs
SET status = 'processing', started_at = NOW()
WHERE id = (
    SELECT id FROM material_import_jobs
    WHERE status = 'pending'
    ORDER BY created_at ASC
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: GetImportJob :one
SELECT * FROM material_import_jobs WHERE id = $1;

-- name: UpdateImportJobProgress :one
UPDATE material_import_jobs
SET rows_total   = $2,
    rows_created = $3,
    rows_updated = $4,
    rows_skipped = $5,
    rows_flagged = $6,
    rows_failed  = $7
WHERE id = $1
RETURNING *;

-- name: CompleteImportJob :one
UPDATE material_import_jobs
SET status = 'complete', completed_at = NOW()
WHERE id = $1
RETURNING *;

-- name: FailImportJob :one
UPDATE material_import_jobs
SET status = 'failed', completed_at = NOW(), error_message = $2
WHERE id = $1
RETURNING *;

-- name: ListImportJobs :many
SELECT j.*, s.name AS supplier_name
FROM material_import_jobs j
JOIN suppliers s ON s.id = j.supplier_id
ORDER BY j.created_at DESC;

-- name: CreateImportFlag :one
INSERT INTO material_import_flags (
    id, import_job_id, supplier_id,
    raw_name, raw_unit, raw_price, raw_supplier_code,
    candidate_material_id, confidence, claude_reasoning
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: ListImportFlags :many
SELECT f.*, m.name AS candidate_name, m.unit_price AS candidate_price
FROM material_import_flags f
LEFT JOIN materials m ON m.id = f.candidate_material_id
WHERE f.import_job_id = $1
ORDER BY f.confidence DESC;

-- name: ListPendingFlags :many
SELECT f.*, m.name AS candidate_name, m.unit_price AS candidate_price
FROM material_import_flags f
LEFT JOIN materials m ON m.id = f.candidate_material_id
WHERE f.import_job_id = $1 AND f.resolution IS NULL
ORDER BY f.confidence DESC;

-- name: GetImportFlag :one
SELECT * FROM material_import_flags WHERE id = $1;

-- name: ResolveFlag :one
UPDATE material_import_flags
SET resolution = $2, resolved_at = NOW()
WHERE id = $1
RETURNING *;
