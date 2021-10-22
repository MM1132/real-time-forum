package pages

import (
	"log"
	"net/http"
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
