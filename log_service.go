package main

import (
	"coffee-log/queries"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type ErrInvalidForm struct {
	ValidationErrors map[string]string
}

func (e *ErrInvalidForm) Error() string {
	return fmt.Sprintf("invalid form: %+v", e.ValidationErrors)
}

type LogService struct {
	db *sql.DB
	q *queries.Queries
}

func NewLogService(db *sql.DB) *LogService {
	return &LogService{
		db: db,
		q: queries.New(db),
	}
}

func (s LogService) ListLogs(ctx context.Context) ([]queries.Log, error) {
	return s.q.ListLogs(ctx)
}

func (s LogService) GetLogBySlug(ctx context.Context, slug string) (queries.Log, error) {
	results, err := s.q.GetLogBySlug(ctx, slug)
	return results, liftError(err)
}

func (s LogService) GetLogAndEntriesBySlugOrderByDateDesc(
	ctx context.Context, slug string,
) (queries.Log, []queries.Entry, error) {
	var log2 queries.Log
	var entries []queries.Entry

	err := transaction(ctx, s.db, func(tx *sql.Tx) error {
		q := queries.New(tx)
		var err error

		log2, err = q.GetLogBySlug(ctx, slug)
		if err != nil {
			return err
		}

		entries, err = q.ListLogEntriesByLogIDOrderByDateDesc(ctx, log2.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return log2, entries, liftError(err)
}

func (s LogService) GetLogAndLastEntry(ctx context.Context, logID string) (queries.Log, queries.Entry, error) {
	var log2 queries.Log
	var lastEntry queries.Entry

	err := transaction(ctx, s.db, func(tx *sql.Tx) error {
		q := queries.New(tx)
		var err error

		log2, err = q.GetLogBySlug(ctx, logID)
		if err != nil {
			return err
		}

		lastEntry, err = q.GetLastLogEntryByLogID(ctx, log2.ID)
		if err != nil {
			return err
		}

		return nil
	})

	return log2, lastEntry, liftError(err)
}

func (s LogService) CreateLogEntry(ctx context.Context, logID int64, form NewLogEntryForm) (queries.Entry, error) {
	validationErrors := make(map[string]string)

	coffeeGrams, err := parseNullInt32(form.CoffeeGrams)
	if err != nil {
		validationErrors["CoffeeGrams"] = "is not an integer"
	}
	waterGrams, err := parseNullInt32(form.WaterGrams)
	if err != nil {
		validationErrors["WaterGrams"] = "is not an integer"
	}

	params := queries.CreateLogEntryParams{
		LogID:       logID,
		Coffee:      strings.TrimSpace(form.Coffee),
		Water:       blankToNullString(form.Water),
		Method:      blankToNullString(form.Method),
		Grind:       blankToNullString(form.Grind),
		Tasting:     blankToNullString(form.Tasting),
		AddlNotes:   blankToNullString(form.AddlNotes),
		CoffeeGrams: coffeeGrams,
		WaterGrams:  waterGrams,
	}

	if params.Coffee == "" {
		validationErrors["Coffee"] = "must not be blank"
	}

	if len(validationErrors) > 0 {
		return queries.Entry{}, &ErrInvalidForm{ValidationErrors: validationErrors}
	}

	entry, err := s.q.CreateLogEntry(ctx, params)
	return entry, liftError(err)
}

func parseNullInt32(input string) (sql.NullInt32, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return sql.NullInt32{
			Int32: 0,
			Valid: false,
		}, nil
	}

	parsed, err := strconv.ParseInt(trimmed, 10, 32)
	if err != nil {
		return sql.NullInt32{}, err
	}

	return sql.NullInt32{
		Int32: int32(parsed),
		Valid: true,
	}, nil
}

func blankToNullString(input string) sql.NullString {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: trimmed,
		Valid:  true,
	}
}

func transaction(ctx context.Context, db *sql.DB, fn func (tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				panic(txErr)
			}
			return
		}
		if txErr := tx.Commit(); txErr != nil {
			panic(txErr)
		}
	}()

	err = fn(tx)
	return err
}

// Map infrastructure error to a service-level error
func liftError(err error) error {
	if err == nil {
		return err
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrRecordNotFound
	}

	return err
}
