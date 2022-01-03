package forumDB

import (
	"database/sql"
	"fmt"
	"time"
)

type Ping struct {
	PingID  int
	UserID  int
	Content string
	Link    string
	Date    time.Time
}

type PingModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewPingModel(db *sql.DB) PingModel {
	model := PingModel{db: db}

	model.statements = makeStatementMap(db, "server/db/sql/models/pings.sql")

	return model
}

func (m PingModel) Send(userID int, content string, link string) (int, error) {
	stmt := m.statements["Send"]

	res, err := stmt.Exec(
		userID,
		content,
		link,
		time.Now(),
	)

	id, _ := res.LastInsertId()

	return int(id), fmt.Errorf("ping send: %w", err)
}

func (m PingModel) Delete(pingID int) error {
	stmt := m.statements["Delete"]

	_, err := stmt.Exec(
		pingID,
	)

	return fmt.Errorf("ping delete: %w", err)
}

func (m PingModel) GetByUser(userID int) (pings []Ping, _ error) {
	stmt := m.statements["GetByUser"]

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var ping Ping

		err = rows.Scan(
			&ping.PingID,
			&ping.UserID,
			&ping.Content,
			&ping.Link,
			&ping.Date,
		)
		if err != nil {
			return nil, err
		}

		pings = append(pings, ping)
	}

	return pings, nil
}
