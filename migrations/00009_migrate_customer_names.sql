-- +goose Up
-- Create clients from unique customer_names
INSERT INTO clients (id, name, created_at)
SELECT
    lower(hex(randomblob(4)) || '-' || hex(randomblob(2)) || '-4' ||
          substr(hex(randomblob(2)),2) || '-' ||
          substr('89ab', abs(random()) % 4 + 1, 1) ||
          substr(hex(randomblob(2)),2) || '-' || hex(randomblob(6))),
    customer_name,
    datetime('now')
FROM jobs
WHERE customer_name IS NOT NULL AND customer_name != ''
GROUP BY customer_name;

-- Link jobs to their new client records
UPDATE jobs
SET client_id = (SELECT id FROM clients WHERE clients.name = jobs.customer_name)
WHERE customer_name IS NOT NULL AND customer_name != '';

-- +goose Down
UPDATE jobs SET client_id = NULL;
DELETE FROM clients;
