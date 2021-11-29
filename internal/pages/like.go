package pages

import (
	"database/sql"
	"fmt"
	"forum/internal/forumEnv"
	"net/http"
)

type Like struct {
	forumEnv.Env
}

type LikeInfo struct {
	LikeCount int
}

func (env Like) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		err := fmt.Errorf("invalid request to %v from %v", r.RequestURI, r.RemoteAddr)
		sendErr(fmt.Errorf(`like: %w`, err), w, http.StatusMethodNotAllowed)
		return
	}

	data := forumEnv.GenericData{}
	data.InitData(env.Env, r)

	err := checkUser(data, r.RemoteAddr)
	if err != nil {
		sendInterfaceAsJson(w, Redirect{RedirectPath: "/login"})
		// http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID, err := GetQueryInt("id", r)
	if err != nil {
		sendErr(fmt.Errorf(`like: %w`, err), w, http.StatusBadRequest)
		return
	}

	switch r.URL.Path {
	case "/like":
		err = env.LikeDislike(postID, data, 1)
	case "/dislike":
		err = env.LikeDislike(postID, data, -1)
	}
	if err != nil {
		sendErr(fmt.Errorf(`like: %w`, err), w, http.StatusInternalServerError)
		return
	} else {
		// Liking was successful
		// Get the number of likes at the moment
		likeCount, err := env.Likes.GetPostTotal(postID)
		if err != nil {
			sendErr(fmt.Errorf(`like: %w`, err), w, http.StatusInternalServerError)
			return
		}
		sendInterfaceAsJson(w, LikeInfo{LikeCount: likeCount})
	}
}

func (env Like) LikeDislike(postID int, data forumEnv.GenericData, value int) error {
	like, err := env.Likes.Get(postID, data.User.UserID)

	// If no (dis)like by the user is found, create a new one
	if err == sql.ErrNoRows {
		err := env.Likes.Insert(postID, data.Session.UserID, value)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil { // Something has gone horribly wrong
		return err
	}

	// If one exists for this post, modify it
	if like.Value == value { // If clicked twice
		err := env.Likes.Delete(like.PostID, like.UserID)
		if err != nil {
			return err
		}
		return nil
	} else { // If clicked liked but after clicked dislike
		err = env.Likes.Update(like.PostID, like.UserID, value)
		if err != nil {
			return err
		}
	}
	return nil
}
