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
type Register struct {
	db   *sql.DB
	tmpl *template.Template
}

// Contains things that are generated per request
type registerData struct {
	Title     string // Title should be on every page
	UserCount int
}

// Initializes a page handler with db and template values, and returns a ready to use http.Handler
func RegisterHandler(db *sql.DB, templates map[string]*template.Template) http.Handler {
	name := "register"

	// Get the right template
	tmpl, ok := templates[name]
	if !ok {
		log.Fatalf("Could not find the template for %v.html\n", name)
	}

	// Define a new Index handler with the db and template set
	handler := &Register{db: db, tmpl: tmpl}
	return handler
}

func (p Register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We must create a new indexData struct because it can't be shared between requests
	data := &registerData{Title: "Forum Register"}
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
		register(p.db, w, r)

	}
}

func (d *registerData) getUserRowCount(db *sql.DB) error {
	row := db.QueryRow("SELECT COUNT(*) FROM users")

	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}
	d.UserCount = count
	return nil
}

func register(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	newUser := forumDB.User{Name: r.FormValue("username"), Email: r.FormValue("email"), Password: r.FormValue("password")}

	_, err := forumDB.InsertUser(db, &newUser)
	if err != nil {
		log.Println(err)
	}
	fmt.Print("New user registered: ")
	fmt.Println(newUser)
	// forumDB.InsertUser(p.db, &newUser)
}
