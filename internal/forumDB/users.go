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
	model := UserModel{db: db}

	model.statements = makeStatementMap(db, "server/db/sql/models/users.sql")

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

func (m UserModel) GetByName(name string) (User, error) {
	stmt := m.statements["GetByName"]
	row := stmt.QueryRow(name)

	user := User{}
	if err := row.Scan(
		&user.UserID, &user.Name, &user.Email, &user.Password, &user.Creation,
	); err != nil {
		return User{}, err
	}

	return user, nil
}
