-- Signature Requests

-- name: CreateSignatureRequest :one
INSERT INTO signature_requests (
    id, estimate_id, recipient_email, recipient_name, token,
    document_hash, quote_snapshot, message, status, expires_at,
    sender_ip, sender_user_agent
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: GetSignatureRequest :one
SELECT * FROM signature_requests WHERE id = $1;

-- name: GetSignatureRequestByToken :one
SELECT * FROM signature_requests WHERE token = $1;

-- name: GetSignatureRequestByEstimate :one
SELECT * FROM signature_requests
WHERE estimate_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetPendingSignatureRequestByEstimate :one
SELECT * FROM signature_requests
WHERE estimate_id = $1 AND status = 'pending'
ORDER BY created_at DESC
LIMIT 1;

-- name: ListSignatureRequestsByEstimate :many
SELECT * FROM signature_requests
WHERE estimate_id = $1
ORDER BY created_at DESC;

-- name: UpdateSignatureRequestStatus :one
UPDATE signature_requests SET status = $1
WHERE id = $2
RETURNING *;

-- name: CancelSignatureRequest :one
UPDATE signature_requests SET status = 'cancelled'
WHERE id = $1 AND status = 'pending'
RETURNING *;

-- name: ExpireOldSignatureRequests :exec
UPDATE signature_requests SET status = 'expired'
WHERE status = 'pending' AND expires_at < datetime('now');

-- Signatures

-- name: CreateSignature :one
INSERT INTO signatures (
    id, request_id, legal_name, consent_text, document_hash,
    signed_at, signer_ip, signer_user_agent, signer_email,
    certificate_pdf_path
)
VALUES ($1, $2, $3, $4, $5, datetime('now'), $6, $7, $8, $9)
RETURNING *;

-- name: GetSignature :one
SELECT * FROM signatures WHERE id = $1;

-- name: GetSignatureByRequest :one
SELECT * FROM signatures WHERE request_id = $1;

-- name: UpdateSignatureCertificatePath :one
UPDATE signatures SET certificate_pdf_path = $1
WHERE id = $2
RETURNING *;

-- Check if estimate has any pending signature requests (for locking)
-- name: HasPendingSignatureRequest :one
SELECT EXISTS(
    SELECT 1 FROM signature_requests
    WHERE estimate_id = $1 AND status = 'pending'
) as has_pending;

-- Check if estimate has been signed
-- name: HasSignedRequest :one
SELECT EXISTS(
    SELECT 1 FROM signature_requests sr
    INNER JOIN signatures s ON s.request_id = sr.id
    WHERE sr.estimate_id = $1
) as has_signed;
