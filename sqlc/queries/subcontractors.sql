-- name: ListTrades :many
SELECT * FROM trades ORDER BY sort_order, name;

-- name: ListSubcontractors :many
SELECT
    s.*,
    t.name AS primary_trade
FROM subcontractors s
LEFT JOIN subcontractor_trades st ON st.subcontractor_id = s.id AND st.position = 0
LEFT JOIN trades t ON t.id = st.trade_id
ORDER BY s.is_favorite DESC, s.name ASC;

-- name: ListSubcontractorsByTrade :many
SELECT
    s.*,
    t.name AS primary_trade
FROM subcontractors s
JOIN subcontractor_trades st ON st.subcontractor_id = s.id AND st.position = 0
JOIN trades t ON t.id = st.trade_id
WHERE t.id = $1
ORDER BY s.is_favorite DESC, s.name ASC;

-- name: SearchSubcontractors :many
SELECT
    s.*,
    t.name AS primary_trade
FROM subcontractors s
LEFT JOIN subcontractor_trades st ON st.subcontractor_id = s.id AND st.position = 0
LEFT JOIN trades t ON t.id = st.trade_id
WHERE s.name ILIKE '%' || @search_term::text || '%'
   OR s.company ILIKE '%' || @search_term::text || '%'
ORDER BY s.is_favorite DESC, s.name ASC;

-- name: GetSubcontractor :one
SELECT * FROM subcontractors WHERE id = $1;

-- name: GetSubcontractorTrades :many
SELECT t.*, st.position
FROM trades t
JOIN subcontractor_trades st ON st.trade_id = t.id
WHERE st.subcontractor_id = $1
ORDER BY st.position ASC;

-- name: CreateSubcontractor :one
INSERT INTO subcontractors (id, name, company, phone, email, address, notes, is_favorite)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateSubcontractor :one
UPDATE subcontractors
SET name = $2, company = $3, phone = $4, email = $5, address = $6, notes = $7,
    is_favorite = $8, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: ToggleFavorite :one
UPDATE subcontractors
SET is_favorite = NOT is_favorite,
    updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteSubcontractor :exec
DELETE FROM subcontractors WHERE id = $1;

-- name: SetSubcontractorTrade :exec
INSERT INTO subcontractor_trades (subcontractor_id, trade_id, position)
VALUES ($1, $2, $3);

-- name: ClearSubcontractorTrades :exec
DELETE FROM subcontractor_trades WHERE subcontractor_id = $1;

-- name: CountSubcontractors :one
SELECT count(*) FROM subcontractors;

-- name: CountFavoriteSubcontractors :one
SELECT count(*) FROM subcontractors WHERE is_favorite = true;
