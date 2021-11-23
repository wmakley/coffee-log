package form

import (
	"database/sql"
)

func blankToNullString(input string) sql.NullString {
	if input == "" {
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
