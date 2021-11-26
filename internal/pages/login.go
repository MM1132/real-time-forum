package pages

import (
	"forum/internal/forumDB"
	"forum/internal/forumEnv"
	"log"
	"net/http"
	"strings"
)

type Login struct {
	forumEnv.Env
}

// Contains things that are generated for every request and passed on to the template
type loginData struct {
	forumEnv.GenericData
	Errors map[string]string // string didn't work for some reason. Maybe we'll add other errors in the future so whatever.
}

func (env Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := loginData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}
	if data.User.UserID != 0 { // access denied if logged in
		http.Redirect(w, r, "/board", http.StatusTemporaryRedirect)
		return
	}
	// We must create a new loginData struct because it can't be shared between requests

	data.AddTitle("Login")
	data.Errors = make(map[string]string)

	if r.Method == "POST" {
		if env.validate(r, data) {
			env.login(w, r)
		}
	}
	// Finally execute the template with the data we got
	tmpl := env.Templates["login"]
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	log.Printf("Served %v to %v\n", data.Title, r.RemoteAddr)
}

func (env Login) login(w http.ResponseWriter, r *http.Request) *forumDB.User { // creates a new session for specified user, only usable in POST request. Returns pointer to user if successful, if not returns nil.
	// username is always capitalized, lowercase
	user, _ := env.Users.GetByName(strings.Title(strings.ToLower(r.FormValue("username"))))

	token, err := env.Sessions.Insert(user.UserID)
	if err != nil {
		log.Panic()
	}

	cookie := &http.Cookie{ // creates new cookie
		Name:   "session",
		Value:  token,
		Path:   "/",   // Otherwise it defaults to /login
		Secure: true,  // true will not work on connections not localhost or HTTPS secured
		MaxAge: 86400, // One day
	}

	w.Header().Add("Set-Cookie", cookie.String())

	log.Printf("%v has logged in.\n", user.Name)
	http.Redirect(w, r, "/board", http.StatusFound)

	return &user
}

func (env Login) validate(r *http.Request, data loginData) bool {
	user, err := env.Users.GetByName(strings.Title(strings.ToLower(r.FormValue("username"))))
	if err != nil {
		data.Errors["Error"] = "Incorrect username or password."
	}
	if r.FormValue("password") != user.Password {
		data.Errors["Error"] = "Incorrect username or password."
	}
	return len(data.Errors) == 0
}
