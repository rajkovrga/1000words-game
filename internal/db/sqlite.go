package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func Connect(dbPath string) {
	db, err := sql.Open("sqlite", dbPath)

	if err != nil {
		log.Fatal("Db connection error", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Db ping error", err)
	}

}

func ExecuteQuery(query string) {

}
