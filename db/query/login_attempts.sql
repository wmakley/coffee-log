-- name: GetLoginAttempt :one
SELECT * FROM login_attempts
WHERE ip_address = $1 LIMIT 1;

-- name: CreateLoginAttempt :one
INSERT INTO login_attempts (ip_address, attempts, created_at, updated_at)
VALUES ($1, 1, NOW(), NOW())
RETURNING *;

-- name: IncrementLoginAttempt :one
UPDATE login_attempts
SET attempts = attempts + 1, updated_at = NOW()
WHERE ip_address = $1
RETURNING *;

-- name: GetBannedIP :one
SELECT * FROM banned_ips
WHERE ip_address = $1
LIMIT 1;

-- name: CreateBannedIP :one
INSERT INTO banned_ips (ip_address, created_at)
VALUES ($1, NOW())
RETURNING *;

-- name: DeleteAllLoginAttempts :exec
DELETE FROM login_attempts;

-- name: DeleteAllBannedIPs :exec
DELETE FROM banned_ips;
