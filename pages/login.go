package pages

import (
	"fmt"
	"forum/forumDB"
	"forum/forumEnv"
	"log"
	"net/http"
)

type Login struct {
	forumEnv.Env
}

// Contains things that are generated for every request and passed on to the template
type loginData struct {
	forumEnv.GenericData
}

func (env Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We must create a new loginData struct because it can't be shared between requests
	data := &loginData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}
	data.AddTitle("Login")

	switch r.Method {
	case "GET":
	//
	case "POST":
		env.login(w, r)

	}
	// Finally execute the template with the data we got
	tmpl := env.Templates["login"]
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	log.Printf("Served %v to %v\n", data.Title, r.RemoteAddr)
}

func (env Login) login(w http.ResponseWriter, r *http.Request) *forumDB.User {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return nil
	}
	user, err := env.Users.ByName(r.FormValue("name"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
		return nil
	}
	if r.FormValue("pass") == user.Password {
		cookie, err := r.Cookie("session")
		if err != nil {
			sessionID, err := env.Sessions.New(&user)
			if err != nil {
				log.Panic()
			} else {
				cookie = &http.Cookie{
					Name:   "session",
					Value:  sessionID,
					Path:   "/", // Otherwise it defaults to /login
					Secure: true,
					MaxAge: 86400, // One day
				}
			}
		}
		http.SetCookie(w, cookie)

		log.Printf("%v has logged in.", user.Name)
		log.Println()

		return &user
	}
	log.Println("Incorrect username or password.")
	w.WriteHeader(http.StatusUnauthorized)
	return nil
}
