package forumDB

import (
	"database/sql"
	"forum/internal/utils"
	"time"
)

type User struct {
	UserID   int
	Name     string
	Email    string
	Password string
	Creation time.Time
}

type UserModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewUserModel(db *sql.DB) UserModel {
	statements := make(map[string]*sql.Stmt)
	model := UserModel{db: db}

	var err error
	statements["Insert"], err = db.Prepare("INSERT INTO users(name, email, password, created) values(?,?,?,?)")
	utils.FatalErr(err)

	statements["Get"], err = db.Prepare("SELECT * FROM users WHERE userID=?")
	utils.FatalErr(err)

	statements["ByName"], err = db.Prepare("SELECT * FROM users WHERE name=?")
	utils.FatalErr(err)

	model.statements = statements
	return model
}

// Insert a user into db, returns the UID of the newly inserted user
func (m UserModel) Insert(newUser User) (int, error) {
	stmt := m.statements["Insert"]
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
	stmt := m.statements["Get"]
	row := stmt.QueryRow(UID)

	user := User{}
	if err := row.Scan(
		&user.UserID, &user.Name, &user.Email, &user.Password, &user.Creation,
	); err != nil {
		return User{}, err
	}

	return user, nil
}

func (m UserModel) ByName(name string) (User, error) {
	stmt := m.statements["ByName"]
	row := stmt.QueryRow(name)

	user := User{}
	if err := row.Scan(
		&user.UserID, &user.Name, &user.Email, &user.Password, &user.Creation,
	); err != nil {
		return User{}, err
	}

	return user, nil
}
