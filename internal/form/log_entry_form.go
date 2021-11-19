package form

import (
	"coffee-log/db/sqlc"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type LogEntryForm struct {
	LogID        int64  `form:"log_id"`
	Coffee       string `form:"coffee"`
	Water        string `form:"water"`
	BrewMethod   string `form:"brew_method"`
	CoffeeGrams  int32  `form:"coffee_grams"`
	WaterGrams   int32  `form:"water_grams"`
	GrindNotes   string `form:"grind_notes"`
	TastingNotes string `form:"tasting_notes"`
	AddlNotes    string `form:"addl_notes"`
}

type ErrInvalidForm struct {
	ValidationErrors map[string]string
}

func (e *ErrInvalidForm) Error() string {
	return fmt.Sprintf("invalid form: %+v", e.ValidationErrors)
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
		String: input,
		Valid:  true,
	}
}

func (form *LogEntryForm)CreateParams() sqlc.CreateLogEntryParams {
	return sqlc.CreateLogEntryParams{
		LogID:        0,
		EntryDate:    time.Now(),
		Coffee:       form.Coffee,
		Water:        blankToNullString(form.Water),
		BrewMethod:   blankToNullString(form.BrewMethod),
		GrindNotes:   blankToNullString(form.GrindNotes),
		TastingNotes: blankToNullString(form.TastingNotes),
		AddlNotes:    blankToNullString(form.AddlNotes),
		CoffeeGrams: sql.NullInt32{
			Int32: form.CoffeeGrams,
			Valid: form.CoffeeGrams > 0,
		},
		WaterGrams: sql.NullInt32{
			Int32: form.WaterGrams,
			Valid: form.WaterGrams > 0,
		},
	}
}
