package forumDB

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	SessionID string
	UserID    int
	Created   time.Time
}

type SessionInterface interface {
	New(user *User) (string, error)
	GetByUID(UID int) (Session, error)
	GetByToken(SID string) (Session, error)
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
	session.SessionID = uuid.NewV4().String()
	session.UserID = user.UserID
	session.Created = time.Now()
	fmt.Print(session)
	_, err = stmt.Exec(session.SessionID, session.UserID, session.Created)
	if err != nil {
		return "0", err
	}

	return string(session.SessionID), err
}

func (m SessionModel) GetByToken(sid string) (Session, error) {
	stmt, err := m.DB.Prepare(
		"SELECT * FROM sessions WHERE token=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(sid)
	session := Session{}
	err = row.Scan(&session.SessionID, &session.UserID, &session.Created)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (m SessionModel) GetByUID(uid int) (Session, error) {
	stmt, err := m.DB.Prepare(
		"SELECT * FROM sessions WHERE userID=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(uid)
	session := Session{}
	err = row.Scan(&session.SessionID, &session.UserID, &session.Created)
	if err != nil {
		return session, err
	}

	return session, nil
}
