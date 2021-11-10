package forumDB

import (
	"database/sql"
	"forum/internal/utils"
)

type Thread struct {
	ThreadID   int
	Title      string
	CategoryID int
}

type ThreadInterface interface {
	Insert(newThread Thread) (int, error)
	Get(threadID int) (Thread, error)
	ByCategory(categoryID int) ([]Thread, error)
}

type ThreadModel struct {
	DB *sql.DB
}

func (m ThreadModel) Insert(newThread Thread) (int, error) {
	stmt, err := m.DB.Prepare(
		"INSERT INTO threads(title, categoryID) values(?,?)",
	)
	utils.FatalErr(err)

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
	stmt, err := m.DB.Prepare(
		"SELECT * FROM threads WHERE threadID=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(threadID)
	thread := Thread{}
	err = row.Scan(
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
	stmt, err := m.DB.Prepare(
		"SELECT * FROM threads WHERE categoryID=?",
	)
	utils.FatalErr(err)

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
