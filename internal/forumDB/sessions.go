package forumDB

import (
	"database/sql"
	"fmt"
	"forum/internal/utils"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	Token   string
	UserID  int
	Created time.Time
}

type SessionInterface interface {
	New(user *User) (string, error)
	GetByUserID(userID int) (Session, error)
	GetByToken(token string) (Session, error)
}

type SessionModel struct {
	DB *sql.DB
}

// creates a new session for specified user
func (m SessionModel) New(user *User) (string, error) {
	stmt, err := m.DB.Prepare(
		"INSERT INTO sessions(token, userID, created) values(?,?,?)",
	)
	utils.FatalErr(err)
	session := Session{}
	session.Token = uuid.NewV4().String()
	session.UserID = user.UserID
	session.Created = time.Now()
	fmt.Print(session)
	_, err = stmt.Exec(session.Token, session.UserID, session.Created)
	if err != nil {
		return "0", err
	}

	return string(session.Token), err
}

func (m SessionModel) GetByToken(token string) (Session, error) {
	stmt, err := m.DB.Prepare(
		"SELECT * FROM sessions WHERE token=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(token)
	session := Session{}
	err = row.Scan(&session.Token, &session.UserID, &session.Created)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (m SessionModel) GetByUserID(userID int) (Session, error) {
	stmt, err := m.DB.Prepare(
		"SELECT * FROM sessions WHERE userID=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(userID)
	session := Session{}
	err = row.Scan(&session.Token, &session.UserID, &session.Created)
	if err != nil {
		return session, err
	}

	return session, nil
}
