package forumDB

import (
	"database/sql"
	"forum/utils"
	"time"
)

type User struct {
	UserID   int
	Name     string
	Email    string
	Password string
	Creation time.Time
}

type UserInterface interface {
	Insert(newUser User) (int, error)
	Get(UID int) (User, error)
	ByName(name string) (User, error)
}

type UserModel struct {
	DB *sql.DB
}

// Insert a user into db, returns the UID of the newly inserted user
func (m UserModel) Insert(newUser User) (int, error) {
	stmt, err := m.DB.Prepare(
		"INSERT INTO users(name, email, password, created) values(?,?,?,?)",
	)
	utils.FatalErr(err)

	res, err := stmt.Exec(newUser.Name, newUser.Email, newUser.Password, time.Now())
	if err != nil {
		return 0, err
	}

	uid, err := res.LastInsertId()
	utils.FatalErr(err)

	return int(uid), nil
}

// Get a user by UID, returns sql.ErrNoRows if not found
func (m UserModel) Get(UID int) (User, error) {
	stmt, err := m.DB.Prepare(
		"SELECT * FROM users WHERE userID=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(UID)
	user := User{}
	err = row.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Creation)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (m UserModel) ByName(name string) (User, error) {
	stmt, err := m.DB.Prepare(
		"SELECT * FROM users WHERE name=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(name)
	user := User{}
	err = row.Scan(&user.UserID, &user.Name, &user.Email, &user.Password, &user.Creation)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
