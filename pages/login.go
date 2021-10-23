package pages

import (
	"fmt"
	"forum/forumEnv"
	"log"
	"net/http"
)

type Login forumEnv.Env

// Contains things that are generated for every request and passed on to the template
type loginData struct {
	Title string // Title should be on every page
}

func (e Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := e.Templates["login"]

	// We must create a new loginData struct because it can't be shared between requests
	data := &loginData{Title: "Forum Login"}

	// Finally execute the template with the data we got
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	log.Printf("Served %v to %v\n", data.Title, r.RemoteAddr)
	switch r.Method {
	case "GET":
	//
	case "POST":
		e.login(w, r)

	}
}

func (e Login) login(w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return false
	}
	user, err := e.Users.ByName(r.FormValue("name"))
	if err != nil {
		log.Print(err)
	}
	if r.FormValue("pass") == user.Password {
		log.Printf("%v has logged in.", user.Name)
		log.Println()
		return true
	}
	log.Println("Incorrect username or password.")
	return false
}
