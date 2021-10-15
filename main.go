package main

import (
	"fmt"
	fdb "forum/forumDB"
	utils "forum/utils"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	fmt.Println("I am running")
	// Initialize sql.DB struct
	db := fdb.InitializeDB()

	// Example on inserting a new user into the DB...
	newUser := fdb.User{Name: "Raigo", Email: "krisimegaemail@gmail.com", Password: "securepassword"}
	uid, err := fdb.InsertUser(db, &newUser)
	utils.CheckErr(err)

	// ...and then getting that same user from the DB
	foundUser, err := fdb.GetUser(db, uid)
	utils.CheckErr(err)

	fmt.Printf("%+v\n", foundUser)

	db.Close()
}

// Helper function to print errors
