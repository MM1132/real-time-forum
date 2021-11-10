package forumEnv

import (
	"fmt"
	fdb "forum/forumDB"
	"net/http"
)

type GenericData struct {
	Title   string
	Session fdb.Session
	User    fdb.User

	Theme string
}

func (data *GenericData) InitData(env Env, r *http.Request) error {
	data.Title = env.SiteName

	data.parseSessionCookie(env, r)

	return nil
}

func (data *GenericData) AddTitle(title string) {
	data.Title = fmt.Sprintf("%v - %v", title, data.Title)
}

func (data *GenericData) parseSessionCookie(env Env, r *http.Request) {
	//	// check for cookie, if doesn't exist set username to guest
	cookie, err := r.Cookie("session")
	if err != nil {
		return
	}

	// if there's no session assosiated with current cookie, set session to nil and username to guest
	if session, err2 := env.Sessions.GetByToken(cookie.Value); err2 != nil {
		return
	} else { // else set user to user assosiated with session and username to username
		cookie.MaxAge = 86400 // TODO: Make resetting cookie age work

		user, err := env.Users.Get(session.UserID)
		if err != nil {
			return
		}

		data.User = user
		data.Session = session
		return
	}
}
