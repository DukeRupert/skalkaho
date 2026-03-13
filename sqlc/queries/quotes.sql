-- name: CreateQuote :one
INSERT INTO quotes (id, project_id, version, status, estimate_snapshot, totals_snapshot, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetQuote :one
SELECT * FROM quotes WHERE id = $1;

-- name: GetQuoteByToken :one
SELECT q.*,
    p.name AS project_name, p.client_name AS project_client_name,
    c.company_name AS client_company, c.contact_name AS client_contact,
    c.email AS client_email, c.phone AS client_phone
FROM quotes q
JOIN projects p ON p.id = q.project_id
LEFT JOIN clients c ON c.id = p.client_id
WHERE q.token = $1;

-- name: ListQuotesByProject :many
SELECT * FROM quotes WHERE project_id = $1 ORDER BY version DESC;

-- name: UpdateQuoteStatus :exec
UPDATE quotes SET status = $2 WHERE id = $1;

-- name: UpdateQuoteSent :exec
UPDATE quotes SET token = $2, sent_at = $3, expires_at = $4, status = 'sent' WHERE id = $1;

-- name: GetLatestQuoteVersion :one
SELECT COALESCE(MAX(version), 0)::integer AS version FROM quotes WHERE project_id = $1;

-- name: SupersedeActiveQuotes :exec
UPDATE quotes SET status = 'superseded' WHERE project_id = $1 AND status IN ('draft', 'sent');

-- name: CreateQuoteSignature :one
INSERT INTO quote_signatures (id, quote_id, signer_name, signer_ip)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetQuoteSignature :one
SELECT * FROM quote_signatures WHERE quote_id = $1;

-- name: CreateQuoteEmail :one
INSERT INTO quote_emails (id, quote_id, recipient, provider_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: ListQuoteEmails :many
SELECT * FROM quote_emails WHERE quote_id = $1 ORDER BY sent_at DESC;
