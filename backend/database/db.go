package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() (err error) {
	// Create database directory if it doesn't exist
	dbDir := "./data"
	if err = os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %v", err)
	}

	dbPath := filepath.Join(dbDir, "forum.db")
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	// Create tables
	err = CreateTables(TableCreationStatements)
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}
	log.Println("Database initialized successfully")
	return
}
