package forumDB

import (
	"database/sql"
	"os"
)

// Initialize the database
func InitializeDB() *sql.DB {
	// Check if the database doesn't exist
	db, err := sql.Open("sqlite3", "db/forum.db")
	fatalErr(err)
	CreateTables(db)
	return db
}

// Create the initial tables
func CreateTables(db *sql.DB) {
	sqlfile, err := os.ReadFile("db/tables.sql")
	fatalErr(err)
	commands := string(sqlfile)
	_, err = db.Exec(commands)
	fatalErr(err)
}
