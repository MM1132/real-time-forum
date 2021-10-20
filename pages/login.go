package pages

import (
	"database/sql"
	"fmt"
	"forum/forumDB"
	"html/template"
	"log"
	"net/http"
)

// A handler for index.html, containing things that are shared between requests.
// It should be initialized with IndexHandler()
type Login struct {
	db   *sql.DB
	tmpl *template.Template
}

// Contains things that are generated per request
type loginData struct {
	Title     string // Title should be on every page
	UserCount int
}

// Initializes a page handler with db and template values, and returns a ready to use http.Handler
func LoginHandler(db *sql.DB, templates map[string]*template.Template) http.Handler {
	name := "login"

	// Get the right template
	tmpl, ok := templates[name]
	if !ok {
		log.Fatalf("Could not find the template for %v.html\n", name)
	}

	// Define a new Login handler with the db and template set
	handler := &Login{db: db, tmpl: tmpl}
	return handler
}

func (p Login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We must create a new loginData struct because it can't be shared between requests
	data := &loginData{Title: "Forum Login"}
	// Fill the new data struct with a value we'll use in the template
	if err := data.getUserRowCount(p.db); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// Finally execute the template with the data we got
	if err := p.tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	log.Printf("Served %v to %v\n", data.Title, r.RemoteAddr)
	switch r.Method {
	case "GET":
	//
	case "POST":
		login(p.db, w, r)

	}
}

func (d *loginData) getUserRowCount(db *sql.DB) error {
	row := db.QueryRow("SELECT COUNT(*) FROM users")

	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	d.UserCount = count
	return nil
}

func login(db *sql.DB, w http.ResponseWriter, r *http.Request) bool {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return false
	}
	user, err := forumDB.GetUserByName(db, r.FormValue("name"))
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
