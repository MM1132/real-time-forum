package pages

import (
	fdb "forum/internal/forumDB"
	"forum/internal/forumEnv"
	"net/http"
)

type User struct {
	forumEnv.Env
}

type UserData struct {
	forumEnv.GenericData
	User       fdb.User
	BreadPosts []BreadPost
}

type BreadPost struct {
	Post     fdb.Post
	Thread   fdb.Thread
	Category fdb.Category
}

func (env User) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := UserData{}
	// Here we initialize the data, whatever that means
	if err := data.InitData(env.Env, r); err != nil {
		return
	}

	if r.Method != "GET" {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the id of the thread we return to the client
	idInt, err := GetQueryInt("id", r)
	if err != nil {
		http.NotFound(w, r)
	}

	// Get all the information about the user by its ID
	if user, err := env.Users.Get(idInt); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	} else {
		data.User = user
	}

	// Get all the posts from that user
	if posts, err := env.Posts.GetByUserID(data.User.UserID); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
	} else {
		for i := range posts {
			data.BreadPosts = append(data.BreadPosts, BreadPost{Post: posts[i]})

			// Get the thread where the post is in
			thread, err := env.Threads.Get(data.BreadPosts[i].Post.ThreadID)
			if err != nil {
				sendErr(err, w, http.StatusInternalServerError)
			}
			data.BreadPosts[i].Thread = thread

			if category, err := env.Categories.Get(thread.CategoryID); err != nil {
				sendErr(err, w, http.StatusInternalServerError)
			} else {
				data.BreadPosts[i].Category = category
			}
		}
	}

	// And finally we are executing the template with the data we got
	template := env.Templates["user"]
	if err := template.ExecuteTemplate(w, "layout", data); err != nil {
		sendErr(err, w, http.StatusInternalServerError)
		return
	}
}
