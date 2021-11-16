package internal

import "fmt"

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
