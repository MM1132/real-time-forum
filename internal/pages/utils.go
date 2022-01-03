package pages

import (
	"encoding/json"
	"errors"
	"fmt"
	"forum/internal/forumDB"
	"forum/internal/forumEnv"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

// Sends an error on the writer with default text
func sendErr(err error, w http.ResponseWriter, code int) {
	log.Println(err)
	http.Error(w, http.StatusText(code), code)
}

// Get the id of the thread we return to the client
func GetQueryInt(key string, r *http.Request) (int, error) {
	idString := r.URL.Query().Get(key)
	// Then we try turning the id into an integer
	if idString == "" {
		return 0, errors.New("no id given")
	}
	if id, err := strconv.Atoi(idString); err != nil {
		return 0, errors.New("could not convert idString")
	} else {
		return id, nil
	}
}

// Convert any golang object into json and send it to the user
func sendInterfaceAsJson(w http.ResponseWriter, response interface{}) error {
	jsonData, err := json.Marshal(response)
	if err != nil {
		return err
	}
	w.Write(jsonData)
	return nil
}

// I guess this function is used to determine whether the user is logged in or not
func checkUser(data forumEnv.GenericData, addr string) error {
	if data.User.UserID == 0 {
		return fmt.Errorf("%v not authorized", addr)
	}
	return nil
}

func writePost(content string, userID int, threadID int, env Thread) (int, error) {
	newPost := forumDB.Post{
		Content:  content,
		UserID:   userID,
		ThreadID: threadID,
	}

	id, err := env.Posts.Insert(newPost)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Used to turn passwords into hashes
func generateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
