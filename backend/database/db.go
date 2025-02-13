package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbname string) (err error) {
	DB, err = sql.Open("sqlite3", dbname)
	if err != nil {
		return
	}
	// Create tables
	err = CreateTables(TableCreationStatements)
	if err != nil {
		return
	}
	log.Println("Database initialized successfully")
	return
}
