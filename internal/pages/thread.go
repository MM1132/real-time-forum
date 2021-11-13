package pages

import (
	"fmt"
	"forum/internal/forumDB"
	"forum/internal/forumEnv"
	"log"
	"net/http"
)

type Thread struct {
	forumEnv.Env
}

type ThreadData struct {
	forumEnv.GenericData
	Thread      forumDB.Thread
	Posts       []forumDB.Post
	Breadcrumbs []forumDB.Category
}

func (env Thread) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := ThreadData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}

	// If there's a POST request for this thread, let another function handle it.
	if r.Method == "POST" {
		err := env.post(w, r, data)
		if err != nil {
			log.Printf("error inserting post: %v", err)
			return
		}
	} else if r.Method != "GET" { // If it's neither GET or POST, don't allow it
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the id of the thread we return to the client
	threadIdInt, err := GetQueryInt("id", r)
	if err != nil && threadIdInt != 0 {
		http.NotFound(w, r)
		return
	}

	// Get all the information about the thread by its ID
	thread, err := env.Threads.Get(threadIdInt)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// Set the thread of the empty date struct
	data.Thread = thread
	data.AddTitle(thread.Title)

	// Get all the posts that match this thread's ID
	if data.Posts, err = env.Posts.GetByThreadID(threadIdInt); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// BreadCrumbs
	// Get the category the thread is in
	if data.Breadcrumbs, err = env.Categories.GetBreadcrumbs(thread.CategoryID); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	}

	// And finally we are executing the template with the data we got
	template := env.Templates["thread"]
	if err := template.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}

func (env Thread) post(w http.ResponseWriter, r *http.Request, data ThreadData) error {
	// Don't allow unathorized users to post

	threadID, err := GetQueryInt("id", r)
	// Whatever the error, this is most definitely the fault of the client
	if err != nil {
		sendErr(err, w, http.StatusBadRequest)
		return err
	}

	if data.User.UserID == 0 {
		return fmt.Errorf("%v not authorized to post in threadID %v", r.RemoteAddr, threadID)
	}

	newPost := forumDB.Post{
		Content:  r.FormValue("post"),
		UserID:   data.User.UserID,
		ThreadID: threadID,
	}

	_, err = env.Posts.Insert(newPost)
	if err != nil {
		sendErr(err, w, http.StatusBadRequest)
		return err
	}

	return nil
}
