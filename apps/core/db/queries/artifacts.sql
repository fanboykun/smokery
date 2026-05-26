-- name: CreateArtifact :one
INSERT INTO artifacts (run_id, type, path) VALUES ($1, $2, $3) RETURNING *;

-- name: ListArtifactsByRun :many
SELECT * FROM artifacts WHERE run_id = $1 ORDER BY created_at;
