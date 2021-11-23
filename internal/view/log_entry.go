package view

import (
	"coffee-log/db/sqlc"
	"fmt"
)

type LogEntryView struct {
	sqlc.LogEntry
	URL string
	JustAdded bool
}

func NewLogEntryView(entry sqlc.LogEntry, logSlug string, justAdded bool) LogEntryView {
	url := fmt.Sprintf("/logs/%s/entries/%d", logSlug, entry.ID)
	return LogEntryView{
		LogEntry: entry,
		JustAdded: justAdded,
		URL: url,
	}
}
