package forumDB

import (
	"database/sql"
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
	fatalErr(err)

	res, err := stmt.Exec(newUser.Name, newUser.Email, newUser.Password, time.Now())
	if err != nil {
		return 0, err
	}

	uid, err := res.LastInsertId()
	fatalErr(err)

	return int(uid), err
}

// Get a user by UID, returns sql.ErrNoRows if not found
func GetUser(db *sql.DB, UID int) (*User, error) {
	stmt, err := db.Prepare(
		"SELECT * FROM users WHERE uid=?",
	)
	fatalErr(err)

	row := stmt.QueryRow(UID)
	user := &User{}
	err = row.Scan(&user.UID, &user.Name, &user.Email, &user.Password, &user.Creation)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByName(db *sql.DB, name string) (*User, error) {
	stmt, err := db.Prepare(
		"SELECT * FROM users WHERE name=?",
	)
	fatalErr(err)

	row := stmt.QueryRow(name)
	user := &User{}
	err = row.Scan(&user.UID, &user.Name, &user.Email, &user.Password, &user.Creation)
	if err != nil {
		return nil, err
	}

	return user, nil
}
