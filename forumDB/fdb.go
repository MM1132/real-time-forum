package forumDB

import (
	"database/sql"
	util "forum/utils"
	"log"
	"os"
)

// Initialize the database
func InitializeDB() *sql.DB {
	// Check if the database exists
	exists := true
	if _, err := os.Stat("db/forum.db"); os.IsNotExist(err) {
		exists = false
	}

	db, err := sql.Open("sqlite3", "db/forum.db")
	util.FatalErr(err)

	if !exists {
		log.Println("Database doesn't exist, initializing it now...")

		execSqlFile(db, "db/tables.sql")
		execSqlFile(db, "db/init.sql")
		execSqlFile(db, "db/initTestData.sql") //! Delete this before submitting the project

		log.Println("Database initialized!")
	}

	return db
}

// Create the initial tables
func execSqlFile(db *sql.DB, path string) {
	sqlfile, err := os.ReadFile(path)
	util.FatalErr(err)

	_, err = db.Exec(string(sqlfile))
	util.FatalErr(err)
}
