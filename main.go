package main

import (
	fdb "forum/internal/forumDB"
	"forum/internal/forumEnv"
	"forum/internal/pages"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize sql.DB struct
	db := fdb.OpenDB("server/db/forum.db")
	defer db.Close()

	// Create a custom mux
	mux := http.NewServeMux()

	// Get templates for page handlers
	templates := forumEnv.CreateTemplates("./server/templates")
	// And create an Env for page DB and Template access
	env := forumEnv.NewEnv(db, templates)
	env.SiteName = "Ubian Debuntu"

	// Then convert the Env into page-specific versions, so they act as handlers
	mux.Handle("/board", pages.Board{env})
	mux.Handle("/thread", pages.Thread{env})
	mux.Handle("/register", pages.Register{env})
	mux.Handle("/login", pages.Login{env})
	mux.Handle("/logout", pages.Logout{env})
	mux.Handle("/user", pages.User{env})
	mux.Handle("/settings", pages.Settings{env})
	mux.Handle("/like", pages.Like{env})
	mux.Handle("/dislike", pages.Like{env})
	mux.Handle("/search", pages.Search{env})

	staticFS := http.FileServer(http.Dir("./server/static"))
	mux.Handle("/", forumEnv.RedirectEmpty("/board", staticFS))

	// Start the server
	host := ":" + getPort()
	log.Printf("Listening on %v\n", host)
	server := HTTPServerMux(mux, host)
	log.Fatal(server.ListenAndServe())
}

func getPort() string {
	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	return port
}

func HTTPServerMux(mux *http.ServeMux, addr string) *http.Server {
	return &http.Server{
		Addr:         addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      forumEnv.Log(mux),
	}
}
