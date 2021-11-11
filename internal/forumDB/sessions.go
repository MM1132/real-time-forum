package forumDB

import (
	"database/sql"
	"forum/internal/utils"
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
	statements := make(map[string]*sql.Stmt)
	model := SessionModel{db: db}

	var err error
	statements["New"], err = db.Prepare("INSERT INTO sessions(token, userID, created) values(?,?,?)")
	utils.FatalErr(err)

	statements["GetByToken"], err = db.Prepare("SELECT * FROM sessions WHERE token=?")
	utils.FatalErr(err)

	statements["GetByUserID"], err = db.Prepare("SELECT * FROM sessions WHERE userID=?")
	utils.FatalErr(err)

	model.statements = statements
	return model
}

// creates a new session for specified user
func (m SessionModel) New(user *User) (string, error) {
	session := Session{}
	session.Token = uuid.NewV4().String()
	session.UserID = user.UserID
	session.Created = time.Now()

	stmt := m.statements["New"]
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
