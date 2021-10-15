package forumDB

import (
	"database/sql"
	utils "forum/utils"
	"time"
)

type User struct {
	UID      int
	Name     string
	Email    string
	Password string
	Creation time.Time
}

// Insert a user into db, returns the UID of the newly inserted user
func InsertUser(db *sql.DB, newUser *User) (int, error) {
	stmt, err := db.Prepare(
		"INSERT INTO users(name, email, password, created) values(?,?,?,?)",
	)
	utils.FatalErr(err)

	res, err := stmt.Exec(newUser.Name, newUser.Email, newUser.Password, time.Now())
	if err != nil {
		return 0, err
	}

	uid, err := res.LastInsertId()
	utils.FatalErr(err)

	return int(uid), err
}

// Get a user by UID, returns sql.ErrNoRows if not found
func GetUser(db *sql.DB, UID int) (*User, error) {
	stmt, err := db.Prepare(
		"SELECT * FROM users WHERE uid=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(UID)
	user := &User{}
	err = row.Scan(&user.UID, &user.Name, &user.Email, &user.Password, &user.Creation)
	if err != nil {
		return nil, err
	}

	return user, nil
}
