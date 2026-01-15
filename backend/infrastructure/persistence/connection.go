package persistence

import (
	"database/sql"
	"personal_bot/internal/core/constants"

	_ "modernc.org/sqlite"
)

func NewConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", constants.Database)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
