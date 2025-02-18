package queries

import (
	"log"
	"os"
	"testing"

	"forum/database"
)

var globaldbname string = "test_forum.db"

func deleteTestDb(dbname string) {
	if _, err := os.Stat(dbname); err == nil {
		err = os.RemoveAll(dbname)
		if err != nil {
			log.Println("Error removing test database")
		}
	}
}

func TestMain(m *testing.M) {
	deleteTestDb(globaldbname)
	database.InitDB(globaldbname)
	defer database.DB.Close()
	code := m.Run()
	os.Exit(code)
}
