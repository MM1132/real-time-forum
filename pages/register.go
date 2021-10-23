package pages

import (
	"fmt"
	"forum/forumDB"
	"forum/forumEnv"
	"log"
	"net/http"
)

type Register forumEnv.Env

// Contains things that are generated per request
type registerData struct {
	Title     string // Title should be on every page
	UserCount int
}

func (e Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := e.Templates["Register"]

	// We must create a new indexData struct because it can't be shared between requests
	data := &registerData{Title: "Forum Register"}

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
		e.register(w, r)

	}
}

func (e Register) register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	newUser := forumDB.User{Name: r.FormValue("username"), Email: r.FormValue("email"), Password: r.FormValue("password")}

	_, err := e.Users.Insert(newUser)
	if err != nil {
		log.Println(err)
	}
	fmt.Print("New user registered: ")
	fmt.Println(newUser)
	// forumDB.InsertUser(p.db, &newUser)
}
