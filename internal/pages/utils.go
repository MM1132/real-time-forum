package pages

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

// Sends an error on the writer with default text
func sendErr(err error, w http.ResponseWriter, code int) {
	log.Println(err)
	http.Error(w, http.StatusText(code), code)
}

// Recover from panicking goroutine if it's a handler.
func recoverHandler(w http.ResponseWriter) {
	if err := recover(); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
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
