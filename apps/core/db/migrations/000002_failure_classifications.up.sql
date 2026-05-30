CREATE TABLE IF NOT EXISTS failure_classifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    run_id UUID NOT NULL REFERENCES runs(id) ON DELETE CASCADE,
    classification TEXT NOT NULL,
    assignee TEXT NOT NULL DEFAULT '',
    note TEXT NOT NULL DEFAULT '',
    author TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_failure_classifications_run_id ON failure_classifications(run_id);
