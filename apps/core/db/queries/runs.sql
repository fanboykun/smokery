-- name: CreateRun :one
INSERT INTO runs (project_id, plan_id, status) VALUES ($1, $2, 'pending') RETURNING *;

-- name: GetRun :one
SELECT * FROM runs WHERE id = $1;

-- name: UpdateRunStatus :one
UPDATE runs SET status = $2, started_at = $3, finished_at = $4 WHERE id = $1 RETURNING *;

-- name: ListRunsByProject :many
SELECT * FROM runs WHERE project_id = $1 ORDER BY created_at DESC;

-- name: CreateRunResult :one
INSERT INTO run_results (run_id, result) VALUES ($1, $2) RETURNING *;

-- name: GetRunResult :one
SELECT * FROM run_results WHERE run_id = $1;
