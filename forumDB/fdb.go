package forumDB

import (
	"database/sql"
)

// Initialize the database
func InitializeDB() *sql.DB {
	// Check if the database doesn't exist
	db, err := sql.Open("sqlite3", "db/forum.db")
	checkErr(err)
	CreateTables(db)
	return db
}

// Create the initial tables
func CreateTables(db *sql.DB) {
	CreateTable(db, "CREATE TABLE IF NOT EXISTS `users` ("+
		"`uid` INTEGER PRIMARY KEY AUTOINCREMENT,"+
		"`name` TEXT UNIQUE NOT NULL,"+
		"`email` TEXT UNIQUE NOT NULL,"+
		"`password` TEXT NOT NULL,"+
		"`created` DATE);")
}

// Create a new table. The string must be a valid sql command, otherwise it will fail.
func CreateTable(db *sql.DB, s string) {
	stmt, err := db.Prepare(s)
	checkErr(err)

	_, err = stmt.Exec()
	checkErr(err)
}
