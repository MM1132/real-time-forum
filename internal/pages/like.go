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

func (env Like) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		err := fmt.Errorf("invalid request to %v from %v", r.RequestURI, r.RemoteAddr)
		sendErr(fmt.Errorf(`like: %w`, err), w, http.StatusMethodNotAllowed)
		return
	}

	data := forumEnv.GenericData{}
	if err := data.InitData(env.Env, r); err != nil {
		sendErr(fmt.Errorf(`like: %w`, err), w, http.StatusInternalServerError)
		return
	}

	err := checkUser(data, r.RemoteAddr)
	if err != nil {
		sendErr(fmt.Errorf(`like: %w`, err), w, http.StatusForbidden)
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
		http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
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
