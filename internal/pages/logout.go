package pages

import (
	"forum/internal/forumEnv"
	"net/http"
)

type Logout struct {
	forumEnv.Env
}

type logoutData struct {
	forumEnv.GenericData
}

func (env Logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := logoutData{}
	data.InitData(env.Env, r)

	if data.User.UserID == 0 { // access denied unless logged in
		http.Redirect(w, r, "/board", http.StatusTemporaryRedirect)
		return
	}

	env.ClearSession(w)
	http.Redirect(w, r, "/board", http.StatusFound)
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
