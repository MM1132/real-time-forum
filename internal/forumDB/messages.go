package forumDB

import (
	"database/sql"
	"time"
)

type Message struct {
	MessageID int       `json:"id"`
	FromID    int       `json:"from"`
	ToID      int       `json:"to"`
	Body      string    `json:"body"`
	Date      time.Time `json:"date"`
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

func (m MessageModel) Insert(fromID, toID int, body string) (int, error) {
	stmt := m.statements["Insert"]

	res, err := stmt.Exec(
		fromID,
		toID,
		body,
		time.Now(),
	)
	if err != nil {
		return 0, nil
	}

	id, _ := res.LastInsertId()

	return int(id), nil
}

// Gets the message by it's ID, I suppose
func (m MessageModel) GetByID(messageID int) (Message, error) {
	stmt := m.statements["GetByID"]

	row := stmt.QueryRow(messageID)

	message := Message{}
	if err := row.Scan(
		&message.MessageID,
		&message.FromID,
		&message.ToID,
		&message.Body,
		&message.Date,
	); err != nil {
		return Message{}, err
	}

	return message, nil
}

// Gets all the messages written by the user
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

// gets 10 messages between fromID and toID ordered by MessageID starting from last - start
func (m MessageModel) GetChat(fromID, toID, start int) (chat []Message, _ error) {
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

func (m MessageModel) GetRecentIDs(userID int) ([]int, error) {
	stmt := m.statements["GetRecentIDs"]

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	var userIDs []int
	for rows.Next() {
		var id int

		err = rows.Scan(
			&id,
		)
		if err != nil {
			return nil, err
		}

		userIDs = append(userIDs, id)
	}

	return userIDs, nil
}

// When client requests the server for messages, this shall be the format of their request
type RequestHistory struct {
	UserID int `json:"id"`
	Index  int `json:"index"`
}
