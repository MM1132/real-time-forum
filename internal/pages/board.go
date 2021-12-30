package pages

import (
	"fmt"
	fdb "forum/internal/forumDB"
	"forum/internal/forumEnv"
	"forum/internal/search"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Board struct {
	forumEnv.Env
}

// Contains things that are generated for every request and passed on to the template
type boardData struct {
	forumEnv.GenericData
	ThisBoard   fdb.Board
	ChildBoards []fdb.Board
	Breadcrumbs []fdb.Board
	BoardID     int

	ThreadSearch search.Searcher

	TagSuggestions []string
}

// ServeHTTP is called with every request this page receives.
func (env Board) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := boardData{}
	data.InitData(env.Env, r)

	// Read query with key "id"
	thisBoardID, err := GetQueryInt("id", r)
	if err != nil && thisBoardID != 0 {
		http.NotFound(w, r)
		return
	}
	data.BoardID = thisBoardID

	// Get our current board
	thisBoard, err := env.Boards.Get(thisBoardID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data.AddTitle(thisBoard.Name)
	data.ThisBoard = thisBoard

	// If it's a post request, let another function handle it first
	if r.Method == "POST" {
		if thisBoardID == 0 || thisBoard.IsGroup { // Don't allow it if root or group
			err = fmt.Errorf("invalid POST request in %v from %v", thisBoard.Name, r.RemoteAddr)
			sendErr(err, w, 405)
			return
		}
		newID, err := env.postThread(w, r, data, thisBoardID)
		if err != nil {
			return // We do all the error sending/logging in the function
		}
		http.Redirect(w, r, fmt.Sprint("/thread?id=", newID), http.StatusSeeOther)
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
	threadSearch := search.NewThreadSearch(env.Env, data.CurrentURL, data.User.UserID)
	threadSearch.Name = "threads-page"
	threadSearch.ProcessRequestBasic(r)
	threadSearch.BoardID.Int64 = int64(thisBoardID)
	threadSearch.BoardID.Valid = true

	err = threadSearch.DoSearch(env.Env)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
	data.ThreadSearch = threadSearch

	// Get suggestions for likes
	data.TagSuggestions = env.Threads.Tags.GetPopular(thisBoardID)

	// Finally, execute the template with the data we got
	tmpl := env.Templates["board"]
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Print(err)
		return
	}
}

func (env Board) postThread(w http.ResponseWriter, r *http.Request, data boardData, boardID int) (int, error) {
	// HTML checks if the form is filled with something, so this is for clients that don't honor the required
	// attribute. However HTML5, for some reason, the tag textarea doesn't have required pattern matching,
	// so this still gives errors to users who decide to put only white spaces in their posts,
	// but those are probably idiots and get what they deserve. Still, I think that kind of error
	// handling should be done in the front-end, so as not to waste any time on frivolus requests
	if strings.TrimSpace(r.FormValue("title")) == "" || strings.TrimSpace(r.FormValue("content")) == "" {
		err := fmt.Errorf("title or content empty from %v", r.RemoteAddr)
		sendErr(err, w, http.StatusBadRequest)
		return 0, err
	}

	err := checkUser(data.GenericData, r.RemoteAddr)
	if err != nil {
		sendErr(err, w, http.StatusForbidden)
		return 0, err
	}

	newThread := fdb.Thread{
		Title:   r.FormValue("title"),
		BoardID: boardID,
	}

	// Process the given tags into a slice (also error check)
	newThread.Tags = env.ProcessTags(r.FormValue("tags"))

	threadID, err := env.Threads.Insert(newThread)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return 0, err
	}

	_, err = writePost(r.FormValue("content"), data.Session.UserID, threadID, Thread(env))
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return 0, err
	}

	return threadID, nil
}

func (_ Board) ProcessTags(rawTags string) (tags []string) {
	if rawTags == "" || len(rawTags) > 100 {
		return nil
	}

	splits := strings.Split(rawTags, "#")
	splits = splits[1:]

	for _, tag := range splits {
		tag = strings.TrimSpace(tag)

		r1 := regexp.MustCompile(`\s`)
		tag = r1.ReplaceAllString(tag, "-")

		r2 := regexp.MustCompile(`[^\w-]`)
		tag = r2.ReplaceAllString(tag, "")

		r3 := regexp.MustCompile(`-+`)
		tag = r3.ReplaceAllString(tag, "-")

		if tag == "" {
			continue
		}

		tags = append(tags, tag)
	}

	return tags
}
