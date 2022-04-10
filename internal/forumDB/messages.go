package forumDB

import (
	"database/sql"
	"fmt"
	"time"
)

type Message struct {
	MessageID int64
	FromID    int64
	ToID      int64
	Body      string
	Date      time.Time
}

type MessageModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewMessageModel(db *sql.DB) MessageModel {
	model := MessageModel{db: db}

	model.statements = makeStatementMap(db, "server/db/sql/models/messages.sql")

	return model
}

func (m MessageModel) NewMessage(fromID, toID int, body string) (int, error) {
	stmt := m.statements["New"]

	res, err := stmt.Exec(
		fromID,
		toID,
		body,
		time.Now(),
	)

	id, _ := res.LastInsertId()

	return int(id), fmt.Errorf("message send: %w", err)
}

func (m MessageModel) GetByFrom(userID int) (messages []Message, _ error) {
	stmt := m.statements["GetByFrom"]

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var message Message

		err = rows.Scan(
			&message.MessageID,
			&message.FromID,
			&message.ToID,
			&message.Body,
			&message.Date,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (m MessageModel) GetByTo(userID int) (messages []Message, _ error) {
	stmt := m.statements["GetByTo"]

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var message Message

		err = rows.Scan(
			&message.MessageID,
			&message.FromID,
			&message.ToID,
			&message.Body,
			&message.Date,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (m MessageModel) GetChat(fromID, toID, start int) (chat []Message, _ error) { // gets all messages between fromID and toID ordered by MessageID starting from last *start*
	stmt := m.statements["GetChat"]

	rows, err := stmt.Query(fromID, toID, fromID, toID, start)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var message Message

		err = rows.Scan(
			&message.MessageID,
			&message.FromID,
			&message.ToID,
			&message.Body,
			&message.Date,
		)
		if err != nil {
			return nil, err
		}

		chat = append(chat, message)
	}
	return chat, nil
}
