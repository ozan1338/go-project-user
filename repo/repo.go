package repo

import "database/sql"

type repo struct {
	db *sql.DB
}