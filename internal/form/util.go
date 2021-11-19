package form

import (
	"database/sql"
	"strings"
)

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
