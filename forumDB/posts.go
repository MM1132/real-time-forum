package forumDB

import (
	"database/sql"
	"time"
)

type Thread struct {
	ThreadID   int
	Title      string
	CategoryID int
}

type Post struct {
	PostID   int
	Content  []byte
	UserID   int
	ThreadID int
	Date     time.Time
}

func InsertThread(db *sql.DB, newThread *Thread, newPost *Post) (int, int, error) {
	stmt, err := db.Prepare(
		"INSERT INTO threads(title, categoryID) values(?,?)",
	)
	fatalErr(err)

	res, err := stmt.Exec(
		newThread.Title,
		newThread.CategoryID,
	)
	fatalErr(err)
	threadID, _ := res.LastInsertId()
	newPost.ThreadID = int(threadID)
	postID, err := InsertPost(db, newPost)
	if err != nil {
		panic(err)
	}
	return int(threadID), postID, err
}

func InsertPost(db *sql.DB, newPost *Post) (int, error) {
	stmt, err := db.Prepare(
		"INSERT INTO posts(content, userID, threadID, date) values(?,?,?,?)",
	)
	fatalErr(err)

	res, err := stmt.Exec(
		newPost.Content,
		newPost.UserID,
		newPost.ThreadID,
		time.Now(),
	)
	fatalErr(err)

	id, _ := res.LastInsertId()
	return int(id), err
}

func GetPost(db *sql.DB, postID int) (*Post, error) {
	stmt, err := db.Prepare(
		"SELECT * FROM posts WHERE postID=?",
	)
	fatalErr(err)

	row := stmt.QueryRow(postID)
	post := &Post{}
	err = row.Scan(
		&post.PostID,
		&post.ThreadID,
		&post.UserID,
		&post.Content,
		&post.Date,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func GetThread(db *sql.DB, threadID int) (*Thread, error) {
	stmt, err := db.Prepare(
		"SELECT * FROM threads WHERE threadID=?",
	)
	fatalErr(err)

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
