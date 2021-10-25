package main

import (
	"fmt"
	fdb "forum/forumDB"
	"forum/forumEnv"
	"forum/pages"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize sql.DB struct
	db := fdb.OpenDB("db/forum.db")

	// Set up handlers
	// Serve static stuff like stylesheets and js on /static/
	staticServer := http.FileServer(http.Dir("./server/static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticServer))

	// Get templates for page handlers
	templates := forumEnv.CreateTemplates("./server/templates")
	// And create an Env for page DB and Template access
	env := forumEnv.NewEnv(db, templates)

	// Then convert the Env into page-specific versions, so they act as handlers
	http.Handle("/forum", pages.Forum(env))
	http.Handle("/register", pages.Register(env))
	http.Handle("/login", pages.Login(env))

	http.HandleFunc("/", handleOther)

	// Start the server
	port := 8080
	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", port), nil))
}

func handleOther(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/forum", http.StatusFound)
		return
	}

	http.NotFound(w, r)
}
