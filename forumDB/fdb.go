package forumDB

import (
	"database/sql"
	utils "forum/utils"
	"os"
)

// Initialize the database
func InitializeDB() *sql.DB {
	// Check if the database doesn't exist
	db, err := sql.Open("sqlite3", "db/forum.db")
	utils.FatalErr(err)
	CreateTables(db)
	return db
}

// Create the initial tables
func CreateTables(db *sql.DB) {
	sqlfile, err := os.ReadFile("db/tables.sql")
	utils.FatalErr(err)
	commands := string(sqlfile)
	db.Exec(commands)
}
