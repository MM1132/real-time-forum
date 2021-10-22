package pages

import (
	"database/sql"
	"forum/utils"
	"html/template"
	"log"
	"net/http"
)

// A handler for index.html, containing things that are shared between requests.
// It should be initialized with IndexHandler()
type Index struct {
	db   *sql.DB
	tmpl *template.Template
}

// Contains things that are generated per request
type indexData struct {
	Title     string // Title should be on every page
	UserCount int
}

// Initializes a page handler with db and template values, and returns a ready to use http.Handler
func IndexHandler(db *sql.DB, templates map[string]*template.Template) http.Handler {
	name := "index"

	// Get the right template
	tmpl, ok := templates[name]
	if !ok {
		utils.FatalErr(noTemplateError(name))
	}

	// Define a new Index handler with the db and template set
	handler := &Index{db: db, tmpl: tmpl}
	return handler
}

func (p Index) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We must create a new indexData struct because it can't be shared between requests
	data := &indexData{Title: "Forum Index"}
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
}

func (d *indexData) getUserRowCount(db *sql.DB) error {
	row := db.QueryRow("SELECT COUNT(*) FROM users")

	var count int
	err := row.Scan(&count)
	if err != nil {
		return err
	}

	d.UserCount = count
	return nil
}
