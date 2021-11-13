package forumDB

import (
	"database/sql"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Session struct {
	Token   string
	UserID  int
	Created time.Time
}

type SessionModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewSessionModel(db *sql.DB) SessionModel {
	model := SessionModel{db: db}

	model.statements = makeStatementMap(db, "server/db/sql/models/sessions.sql")

	return model
}

// creates a new session for specified user
func (m SessionModel) Insert(userID int) (string, error) {
	session := Session{}
	session.Token = uuid.NewV4().String()
	session.UserID = userID
	session.Created = time.Now()

	stmt := m.statements["Insert"]
	if _, err := stmt.Exec(session.Token, session.UserID, session.Created); err != nil {
		return "", err
	}

	return session.Token, nil
}

func (m SessionModel) GetByToken(token string) (Session, error) {
	row := m.statements["GetByToken"].QueryRow(token)
	session := Session{}
	err := row.Scan(&session.Token, &session.UserID, &session.Created)
	if err != nil {
		return Session{}, err
	}

	return session, nil
}

func (m SessionModel) GetByUserID(userID int) (Session, error) {
	row := m.statements["GetByUserID"].QueryRow(userID)
	session := Session{}
	err := row.Scan(&session.Token, &session.UserID, &session.Created)
	if err != nil {
		return session, err
	}

	return session, nil
}
