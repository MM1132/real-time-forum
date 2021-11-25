package pages

import (
	"log"
	"net/http"
)

func (env Board) ServeTagResults(w http.ResponseWriter, r *http.Request) {
	data := boardData{}
	if err := data.InitData(env.Env, r); err != nil {
		return
	}

	tag := r.URL.Query().Get("tag")

	// Get a threads page, except by a tag this time
	threadsPage, err := env.GetThreadsPage(tag, r)
	if err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
	data.ThreadsPage = threadsPage

	// Finally, execute the template with the data we got
	tmpl := env.Templates["board"]
	if err := tmpl.ExecuteTemplate(w, "layout", data); err != nil {
		log.Print(err)
		return
	}
}
