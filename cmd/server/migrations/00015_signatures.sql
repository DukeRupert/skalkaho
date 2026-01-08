-- +goose Up

-- Company profile: contractor/business information
CREATE TABLE company_profile (
    id TEXT PRIMARY KEY DEFAULT 'default',
    name TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    address TEXT,
    city TEXT,
    state TEXT,
    zip TEXT,
    logo_path TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Insert default row
INSERT INTO company_profile (id, name) VALUES ('default', 'Your Company Name');

-- Signature requests: tracks outbound signature requests
CREATE TABLE signature_requests (
    id TEXT PRIMARY KEY,
    estimate_id TEXT NOT NULL REFERENCES estimates(id) ON DELETE CASCADE,

    -- Recipient info
    recipient_email TEXT NOT NULL,
    recipient_name TEXT NOT NULL,

    -- Security
    token TEXT NOT NULL UNIQUE,

    -- Document snapshot (immutable after sending)
    document_hash TEXT NOT NULL,
    quote_snapshot TEXT NOT NULL,

    -- Optional message from sender
    message TEXT,

    -- Lifecycle
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'signed', 'expired', 'cancelled')),
    expires_at TEXT NOT NULL,

    -- Sender context (for audit)
    sender_ip TEXT,
    sender_user_agent TEXT,

    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_signature_requests_token ON signature_requests(token);
CREATE INDEX idx_signature_requests_estimate ON signature_requests(estimate_id);
CREATE INDEX idx_signature_requests_status ON signature_requests(status);

-- Signatures: captures the actual signature event (immutable)
CREATE TABLE signatures (
    id TEXT PRIMARY KEY,
    request_id TEXT NOT NULL REFERENCES signature_requests(id) ON DELETE RESTRICT,

    -- What they signed
    legal_name TEXT NOT NULL,
    consent_text TEXT NOT NULL,
    document_hash TEXT NOT NULL,

    -- Audit trail
    signed_at TEXT NOT NULL DEFAULT (datetime('now')),
    signer_ip TEXT NOT NULL,
    signer_user_agent TEXT NOT NULL,
    signer_email TEXT NOT NULL,

    -- Generated artifacts
    certificate_pdf_path TEXT,

    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_signatures_request ON signatures(request_id);

-- +goose Down
DROP INDEX IF EXISTS idx_signatures_request;
DROP TABLE IF EXISTS signatures;
DROP INDEX IF EXISTS idx_signature_requests_status;
DROP INDEX IF EXISTS idx_signature_requests_estimate;
DROP INDEX IF EXISTS idx_signature_requests_token;
DROP TABLE IF EXISTS signature_requests;
DROP TABLE IF EXISTS company_profile;
