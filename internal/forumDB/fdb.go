package forumDB

import (
	"database/sql"
	util "forum/internal/utils"
	"log"
	"os"
)

// Initialize the database
func OpenDB(dbPath string) *sql.DB {
	// Check if the database exists
	exists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		exists = false
	}

	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	util.FatalErr(err)

	if !exists {
		log.Println("Database doesn't exist, initializing it now...")
		initDB(db, dbPath)
		log.Println("Database initialized!")
	}

	return db
}

func initDB(db *sql.DB, dbPath string) {
	pre := "server/db/sql/init/"
	err := execSqlFiles(db,
		pre+"tables.sql",
		pre+"init.sql",
		pre+"initTestUsers.sql",
		pre+"initTestBoards.sql",
		pre+"initTestThreads.sql",
		pre+"initTestPosts.sql",
	)
	if err != nil {
		log.Println(err)
		log.Println("Failed to initalize database :(")

		log.Print("Attempting to delete...")
		closeErr := db.Close()
		delErr := os.Remove(dbPath)

		if closeErr != nil || delErr != nil {
			log.Println("Failed, please delete the database manually!!!")
			log.Printf("DB Close Error: %v\n", closeErr)
			log.Printf("DB Delete Error: %v\n", delErr)
		} else {
			log.Println("Database deleted")
		}

		os.Exit(1)
	}
}

func execSqlFiles(db *sql.DB, paths ...string) error {
	for _, path := range paths {
		if err := execSqlFile(db, path); err != nil {
			return err
		}
	}

	return nil
}

func execSqlFile(db *sql.DB, path string) error {
	sqlfile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(sqlfile))
	if err != nil {
		return err
	}

	return nil
}
