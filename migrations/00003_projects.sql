-- +goose Up
CREATE TABLE projects (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    client_id TEXT REFERENCES clients(id) ON DELETE SET NULL,
    client_name TEXT,
    status TEXT NOT NULL DEFAULT 'Draft' CHECK (status IN ('Draft', 'In Review', 'Active', 'Completed')),
    total REAL NOT NULL DEFAULT 0,
    description TEXT,
    materials_markup REAL NOT NULL DEFAULT 20,
    labor_markup REAL NOT NULL DEFAULT 25,
    equipment_markup REAL NOT NULL DEFAULT 15,
    subs_markup REAL NOT NULL DEFAULT 10,
    other_markup REAL NOT NULL DEFAULT 10,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_projects_client_id ON projects(client_id);
CREATE INDEX idx_projects_status ON projects(status);

-- +goose Down
DROP TABLE IF EXISTS projects;
