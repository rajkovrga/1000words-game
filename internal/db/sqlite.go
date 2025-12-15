package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect(dbPath string) {
	database, err := sql.Open("sqlite", dbPath)

	if err != nil {
		log.Fatal("Db connection error", err)
	}

	if err = database.Ping(); err != nil {
		log.Fatal("Db ping error", err)
	}

	DB = database

}

func ExecuteQuery(query string) {

}
