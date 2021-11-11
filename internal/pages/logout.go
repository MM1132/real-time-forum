package pages

import (
	"forum/internal/forumEnv"
	"net/http"
)

type Logout struct {
	forumEnv.Env
}

func (env Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		env.ClearSession(w)
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

func (env Logout) ClearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
