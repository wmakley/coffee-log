package view

import (
	"coffee-log/db/sqlc"
	"fmt"
)

type LogTab struct {
	URL string
	Active bool
	Title string
}

func NewLogTab(log *sqlc.Log, active bool) LogTab {
	return LogTab{
		Title: log.Title,
		Active: active,
		URL: fmt.Sprintf("/logs/%s/entries/", log.Slug),
	}
}
