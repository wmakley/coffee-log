package view

import "coffee-log/db/sqlc"

type LogEntryView struct {
	sqlc.LogEntry
	JustAdded bool
}

func NewLogEntryView(entry sqlc.LogEntry, justAdded bool) LogEntryView {
	return LogEntryView{
		LogEntry: entry,
		JustAdded: justAdded,
	}
}
