package pages

import (
	"forum/internal/forumDB"
	fdb "forum/internal/forumDB"
	"forum/internal/forumEnv"
	"log"
	"net/http"
)

type Board struct {
	forumEnv.Env
}

// Contains things that are generated for every request and passed on to the template
type boardData struct {
	forumEnv.GenericData
	ChildBoards []fdb.Board
	Breadcrumbs []fdb.Board

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
		http.NotFound(w, r)
	}

	// Get our current board
	thisBoard, err := env.Boards.Get(thisBoardID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data.AddTitle(thisBoard.Name)

	// And then its children
	childBoards, err := env.Boards.GetChildren(thisBoard.BoardID)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// Then set their extra data
	if err := env.Boards.SetSliceExtras(childBoards); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

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
