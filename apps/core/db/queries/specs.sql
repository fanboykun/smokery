-- name: CreateSpec :one
INSERT INTO specs (project_id, version, title, raw, analysis) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetSpec :one
SELECT * FROM specs WHERE id = $1;

-- name: ListSpecsByProject :many
SELECT * FROM specs WHERE project_id = $1 ORDER BY created_at DESC;
