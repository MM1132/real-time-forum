package pages

import (
	fdb "forum/forumDB"
	"forum/forumEnv"
	"log"
	"net/http"
	"strconv"
)

type Forum forumEnv.Env

// Contains things that are generated for every request and passed on to the template
type forumData struct {
	Title       string // Title should be on every page
	ThisCat     fdb.Category
	ChildCats   []fdb.Category
	Breadcrumbs []fdb.Category

	Threads []fdb.Thread
}

// ServeHTTP is called with every request this page receives.
func (e Forum) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := e.Templates["forum"]

	data := forumData{}
	// Read query with key "id"
	query := r.URL.Query()
	queryString := query.Get("id")

	// Then turn it into an int (0 if no query)
	thisCatID := 0
	if queryString != "" {
		var err error
		thisCatID, err = strconv.Atoi(queryString)
		if err != nil {
			http.NotFound(w, r)
			return
		}
	}

	// Get our current category
	thisCat, err := e.Categories.Get(thisCatID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data.Title = thisCat.Name
	data.ThisCat = thisCat

	// And then its children
	childCats, err := e.Categories.GetChildern(thisCat.CategoryID)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	}
	data.ChildCats = childCats

	// Get a list of this category's parents. Outermost first
	tempCat := thisCat
	for tempCat.ParentID.Valid {
		var err error
		tempCat, err = e.Categories.Get(int(tempCat.ParentID.Int64))
		if err != nil {
			sendErr(err, w, http.StatusInternalServerError)
		}
		data.Breadcrumbs = append([]fdb.Category{tempCat}, data.Breadcrumbs...)
	}

	// Get a list of threads in this category
	threads, err := e.Threads.ByCategory(thisCat.CategoryID)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	}
	data.Threads = threads

	// Finally, execute the template with the data we got
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	log.Printf("Served %v to %v\n", data.Title, r.RemoteAddr)
}
