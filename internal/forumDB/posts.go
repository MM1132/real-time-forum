package forumDB

import (
	"database/sql"
	"forum/internal/utils"
	"time"
)

type Post struct {
	PostID   int
	Content  string
	UserID   int
	ThreadID int
	Date     time.Time
	User     User
}

type PostInterface interface {
	Insert(newPost Post) (int, error)
	Get(postID int) (Post, error)
	GetByThreadID(threadID int) ([]Post, error)
	GetByUserID(userID int) ([]Post, error)
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

// Return the post by its id
func (m PostModel) Get(postID int) (Post, error) {
	stmt, err := m.DB.Prepare(
		"SELECT * FROM posts p JOIN users u ON p.userID = u.userID WHERE postID=?",
	)
	utils.FatalErr(err)

	row := stmt.QueryRow(postID)
	post := Post{}
	if err = row.Scan(
		&post.PostID,
		&post.ThreadID,
		&post.UserID,
		&post.Content,
		&post.Date,
		&post.User.UserID,
		&post.User.Name,
		&post.User.Email,
		&post.User.Password,
		&post.User.Creation,
	); err != nil {
		return Post{}, err
	}

	return post, nil
}

// Get all the posts with the threadID
func (m PostModel) GetByThreadID(threadID int) ([]Post, error) {
	// Prepare the statement
	statement, err := m.DB.Prepare(
		"SELECT * FROM posts p JOIN users u ON p.userID = u.userID WHERE threadID=? ORDER BY date",
	)
	utils.FatalErr(err)

	// Get all the rows where the threadID matches
	rows, err := statement.Query(threadID)
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(
			&post.PostID,
			&post.ThreadID,
			&post.UserID,
			&post.Content,
			&post.Date,
			&post.User.UserID,
			&post.User.Name,
			&post.User.Email,
			&post.User.Password,
			&post.User.Creation,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (m PostModel) GetByUserID(userID int) ([]Post, error) {
	statement, err := m.DB.Prepare(
		"SELECT * FROM posts WHERE userID=? ORDER BY date DESC",
	)
	utils.FatalErr(err)

	rows, err := statement.Query(userID)
	if err != nil {
		return nil, err
	}

	posts := []Post{}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(
			&post.PostID,
			&post.ThreadID,
			&post.UserID,
			&post.Content,
			&post.Date,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
