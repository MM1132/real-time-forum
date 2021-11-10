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

	Users      forumDB.UserInterface
	Posts      forumDB.PostInterface
	Threads    forumDB.ThreadInterface
	Categories forumDB.CategoryInterface
	Sessions   forumDB.SessionInterface
}

// NewEnv creates a new Env for use in http handlers
func NewEnv(db *sql.DB, templates map[string]*template.Template) Env {
	return Env{
		DB:        db,
		Templates: templates,

		Users:      forumDB.UserModel{DB: db},
		Posts:      forumDB.PostModel{DB: db},
		Threads:    forumDB.ThreadModel{DB: db},
		Categories: forumDB.CategoryModel{DB: db},
		Sessions:   forumDB.SessionModel{DB: db},
	}
}
