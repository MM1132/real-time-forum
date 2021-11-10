package pages

import (
	"fmt"
	"forum/internal/forumDB"
	"forum/internal/forumEnv"
	"log"
	"net/http"
)

type Register struct {
	forumEnv.Env
}

// Contains things that are generated for every request and passed on to the template
type registerData struct {
	forumEnv.GenericData
}

func (env Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We must create a new indexData struct because it can't be shared between requests
	data := &registerData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}
	switch r.Method {
	case "GET":
	//
	case "POST":
		env.register(w, r)

	}
	data.AddTitle("Register")

	// Finally execute the template with the data we got
	tmpl := env.Templates["register"]
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}

func (env Register) register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	newUser := forumDB.User{Name: r.FormValue("username"), Email: r.FormValue("email"), Password: r.FormValue("password")}

	_, err := env.Users.Insert(newUser)
	if err != nil {
		log.Println(err)
	}
	log.Print("New user registered: ")
	log.Println(newUser)
	http.Redirect(w, r, "/login", http.StatusFound)

	// forumDB.InsertUser(p.db, &newUser)
}
