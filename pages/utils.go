package pages

import (
	"log"
	"net/http"
)

func sendErr(err error, w http.ResponseWriter, code int) {
	log.Println(err)
	http.Error(w, http.StatusText(code), code)
}

func recoverHandler(w http.ResponseWriter) {
	if err := recover(); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}
