package forumDB

import (
	"database/sql"
	"forum/internal/utils"
	"time"
)

type User struct {
	UserID      int
	Name        string
	Email       string
	Password    string
	Image       string
	Description string
	Creation    time.Time

	Extras *UserExtras
}

type UserExtras struct {
	TotalPosts int
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
		&user.UserID, &user.Name, &user.Email, &user.Password, &user.Image, &user.Description, &user.Creation,
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
		&user.UserID, &user.Name, &user.Email, &user.Password, &user.Image, &user.Description, &user.Creation,
	); err != nil {
		return User{}, err
	}

	return user, nil
}

func (m UserModel) GetByEmail(email string) (User, error) {
	stmt := m.statements["GetByEmail"]
	row := stmt.QueryRow(email)

	user := User{}
	if err := row.Scan(
		&user.UserID, &user.Name, &user.Email, &user.Password, &user.Image, &user.Description, &user.Creation,
	); err != nil {
		return User{}, err
	}

	return user, nil
}

func (m UserModel) SetExtras(user *User) error {
	stmt := m.statements["SetExtras"]
	row := stmt.QueryRow(user.UserID)

	extras := UserExtras{}
	if err := row.Scan(
		&extras.TotalPosts,
	); err != nil {
		return err
	}
	user.Extras = &extras
	return nil
}

// This function is for changing the profile picture of the user
func (m UserModel) SetImage(image string, userID int) error {
	stmt := m.statements["SetImage"]
	_, err := stmt.Exec(image, userID)
	if err != nil {
		return err
	}
	return nil
}

// This function shall be called from the settings page, for when the user wants their password to be changed
func (m UserModel) SetPassword(password string, userID int) error {
	stmt := m.statements["SetPassword"]
	_, err := stmt.Exec(password, userID)
	if err != nil {
		return err
	}
	return nil
}

// This is to change the Description the user has on their profile
func (m UserModel) SetDescription(description string, userID int) error {
	stmt := m.statements["SetDescription"]
	_, err := stmt.Exec(description, userID)
	if err != nil {
		return err
	}
	return nil
}
