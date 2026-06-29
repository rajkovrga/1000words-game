package database

import (
	"database/sql"
	"log/slog"

	_ "modernc.org/sqlite"
)

func Connect(dbPath string, logger *slog.Logger) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	logger.Info("Database connected", "path", dbPath)

	return db, nil
}
