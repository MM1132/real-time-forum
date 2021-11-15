package pages

import (
	"fmt"
	"forum/internal/forumDB"
	fdb "forum/internal/forumDB"
	"forum/internal/forumEnv"
	"log"
	"net/http"
	"strings"
)

type Board struct {
	forumEnv.Env
}

// Contains things that are generated for every request and passed on to the template
type boardData struct {
	forumEnv.GenericData
	ChildBoards []fdb.Board
	Breadcrumbs []fdb.Board
	BoardID     int

	Threads []forumDB.Thread
}

// ServeHTTP is called with every request this page receives.
func (env Board) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := boardData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}

	// Read query with key "id"
	thisBoardID, err := GetQueryInt("id", r)
	if err != nil && thisBoardID != 0 {
		if r.Method == "POST" {
			sendErr(err, w, http.StatusBadRequest)
			return
		}
		http.NotFound(w, r)
		return
	}
	data.BoardID = thisBoardID

	// Get our current board
	thisBoard, err := env.Boards.Get(thisBoardID)
	if err != nil {
		if r.Method == "POST" {
			sendErr(err, w, http.StatusBadRequest)
			return
		}
		sendErr(err, w, http.StatusNotFound)
		return
	}
	data.AddTitle(thisBoard.Name)

	// If it's a post request, let another function handle it first
	if r.Method == "POST" {
		if thisBoardID == 0 || thisBoard.IsGroup { // Don't allow it if root or group
			err = fmt.Errorf("invalid POST request in %v from %v", thisBoard.Name, r.RemoteAddr)
			sendErr(err, w, 405)
			return
		}
		err := env.postThread(w, r, data, thisBoardID)
		if err != nil {
			return // We do all the error sending/logging in the function
		}
		http.Redirect(w, r, r.RequestURI, http.StatusSeeOther)
		return
	}

	// Check if current board is a group
	if thisBoard.IsGroup {
		http.Redirect(w, r, fmt.Sprintf("%v?id=%v", r.URL.Path, thisBoard.ParentID.Int64), http.StatusTemporaryRedirect)
		return
	}

	// And then its children
	childBoards, err := env.Boards.GetChildren(thisBoard.BoardID)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// // Then set their extra data
	// if err := env.Boards.SetSliceExtras(childBoards); err != nil {
	// 	sendErr(err, w, http.StatusInternalServerError)
	// 	return
	// }

	data.ChildBoards = childBoards

	// BREAD-CRUMBS!!!!
	if data.Breadcrumbs, err = env.Boards.GetBreadcrumbs(thisBoard.BoardID); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// Get a list of threads in this board
	threads, err := env.Threads.ByBoard(thisBoard.BoardID)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
	data.Threads = threads

	// Finally, execute the template with the data we got
	tmpl := env.Templates["board"]
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Print(err)
		return
	}
}

func (env Board) postThread(w http.ResponseWriter, r *http.Request, data boardData, boardID int) error {
	// HTML checks if the form is filled with something, so this is for clients that don't honor the required
	// attribute. However HTML5, for some reason, the tag textarea doesn't have required pattern matching,
	// so this still gives errors to users who decide to put only white spaces in their posts,
	// but those are probably idiots and get what they deserve. Still, I think that kind of error
	// handling should be done in the front-end, so as not to waste any time on frivolus requests
	if strings.TrimSpace(r.FormValue("title")) == "" || strings.TrimSpace(r.FormValue("post")) == "" {
		err := fmt.Errorf("title or content empty from %v", r.RemoteAddr)
		sendErr(err, w, http.StatusBadRequest)
		return err
	}

	err := checkUser(data.GenericData, r.RemoteAddr)
	if err != nil {
		sendErr(err, w, http.StatusForbidden)
		return err
	}

	newThread := forumDB.Thread{
		Title:   r.FormValue("title"),
		BoardID: boardID,
	}

	threadID, err := env.Threads.Insert(newThread)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return err
	}

	writePost(r.FormValue("post"), data.Session.UserID, threadID, Thread(env))
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return err
	}

	return nil
}
