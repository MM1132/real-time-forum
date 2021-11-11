package pages

import (
	"forum/internal/forumDB"
	"forum/internal/forumEnv"
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

	// Check for if the request is of the right type
	if r.Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the id of the thread we return to the client
	threadIdInt, err := GetThreadID(r)
	if err != nil {
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
