package forumDB

import (
	"database/sql"
)

type Thread struct {
	ThreadID   int
	Title      string
	CategoryID int
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
		newThread.CategoryID,
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
		&thread.CategoryID,
	)
	if err != nil {
		return Thread{}, err
	}

	return thread, nil
}

func (m ThreadModel) ByCategory(categoryID int) ([]Thread, error) {
	stmt := m.statements["ByCategory"]
	rows, err := stmt.Query(categoryID)
	if err != nil {
		return nil, err
	}

	var threads []Thread
	for rows.Next() {
		thread := Thread{}
		err = rows.Scan(
			&thread.ThreadID,
			&thread.Title,
			&thread.CategoryID,
		)
		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}
