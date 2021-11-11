package forumEnv

import (
	"database/sql"
	"forum/internal/forumDB"
	"html/template"
)

// Env is used in all page handlers to give database and template access.
// It must be initialized with NewEnv()
type Env struct {
	SiteName string

	DB        *sql.DB
	Templates map[string]*template.Template

	Users      forumDB.UserModel
	Posts      forumDB.PostModel
	Threads    forumDB.ThreadModel
	Categories forumDB.CategoryModel
	Sessions   forumDB.SessionModel
}

// NewEnv creates a new Env for use in http handlers
func NewEnv(db *sql.DB, templates map[string]*template.Template) Env {
	return Env{
		DB:        db,
		Templates: templates,

		Users:      forumDB.NewUserModel(db),
		Posts:      forumDB.NewPostModel(db),
		Threads:    forumDB.NewThreadModel(db),
		Categories: forumDB.NewCategoryModel(db),
		Sessions:   forumDB.NewSessionModel(db),
	}
}
