// Code generated by sqlc. DO NOT EDIT.
// source: query.sql

package queries

import (
	"context"
	"database/sql"
)

const createLogEntry = `-- name: CreateLogEntry :one
INSERT INTO entries
    (log_id, coffee, water, method, grind, tasting, addl_notes, coffee_grams, water_grams)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING id, log_id, coffee, water, method, grind, tasting, addl_notes, coffee_grams, water_grams, created_at, updated_at
`

type CreateLogEntryParams struct {
	LogID       int64
	Coffee      string
	Water       sql.NullString
	Method      sql.NullString
	Grind       sql.NullString
	Tasting     sql.NullString
	AddlNotes   sql.NullString
	CoffeeGrams sql.NullInt32
	WaterGrams  sql.NullInt32
}

func (q *Queries) CreateLogEntry(ctx context.Context, arg CreateLogEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createLogEntry,
		arg.LogID,
		arg.Coffee,
		arg.Water,
		arg.Method,
		arg.Grind,
		arg.Tasting,
		arg.AddlNotes,
		arg.CoffeeGrams,
		arg.WaterGrams,
	)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.LogID,
		&i.Coffee,
		&i.Water,
		&i.Method,
		&i.Grind,
		&i.Tasting,
		&i.AddlNotes,
		&i.CoffeeGrams,
		&i.WaterGrams,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLastLogEntryByLogID = `-- name: GetLastLogEntryByLogID :one
SELECT id, log_id, coffee, water, method, grind, tasting, addl_notes, coffee_grams, water_grams, created_at, updated_at FROM entries
WHERE log_id = $1
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetLastLogEntryByLogID(ctx context.Context, logID int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getLastLogEntryByLogID, logID)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.LogID,
		&i.Coffee,
		&i.Water,
		&i.Method,
		&i.Grind,
		&i.Tasting,
		&i.AddlNotes,
		&i.CoffeeGrams,
		&i.WaterGrams,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLog = `-- name: GetLog :one
SELECT id, name, slug, created_at, updated_at FROM logs
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetLog(ctx context.Context, id int64) (Log, error) {
	row := q.db.QueryRowContext(ctx, getLog, id)
	var i Log
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLogBySlug = `-- name: GetLogBySlug :one
SELECT id, name, slug, created_at, updated_at FROM logs
WHERE slug = $1 LIMIT 1
`

func (q *Queries) GetLogBySlug(ctx context.Context, slug string) (Log, error) {
	row := q.db.QueryRowContext(ctx, getLogBySlug, slug)
	var i Log
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Slug,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listLogEntriesByLogIDOrderByDateDesc = `-- name: ListLogEntriesByLogIDOrderByDateDesc :many
SELECT id, log_id, coffee, water, method, grind, tasting, addl_notes, coffee_grams, water_grams, created_at, updated_at FROM entries
WHERE log_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListLogEntriesByLogIDOrderByDateDesc(ctx context.Context, logID int64) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listLogEntriesByLogIDOrderByDateDesc, logID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.LogID,
			&i.Coffee,
			&i.Water,
			&i.Method,
			&i.Grind,
			&i.Tasting,
			&i.AddlNotes,
			&i.CoffeeGrams,
			&i.WaterGrams,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listLogs = `-- name: ListLogs :many
SELECT id, name, slug, created_at, updated_at FROM logs
ORDER BY name ASC
`

func (q *Queries) ListLogs(ctx context.Context) ([]Log, error) {
	rows, err := q.db.QueryContext(ctx, listLogs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Log
	for rows.Next() {
		var i Log
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Slug,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const recordLogEntryHistory = `-- name: RecordLogEntryHistory :exec
INSERT INTO entries_history
(action, stamp, log_id, coffee, water, method, grind, tasting, addl_notes, coffee_grams, water_grams)
VALUES ($1, NOW(), $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

type RecordLogEntryHistoryParams struct {
	Action      string
	LogID       int64
	Coffee      sql.NullString
	Water       sql.NullString
	Method      sql.NullString
	Grind       sql.NullString
	Tasting     sql.NullString
	AddlNotes   sql.NullString
	CoffeeGrams sql.NullInt32
	WaterGrams  sql.NullInt32
}

func (q *Queries) RecordLogEntryHistory(ctx context.Context, arg RecordLogEntryHistoryParams) error {
	_, err := q.db.ExecContext(ctx, recordLogEntryHistory,
		arg.Action,
		arg.LogID,
		arg.Coffee,
		arg.Water,
		arg.Method,
		arg.Grind,
		arg.Tasting,
		arg.AddlNotes,
		arg.CoffeeGrams,
		arg.WaterGrams,
	)
	return err
}
