package forumDB

import (
	"database/sql"
	"forum/utils"
)

type Thread struct {
	ThreadID   int
	Title      string
	CategoryID int
}

func InsertThread(db *sql.DB, newThread *Thread, newPost *Post) (int, int, error) {
	stmt, err := db.Prepare(
		"INSERT INTO threads(title, categoryID) values(?,?)",
	)
	utils.FatalErr(err)

	res, err := stmt.Exec(
		newThread.Title,
		newThread.CategoryID,
	)
	utils.FatalErr(err)
	threadID, _ := res.LastInsertId()
	newPost.ThreadID = int(threadID)
	postID, err := InsertPost(db, newPost)
	if err != nil {
		panic(err)
	}
	return int(threadID), postID, err
}

func GetThread(db *sql.DB, threadID int) (*Thread, error) {
	stmt, err := db.Prepare(
		"SELECT * FROM threads WHERE threadID=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(threadID)
	thread := &Thread{}
	err = row.Scan(
		&thread.ThreadID,
		&thread.Title,
		&thread.CategoryID,
	)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func GetThreadsByCategoryID(db *sql.DB, categoryID int) ([]Thread, error) {
	stmt, err := db.Prepare(
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
