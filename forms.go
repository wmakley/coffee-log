package main

type NewLogEntryForm struct {
	LogID       int64  `form:"log_id"`
	Coffee      string `form:"coffee"`
	Water		string `form:"water"`
	Method      string `form:"method"`
	Grind       string `form:"grind"`
	Tasting     string `form:"tasting"`
	AddlNotes   string `form:"addl_notes"`
	CoffeeGrams string `form:"coffee_grams"`
	WaterGrams  string `form:"water_grams"`
}
