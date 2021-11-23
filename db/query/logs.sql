-- name: ListLogs :many
SELECT *
FROM logs
ORDER BY id ASC;

-- name: GetLog :one
SELECT *
FROM logs
WHERE id = $1
LIMIT 1;

-- name: GetLogBySlug :one
SELECT *
FROM logs
WHERE slug = $1
LIMIT 1;

-- name: GetLogByUserId :one
SELECT *
FROM logs
WHERE user_id = $1
LIMIT 1;

-- name: CreateLog :one
INSERT INTO logs (user_id, slug, title)
VALUES ($1, $2, $3)
RETURNING *;
