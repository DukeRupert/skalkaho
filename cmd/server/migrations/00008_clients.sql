-- +goose Up
CREATE TABLE clients (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    company TEXT,
    email TEXT,
    phone TEXT,
    address TEXT,
    city TEXT,
    state TEXT,
    zip TEXT,
    tax_id TEXT,
    notes TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX idx_clients_name ON clients(name);
CREATE INDEX idx_clients_created_at ON clients(created_at);

-- Add client_id to jobs
ALTER TABLE jobs ADD COLUMN client_id TEXT REFERENCES clients(id) ON DELETE SET NULL;
CREATE INDEX idx_jobs_client ON jobs(client_id);

-- +goose Down
DROP INDEX IF EXISTS idx_jobs_client;
DROP INDEX IF EXISTS idx_clients_created_at;
DROP INDEX IF EXISTS idx_clients_name;
DROP TABLE IF EXISTS clients;
