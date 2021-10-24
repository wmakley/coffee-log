-- name: ListLogs :many
SELECT * FROM logs
ORDER BY name ASC;

-- name: GetLog :one
SELECT * FROM logs
WHERE id = $1 LIMIT 1;

-- name: GetLogBySlug :one
SELECT * FROM logs
WHERE slug = $1 LIMIT 1;

-- name: ListLogEntriesByLogIdOrderByDateDesc :many
SELECT * FROM entries
WHERE log_id = $1
ORDER BY created_at DESC;
