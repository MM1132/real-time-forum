package forumEnv

import (
	"database/sql"
	fdb "forum/forumDB"
	"html/template"
)

// Env is used in all page handlers to give database and template access.
// It must be initialized with NewEnv()
type Env struct {
	SiteName string

	DB        *sql.DB
	Templates map[string]*template.Template

	Users      fdb.UserInterface
	Posts      fdb.PostInterface
	Threads    fdb.ThreadInterface
	Categories fdb.CategoryInterface
	Sessions   fdb.SessionInterface
}

// NewEnv creates a new Env for use in http handlers
func NewEnv(db *sql.DB, templates map[string]*template.Template) Env {
	return Env{
		DB:        db,
		Templates: templates,

		Users:      fdb.UserModel{DB: db},
		Posts:      fdb.PostModel{DB: db},
		Threads:    fdb.ThreadModel{DB: db},
		Categories: fdb.CategoryModel{DB: db},
		Sessions:   fdb.SessionModel{DB: db},
	}
}
