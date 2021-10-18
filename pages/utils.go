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
