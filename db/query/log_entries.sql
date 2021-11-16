-- name: ListLogEntriesByLogIDOrderByDateDesc :many
SELECT * FROM log_entries
WHERE deleted_at IS NOT NULL AND log_id = $1
ORDER BY entry_date DESC;

-- name: CreateLogEntry :one
INSERT INTO log_entries
(log_id, entry_date, coffee, water, brew_method, grind_notes, tasting_notes, addl_notes, coffee_grams, water_grams)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: UpdateLogEntry :one
UPDATE log_entries
SET entry_date = $2,
	coffee = $3,
	water = $4,
	brew_method = $5,
	grind_notes = $6,
	tasting_notes = $7,
	addl_notes = $8,
	coffee_grams = $9,
	water_grams = $10,
	updated_at = $11
WHERE id = $1
RETURNING *;
