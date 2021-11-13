package forumDB

import (
	"database/sql"
)

type Thread struct {
	ThreadID int
	Title    string
	BoardID  int
}

type ThreadModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewThreadModel(db *sql.DB) ThreadModel {
	model := ThreadModel{db: db}

	model.statements = makeStatementMap(db, "server/db/sql/models/threads.sql")

	return model
}

func (m ThreadModel) Insert(newThread Thread) (int, error) {
	stmt := m.statements["Insert"]

	res, err := stmt.Exec(
		newThread.Title,
		newThread.BoardID,
	)
	if err != nil {
		return 0, err
	}

	threadID, _ := res.LastInsertId()
	return int(threadID), err
}

// Get the thread by its ID
func (m ThreadModel) Get(threadID int) (Thread, error) {
	stmt := m.statements["Get"]
	row := stmt.QueryRow(threadID)

	thread := Thread{}
	err := row.Scan(
		&thread.ThreadID,
		&thread.Title,
		&thread.BoardID,
	)
	if err != nil {
		return Thread{}, err
	}

	return thread, nil
}

func (m ThreadModel) ByBoard(boardID int) ([]Thread, error) {
	stmt := m.statements["ByBoard"]
	rows, err := stmt.Query(boardID)
	if err != nil {
		return nil, err
	}

	var threads []Thread
	for rows.Next() {
		thread := Thread{}
		err = rows.Scan(
			&thread.ThreadID,
			&thread.Title,
			&thread.BoardID,
		)
		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}
