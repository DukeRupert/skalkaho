-- +goose Up
-- Add status and expires_at fields to jobs
ALTER TABLE jobs ADD COLUMN status TEXT NOT NULL DEFAULT 'draft'
    CHECK (status IN ('draft', 'sent', 'accepted', 'rejected', 'expired'));
ALTER TABLE jobs ADD COLUMN expires_at TEXT;

-- Add indexes for filtering/sorting
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_created_at ON jobs(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_jobs_created_at;
DROP INDEX IF EXISTS idx_jobs_status;
ALTER TABLE jobs DROP COLUMN expires_at;
ALTER TABLE jobs DROP COLUMN status;
