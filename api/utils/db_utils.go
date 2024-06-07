package utils

import (
	"database/sql"
)

func DbExists(db *sql.DB, query string, params ...any) bool {
	var column any

	if err := db.QueryRow(query, params...).Scan(&column); err != nil {
		if err == sql.ErrNoRows {
			return false
		}

		return false
	}

	return true
}
