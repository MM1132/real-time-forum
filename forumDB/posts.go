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

type PostInterface interface {
	Insert(newPost Post) (int, error)
	Get(postID int) (Post, error)
}

type PostModel struct {
	DB *sql.DB
}

func (m PostModel) Insert(newPost Post) (int, error) {
	stmt, err := m.DB.Prepare(
		"INSERT INTO posts(content, userID, threadID, date) values(?,?,?,?)",
	)
	utils.FatalErr(err)

	res, err := stmt.Exec(
		newPost.Content,
		newPost.UserID,
		newPost.ThreadID,
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return int(id), err
}

func (m PostModel) Get(postID int) (Post, error) {
	stmt, err := m.DB.Prepare(
		"SELECT * FROM posts WHERE postID=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(postID)
	post := Post{}
	err = row.Scan(
		&post.PostID,
		&post.ThreadID,
		&post.UserID,
		&post.Content,
		&post.Date,
	)
	if err != nil {
		return Post{}, err
	}

	return post, nil
}
