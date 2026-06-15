package persistence

import (
	"database/sql"
	"fmt"
	"personal_bot/backend/internal/core/constants"

	_ "modernc.org/sqlite"
)

func NewConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite", constants.Database)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
