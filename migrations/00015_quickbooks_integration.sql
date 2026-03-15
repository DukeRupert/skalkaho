-- +goose Up

-- OAuth tokens and connection state for the QBO integration.
-- Single-row table — one QBO connection per deployment.
CREATE TABLE qb_connection (
    id                   INTEGER PRIMARY KEY DEFAULT 1,
    realm_id             TEXT NOT NULL,
    access_token         TEXT NOT NULL,
    refresh_token        TEXT NOT NULL,
    access_token_expiry  TIMESTAMPTZ NOT NULL,
    refresh_token_expiry TIMESTAMPTZ NOT NULL,
    connected_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_refreshed_at    TIMESTAMPTZ,
    CHECK (id = 1)
);

-- QBO customer link on clients (manually set by owner)
ALTER TABLE clients
    ADD COLUMN qb_customer_id TEXT;

ALTER TABLE clients
    ADD COLUMN qb_customer_name TEXT;

-- Invoice tracking against signed quotes
CREATE TABLE invoices (
    id              TEXT PRIMARY KEY,
    quote_id        TEXT NOT NULL REFERENCES quotes(id) ON DELETE RESTRICT,
    qb_invoice_id   TEXT,
    qb_invoice_num  TEXT,
    qb_sync_token   TEXT,
    amount          REAL NOT NULL,
    status          TEXT NOT NULL DEFAULT 'draft'
                        CHECK (status IN ('draft', 'pushed', 'partially_paid', 'paid', 'voided')),
    due_date        DATE,
    memo            TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Payment records synced back from QBO
CREATE TABLE invoice_payments (
    id              TEXT PRIMARY KEY,
    invoice_id      TEXT NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    qb_payment_id   TEXT NOT NULL UNIQUE,
    amount          REAL NOT NULL,
    payment_date    DATE NOT NULL,
    synced_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for webhook lookups
CREATE INDEX idx_invoices_qb_invoice_id ON invoices(qb_invoice_id)
    WHERE qb_invoice_id IS NOT NULL;

CREATE INDEX idx_invoice_payments_qb_payment_id ON invoice_payments(qb_payment_id);

-- +goose Down

DROP INDEX IF EXISTS idx_invoice_payments_qb_payment_id;
DROP INDEX IF EXISTS idx_invoices_qb_invoice_id;
DROP TABLE IF EXISTS invoice_payments;
DROP TABLE IF EXISTS invoices;
ALTER TABLE clients DROP COLUMN IF EXISTS qb_customer_name;
ALTER TABLE clients DROP COLUMN IF EXISTS qb_customer_id;
DROP TABLE IF EXISTS qb_connection;
