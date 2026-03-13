-- +goose Up
CREATE TABLE quotes (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'sent', 'signed', 'expired', 'superseded')),
    estimate_snapshot JSONB,
    totals_snapshot JSONB,
    token TEXT UNIQUE,
    expires_at TIMESTAMPTZ,
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_by TEXT REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(project_id, version)
);

CREATE INDEX idx_quotes_project_id ON quotes(project_id);
CREATE INDEX idx_quotes_token ON quotes(token);

CREATE TABLE quote_signatures (
    id TEXT PRIMARY KEY,
    quote_id TEXT NOT NULL REFERENCES quotes(id) ON DELETE RESTRICT,
    signer_name TEXT NOT NULL,
    signer_ip TEXT,
    signed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_quote_signatures_quote_id ON quote_signatures(quote_id);

CREATE TABLE quote_emails (
    id TEXT PRIMARY KEY,
    quote_id TEXT NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    recipient TEXT NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    provider_id TEXT
);

CREATE INDEX idx_quote_emails_quote_id ON quote_emails(quote_id);

-- +goose Down
DROP TABLE IF EXISTS quote_emails;
DROP TABLE IF EXISTS quote_signatures;
DROP TABLE IF EXISTS quotes;
