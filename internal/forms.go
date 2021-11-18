package internal

import (
	"database/sql"
	"fmt"
	"strings"
)

type NewLogEntryForm struct {
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
