package pages

import (
	"forum/internal/forumEnv"
	"forum/internal/search"
	"net/http"
	"net/url"
	"regexp"
)

type Search struct {
	forumEnv.Env
}

type searchData struct {
	forumEnv.GenericData
	Search search.Searcher
}

func (env Search) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := searchData{}
	data.InitData(env.Env, r)

	// Convert post request fields to query string
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			sendErr(err, w, http.StatusInternalServerError)
			return
		}

		var kvp []string
		for k, v := range r.Form {
			kvp = append(kvp, k, v[0])
		}

		uri, err := data.Query(kvp...)
		if err != nil {
			sendErr(err, w, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, string(uri), http.StatusSeeOther)
		return
	}

	searchStr := r.FormValue("search")
	searchType := r.FormValue("type")

	// Get type of search if it's set in the search string
	reg := regexp.MustCompile(` ?type\((.+)\)`)
	if match := reg.FindStringSubmatch(searchStr); match != nil {
		searchType = match[1]
		searchStr = reg.ReplaceAllString(searchStr, "")

		uri, err := data.Query("type", searchType, "search", searchStr)
		if err != nil {
			goto SKIP
		}

		newUrl, err := url.Parse(string(uri))
		if err != nil {
			goto SKIP
		}

		data.CurrentURL = *newUrl
	}
SKIP:

	// Set searcher in data, of the type we want
	switch searchType {
	case "thread":
		data.Search = search.NewThreadSearch(env.Env, data.CurrentURL, data.User.UserID)

	default:
		fallthrough // Default to post
	case "post":
		postSearch := search.NewPostSearch(env.Env, data.CurrentURL, data.User.UserID)
		postSearch.Breadcrumbs = true
		data.Search = postSearch
	}

	data.Search.ProcessString(searchStr)
	data.Search.ProcessRequest(r)

	if err := data.Search.DoSearch(env.Env); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}

	// And finally we are executing the template with the data we got
	template := env.Templates["search"]
	if err := template.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}
