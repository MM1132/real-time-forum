package pages

import (
	fdb "forum/forumDB"
	"forum/forumEnv"
	"net/http"
	"strconv"
)

type Thread struct {
	forumEnv.Env
}

type ThreadData struct {
	forumEnv.GenericData
	Thread      fdb.Thread
	Posts       []fdb.Post
	Breadcrumbs []fdb.Category
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
	threadIdString := r.URL.Query().Get("id")
	// Then we try turning the id into an integer
	if threadIdString == "" {
		http.NotFound(w, r)
	}
	threadIdInt := 0
	if id, err := strconv.Atoi(threadIdString); err != nil {
		http.NotFound(w, r)
	} else {
		threadIdInt = id
	}

	// Get all the information about the thread by its ID
	thread, err := env.Threads.Get(threadIdInt)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	}

	// Create the empty struct for storing the data for the thread page
	data.Thread = thread
	data.AddTitle(thread.Title)

	// Get all the posts that match this thread's ID
	if data.Posts, err = env.Posts.GetByThreadID(threadIdInt); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	}

	// BreadCrumbs
	// Get the category the thread is in
	tempCategory, err := env.Categories.Get(thread.CategoryID)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	}
	data.Breadcrumbs = []fdb.Category{tempCategory}
	// Now loop through all the categories and add them all to the amazing breadcrumbs
	for tempCategory.ParentID.Valid {
		tempCategory, err = env.Categories.Get(int(tempCategory.ParentID.Int64))
		if err != nil {
			sendErr(err, w, http.StatusInternalServerError)
			continue
		}
		// And of course, append the category to the end of breadcrumbs
		data.Breadcrumbs = append([]fdb.Category{tempCategory}, data.Breadcrumbs...)
	}

	// And finally we are executing the template with the data we got
	template := env.Templates["thread"]
	if err := template.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}
