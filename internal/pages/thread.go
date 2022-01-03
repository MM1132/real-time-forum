package pages

import (
	"fmt"
	"forum/internal/forumDB"
	"forum/internal/forumEnv"
	"forum/internal/search"
	"log"
	"net/http"
	"strings"
)

type Thread struct {
	forumEnv.Env
}

type ThreadData struct {
	forumEnv.GenericData
	Thread      forumDB.Thread
	Breadcrumbs []forumDB.Board

	PostsSearch *search.PostSearch

	HighlightPost int
}

func (env Thread) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := ThreadData{}
	data.InitData(env.Env, r)

	data.HighlightPost, _ = GetQueryInt("post", r)

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

	// If there's a POST request for this thread, let another function handle it.
	if r.Method == "POST" {
		err := env.post(w, r, data)
		if err != nil {
			log.Printf("error inserting post: %v", err)
			return
		}
		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
		return
	} else if r.Method != "GET" { // If it's neither GET or POST, don't allow it
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get a list of threads in this board
	postSearch := search.NewPostSearch(env.Env, data.CurrentURL, data.User.UserID)
	postSearch.Name = "posts-page"
	postSearch.ProcessRequestBasic(r)
	postSearch.ThreadID.Int64 = int64(threadIdInt)
	postSearch.ThreadID.Valid = true
	postSearch.ProcessOrder("date-asc")

	err = postSearch.DoSearch(env.Env)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
	data.PostsSearch = postSearch

	// BreadCrumbs
	// Get the board the thread is in
	if data.Breadcrumbs, err = env.Boards.GetBreadcrumbs(thread.BoardID); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// And finally we are executing the template with the data we got
	template := env.Templates["thread"]
	if err := template.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}

func (env Thread) post(w http.ResponseWriter, r *http.Request, data ThreadData) error {
	// Don't allow unauthorized users to post

	threadID := data.Thread.ThreadID

	err := checkUser(data.GenericData, r.RemoteAddr)
	if err != nil {
		sendErr(err, w, http.StatusForbidden)
		return err
	}

	// See the postThread() comment
	if strings.TrimSpace(r.FormValue("post")) == "" {
		err = fmt.Errorf("empty post content from %v", r.RemoteAddr)
		sendErr(err, w, http.StatusBadRequest)
		return err
	}

	id, err := writePost(r.FormValue("post"), data.User.UserID, threadID, env)

	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return err
	}

	// Notify OP author of a new reply
	op, err := env.Threads.GetOP(threadID)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return err
	}

	// Don't notify when replying to own thread
	if op.User.UserID != data.User.UserID {
		_, err = env.Pings.Send(
			op.User.UserID,
			fmt.Sprintf(`%v has replied to your thread "%v"`, data.User.Name, data.Thread.Title),
			fmt.Sprintf("/thread?id=%[1]v&post=%[2]v#%[2]v", threadID, id),
		)
		if err != nil {
			log.Printf(`failed to make reply notification to %v from %v: %v`, op.User.UserID, data.User.UserID, err)
		}
	}

	return nil
}
