package forumDB

import (
	"database/sql"
	"forum/utils"
	"time"
)

type Post struct {
	PostID   int
	Content  []byte
	UserID   int
	ThreadID int
	Date     time.Time
}

func InsertPost(db *sql.DB, newPost *Post) (int, error) {
	stmt, err := db.Prepare(
		"INSERT INTO posts(content, userID, threadID, date) values(?,?,?,?)",
	)
	utils.FatalErr(err)

	res, err := stmt.Exec(
		newPost.Content,
		newPost.UserID,
		newPost.ThreadID,
		time.Now(),
	)
	utils.FatalErr(err)

	id, _ := res.LastInsertId()
	return int(id), err
}

func GetPost(db *sql.DB, postID int) (*Post, error) {
	stmt, err := db.Prepare(
		"SELECT * FROM posts WHERE postID=?",
	)
	utils.FatalErr(err)

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
