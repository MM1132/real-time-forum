package utils

import (
	"log"
	"net/http"
)

func FatalErr(err error) {
	if err != nil {
		panic(err)
	}
}

func SendErr(err error, w http.ResponseWriter, code int) {
	log.Println(err)
	http.Error(w, http.StatusText(code), code)
}

func Recover(err error, w http.ResponseWriter) {
	if recover() != nil {
		SendErr(err, w, http.StatusInternalServerError)
	}
}
