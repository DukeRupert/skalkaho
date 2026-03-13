-- +goose Up
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
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS company_profile;
