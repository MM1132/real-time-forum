package pages

import (
	"database/sql"
	"fmt"
	"forum/internal/forumDB"
	"forum/internal/forumEnv"
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
	Posts       []forumDB.Post
	Breadcrumbs []forumDB.Board

	HighlightPost int
}

func (env Thread) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := ThreadData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}

	data.HighlightPost, _ = GetQueryInt("post", r)

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

	for i, post := range data.Posts {
		sum, err := env.Likes.GetPostTotal(post.PostID)

		if err == sql.ErrNoRows {
			data.Posts[i].Likes = 0
			continue
		} else if err != nil {
			sendErr(err, w, http.StatusInternalServerError)
			return
		}
		data.Posts[i].Likes = sum
	}

	// Set the extras for the user
	for i := range data.Posts {
		if err := env.Users.SetExtras(&data.Posts[i].User); err != nil {
			sendErr(err, w, http.StatusInternalServerError)
		}
	}

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
	// Don't allow unathorized users to post

	threadID, err := GetQueryInt("id", r)
	// Whatever the error, this is most definitely the fault of the client
	if err != nil {
		sendErr(err, w, http.StatusBadRequest)
		return err
	}

	err = checkUser(data.GenericData, r.RemoteAddr)
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

	err = writePost(r.FormValue("post"), data.User.UserID, threadID, env)

	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return err
	}

	return nil
}
