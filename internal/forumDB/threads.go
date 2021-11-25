package forumDB

import (
	"database/sql"
	"time"
)

type Thread struct {
	ThreadID int
	Title    string
	BoardID  int

	Extras *ThreadExtras
}

type ThreadExtras struct {
	CountPosts int
	CountUsers int

	LatestID       int
	LatestAuthorID int
	LatestAuthor   string
	LatestDate     time.Time

	OldestID       int
	OldestAuthorID int
	OldestAuthor   string
	OldestDate     time.Time
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

// GetPageThreads returns all threads that are supposed to be in a single page
func (m ThreadModel) GetPageThreads(boardID, page, pageSize int, orderKey string) ([]Thread, error) {
	key := "GetPageThreads-" + orderKey
	stmt := m.statements[key]

	rows, err := stmt.Query(
		boardID,
		pageSize,
		(page-1)*pageSize,
	)
	if err != nil {
		return nil, err
	}

	var threads []Thread
	for rows.Next() {
		thread := Thread{Extras: &ThreadExtras{}}
		err = rows.Scan(
			&thread.ThreadID,
			&thread.Title,
			&thread.BoardID,

			&thread.Extras.CountPosts,
			&thread.Extras.CountUsers,

			&thread.Extras.LatestID,
			&thread.Extras.LatestAuthorID,
			&thread.Extras.LatestAuthor,
			&thread.Extras.LatestDate,

			&thread.Extras.OldestID,
			&thread.Extras.OldestAuthorID,
			&thread.Extras.OldestAuthor,
			&thread.Extras.OldestDate,
		)
		if err != nil {
			return nil, err
		}

		threads = append(threads, thread)
	}

	return threads, nil
}

func (m ThreadModel) ThreadCount(boardID int) int {
	stmt := m.statements["ThreadCount"]

	row := stmt.QueryRow(boardID)

	var count int
	_ = row.Scan(&count)

	return count
}
