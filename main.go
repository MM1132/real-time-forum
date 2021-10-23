package main

import (
	"fmt"
	fdb "forum/forumDB"
	"forum/pages"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize sql.DB struct
	db := fdb.InitializeDB()
	defer db.Close()

	// Set up handlers
	// Serve static stuff like stylesheets and js on /static/
	staticServer := http.FileServer(http.Dir("./server/static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticServer))

	// Create templates for page handlers
	templates := pages.CreateTemplates("./server/templates")
	http.Handle("/index", pages.IndexHandler(db, templates))
	http.Handle("/register", pages.RegisterHandler(db, templates))
	http.Handle("/login", pages.LoginHandler(db, templates))

	http.HandleFunc("/", handleOther)

	// Start the server
	port := 8080
	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}

func handleOther(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/forum", http.StatusFound)
		return
	}

	http.NotFound(w, r)
}
