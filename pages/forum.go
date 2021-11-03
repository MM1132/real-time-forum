package pages

import (
	fdb "forum/forumDB"
	"forum/forumEnv"
	"net/http"
	"strconv"
)

type Forum struct {
	forumEnv.Env
}

// Contains things that are generated for every request and passed on to the template
type forumData struct {
	forumEnv.GenericData
	ThisCat     fdb.Category
	ChildCats   []fdb.Category
	Breadcrumbs []fdb.Category

	Threads []fdb.Thread
}

// ServeHTTP is called with every request this page receives.
func (env Forum) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := forumData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}

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
	thisCat, err := env.Categories.Get(thisCatID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	data.AddTitle(thisCat.Name)
	data.ThisCat = thisCat

	// And then its children
	childCats, err := env.Categories.GetChildern(thisCat.CategoryID)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	}
	data.ChildCats = childCats

	// BREAD-CRUMBS!!!!
	// Get a list of this category's parents. Outermost first
	tempCat := thisCat
	for tempCat.ParentID.Valid {
		var err error
		tempCat, err = env.Categories.Get(int(tempCat.ParentID.Int64))
		if err != nil {
			sendErr(err, w, http.StatusInternalServerError)
		}
		data.Breadcrumbs = append([]fdb.Category{tempCat}, data.Breadcrumbs...)
	}

	// Get a list of threads in this category
	threads, err := env.Threads.ByCategory(thisCat.CategoryID)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	}
	data.Threads = threads

	// Finally, execute the template with the data we got
	tmpl := env.Templates["forum"]
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}
