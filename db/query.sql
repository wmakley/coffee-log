-- name: ListLogs :many
SELECT * FROM logs
ORDER BY name ASC;

-- name: GetLog :one
SELECT * FROM logs
WHERE id = $1 LIMIT 1;

-- name: GetLogBySlug :one
SELECT * FROM logs
WHERE slug = $1 LIMIT 1;

-- name: ListLogEntriesByLogIDOrderByDateDesc :many
SELECT * FROM entries
WHERE log_id = $1
ORDER BY created_at DESC;

-- name: GetLastLogEntryByLogID :one
SELECT * FROM entries
WHERE log_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: CreateLogEntry :one
INSERT INTO entries
    (log_id, coffee, water, method, grind, tasting, addl_notes, coffee_grams, water_grams)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING *;

-- name: RecordLogEntryHistory :exec
INSERT INTO entries_history
(action, stamp, log_id, coffee, water, method, grind, tasting, addl_notes, coffee_grams, water_grams)
VALUES ($1, NOW(), $2, $3, $4, $5, $6, $7, $8, $9, $10);
