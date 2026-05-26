-- name: CreateComment :one
INSERT INTO comments (run_id, author, body) VALUES ($1, $2, $3) RETURNING *;

-- name: ListCommentsByRun :many
SELECT * FROM comments WHERE run_id = $1 ORDER BY created_at ASC;
