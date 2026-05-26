-- name: CreateOperation :one
INSERT INTO operations (spec_id, operation_id, method, path, summary, tags, classification, is_destructive) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: ListOperationsBySpec :many
SELECT * FROM operations WHERE spec_id = $1 ORDER BY path, method;

-- name: UpdateOperationOverrides :one
UPDATE operations SET overrides = $2 WHERE id = $1 RETURNING *;

-- name: UpdateOperationClassification :one
UPDATE operations SET classification = $2, is_destructive = $3 WHERE id = $1 RETURNING *;
