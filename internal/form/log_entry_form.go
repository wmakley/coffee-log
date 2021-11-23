package form

import (
	"coffee-log/db/sqlc"
	"database/sql"
	"strings"
	"time"
)

type LogEntryForm struct {
	URL          string
	EntryDate    time.Time `form:"entry_date" binding:"required"`
	Coffee       string    `form:"coffee" binding:"required"`
	Water        string    `form:"water"`
	BrewMethod   string    `form:"brew_method"`
	CoffeeGrams  int32     `form:"coffee_grams"`
	WaterGrams   int32     `form:"water_grams"`
	GrindNotes   string    `form:"grind_notes"`
	TastingNotes string    `form:"tasting_notes"`
	AddlNotes    string    `form:"addl_notes"`
	Errors       map[string]string
}

func (form *LogEntryForm) IsValid() bool {
	form.Coffee = strings.TrimSpace(form.Coffee)
	form.Water = strings.TrimSpace(form.Water)
	form.BrewMethod = strings.TrimSpace(form.BrewMethod)
	form.GrindNotes = strings.TrimSpace(form.GrindNotes)
	form.TastingNotes = strings.TrimSpace(form.TastingNotes)
	form.AddlNotes = strings.TrimSpace(form.AddlNotes)
	form.Errors = make(map[string]string)

	if form.Coffee == "" {
		form.Errors["Coffee"] = "must not be blank"
	}

	return len(form.Errors) == 0
}

func (form *LogEntryForm) CreateParams() sqlc.CreateLogEntryParams {
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

func (form *LogEntryForm) UpdateParams() sqlc.CreateLogEntryParams {
	return sqlc.CreateLogEntryParams{
		LogID:        0,
		EntryDate:    form.EntryDate,
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
