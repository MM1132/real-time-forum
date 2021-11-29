package forumEnv

import (
	"fmt"
	"forum/internal/forumDB"
	"html/template"
	"net/http"
	"net/url"
)

type GenericData struct {
	SiteName string
	Title    string
	Session  forumDB.Session
	User     forumDB.User

	CurrentURL  url.URL
	SearchValue string

	Theme string
}

func (data *GenericData) InitData(env Env, r *http.Request) {
	data.SiteName = env.SiteName
	data.Title = env.SiteName
	data.CurrentURL = *r.URL
	data.SearchValue = r.FormValue("search")

	data.parseSessionCookie(env, r)

	return
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

// Query returns the current URI, except with new key/value pairs added. Setting an already existing key will replace its value.
// This function is meant to be used in templates for hrefs.
func (data GenericData) Query(kvp ...string) (template.URL, error) {
	if len(kvp)%2 == 1 {
		return "", fmt.Errorf(`need an even number of args`)
	}

	currentURL := data.CurrentURL

	query := currentURL.Query()
	for i := 0; i < len(kvp); i += 2 {
		query.Set(kvp[i], kvp[i+1])
	}

	currentURL.RawQuery = query.Encode()
	return template.URL(currentURL.RequestURI()), nil
}
