-- QBO Connection ---------------------------------------------------------

-- name: GetQBConnection :one
SELECT * FROM qb_connection WHERE id = 1;

-- name: UpsertQBConnection :one
INSERT INTO qb_connection (id, realm_id, access_token, refresh_token,
    access_token_expiry, refresh_token_expiry)
VALUES (1, $1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE SET
    realm_id              = EXCLUDED.realm_id,
    access_token          = EXCLUDED.access_token,
    refresh_token         = EXCLUDED.refresh_token,
    access_token_expiry   = EXCLUDED.access_token_expiry,
    refresh_token_expiry  = EXCLUDED.refresh_token_expiry,
    last_refreshed_at     = NOW()
RETURNING *;

-- name: UpdateQBTokens :one
UPDATE qb_connection
SET access_token         = $1,
    refresh_token        = $2,
    access_token_expiry  = $3,
    refresh_token_expiry = $4,
    last_refreshed_at    = NOW()
WHERE id = 1
RETURNING *;

-- name: DeleteQBConnection :exec
DELETE FROM qb_connection WHERE id = 1;

-- Client ↔ QBO Linking ---------------------------------------------------

-- name: LinkClientToQB :one
UPDATE clients
SET qb_customer_id = $1, qb_customer_name = $2
WHERE id = $3
RETURNING *;

-- name: UnlinkClientFromQB :one
UPDATE clients
SET qb_customer_id = NULL, qb_customer_name = NULL
WHERE id = $1
RETURNING *;

-- name: ListUnlinkedClients :many
SELECT * FROM clients
WHERE qb_customer_id IS NULL
ORDER BY company_name ASC;

-- name: ListLinkedClients :many
SELECT * FROM clients
WHERE qb_customer_id IS NOT NULL
ORDER BY company_name ASC;

-- Invoices ---------------------------------------------------------------

-- name: CreateInvoice :one
INSERT INTO invoices (id, quote_id, amount, due_date, memo)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetInvoice :one
SELECT * FROM invoices WHERE id = $1;

-- name: GetInvoiceByQBId :one
SELECT * FROM invoices WHERE qb_invoice_id = $1;

-- name: ListInvoicesByQuote :many
SELECT i.*, COALESCE(SUM(p.amount), 0)::real AS total_paid
FROM invoices i
LEFT JOIN invoice_payments p ON p.invoice_id = i.id
WHERE i.quote_id = $1
GROUP BY i.id
ORDER BY i.created_at ASC;

-- name: MarkInvoicePushed :one
UPDATE invoices
SET qb_invoice_id  = $1,
    qb_invoice_num = $2,
    qb_sync_token  = $3,
    status         = 'pushed',
    updated_at     = NOW()
WHERE id = $4
RETURNING *;

-- name: UpdateInvoiceStatus :one
UPDATE invoices
SET status = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: UpdateInvoiceSyncToken :one
UPDATE invoices
SET qb_sync_token = $1, updated_at = NOW()
WHERE id = $2
RETURNING *;

-- name: SumInvoicedByQuote :one
SELECT COALESCE(SUM(amount), 0)::real AS total_invoiced
FROM invoices
WHERE quote_id = $1
  AND status != 'voided';

-- Payments ---------------------------------------------------------------

-- name: CreateInvoicePayment :one
INSERT INTO invoice_payments (id, invoice_id, qb_payment_id, amount, payment_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetInvoicePaymentByQBId :one
SELECT * FROM invoice_payments WHERE qb_payment_id = $1;

-- name: ListPaymentsByInvoice :many
SELECT * FROM invoice_payments
WHERE invoice_id = $1
ORDER BY payment_date ASC;

-- name: SumPaymentsByInvoice :one
SELECT COALESCE(SUM(amount), 0)::real AS total_paid
FROM invoice_payments
WHERE invoice_id = $1;
