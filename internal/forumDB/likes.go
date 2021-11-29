package forumDB

import (
	"database/sql"
	"time"
)

type Like struct {
	PostID int
	UserID int
	Value  int
	Date   time.Time
}

type LikeModel struct {
	db         *sql.DB
	statements map[string]*sql.Stmt
}

func NewLikeModel(db *sql.DB) LikeModel {
	model := LikeModel{db: db}

	model.statements = makeStatementMap(db, "server/db/sql/models/likes.sql")

	return model
}

func (m LikeModel) Insert(postID int, userID int, value int) error {
	stmt := m.statements["Insert"]

	_, err := stmt.Exec(
		postID,
		userID,
		value,
		time.Now(),
	)

	return err
}

func (m LikeModel) GetPostTotal(postID int) (int, error) {
	stmt := m.statements["GetPostTotal"]

	row := stmt.QueryRow(postID)

	var sum int
	err := row.Scan(&sum)

	return sum, err
}

func (m LikeModel) GetAllUser(userID int) ([]Like, error) {
	stmt := m.statements["GetAllUser"]

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	var likes []Like
	for rows.Next() {
		like := Like{}
		err = rows.Scan(
			&like.PostID,
			&like.UserID,
			&like.Value,
			&like.Date,
		)
		if err != nil {
			return nil, err
		}
		likes = append(likes, like)
	}
	return likes, nil
}

func (m LikeModel) Delete(postID int, userID int) error {
	stmt := m.statements["Delete"]
	_, err := stmt.Exec(postID, userID)
	return err
}

func (m LikeModel) Get(postID int, userID int) (Like, error) {
	stmt := m.statements["Get"]
	row := stmt.QueryRow(postID, userID)

	var like Like
	err := row.Scan(
		&like.PostID,
		&like.UserID,
		&like.Value,
		&like.Date,
	)

	return like, err
}

func (m LikeModel) Update(postID int, userID int, value int) error {
	stmt := m.statements["Update"]
	_, err := stmt.Exec(value, time.Now(), userID, postID)
	return err
}

type LikedPost struct {
	Value  int
	Post   Post
	Author User
}

func (m LikeModel) GetLikedPosts(userID int) ([]LikedPost, error) {
	stmt := m.statements["GetLikedPosts"]
	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	var posts []LikedPost
	for rows.Next() {
		likedPost := LikedPost{}
		err = rows.Scan(
			&likedPost.Value,
			&likedPost.Post.PostID,
			&likedPost.Post.ThreadID,
			&likedPost.Post.UserID,
			&likedPost.Post.Content,
			&likedPost.Post.Date,
			&likedPost.Author.UserID,
			&likedPost.Author.Name,
			&likedPost.Author.Email,
			&likedPost.Author.Password,
			&likedPost.Author.Image,
			&likedPost.Author.Description,
			&likedPost.Author.Creation,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, likedPost)
	}

	return posts, nil
}
