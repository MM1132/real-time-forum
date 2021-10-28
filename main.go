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

	// Create a custom mux
	mux := http.NewServeMux()

	// Get templates for page handlers
	templates := forumEnv.CreateTemplates("./server/templates")
	// And create an Env for page DB and Template access
	env := forumEnv.NewEnv(db, templates)

	// Then convert the Env into page-specific versions, so they act as handlers
	mux.Handle("/forum", pages.Forum(env))
	mux.Handle("/thread", pages.Thread{env})
	mux.Handle("/register", pages.Register(env))
	mux.Handle("/login", pages.Login(env))

	staticFS := http.FileServer(http.Dir("./server/static"))
	mux.Handle("/", forumEnv.RedirectEmpty("/forum", staticFS))

	// Start the server
	port := 8080
	log.Printf("Listening on port %v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%v", port), forumEnv.Log(mux)))
}
