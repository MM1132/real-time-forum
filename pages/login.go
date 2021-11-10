package pages

import (
	"fmt"
	"forum/forumDB"
	"forum/forumEnv"
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
}

func (env Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Referer(), "localhost") && !strings.Contains(r.Referer(), "127.0.0.1") {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		fmt.Println(r.Referer())
		return
	}
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
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Incorrect username or password.")
		return nil
	} else if r.FormValue("pass") == user.Password {
		token, err := env.Sessions.New(&user)
		if err != nil {
			log.Panic()
		}
		cookie := &http.Cookie{
			Name:   "session",
			Value:  token,
			Path:   "/", // Otherwise it defaults to /login
			Secure: true,
			MaxAge: 86400, // One day
		}

		http.SetCookie(w, cookie)

		log.Printf("%v has logged in.", user.Name)
		log.Println()
		http.Redirect(w, r, "/forum", http.StatusFound)

		return &user
	}
	log.Println("Incorrect username or password.")
	w.WriteHeader(http.StatusUnauthorized)
	return nil
}
