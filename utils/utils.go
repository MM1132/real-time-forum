package utils

import (
	"log"
	"net/http"
)

func FatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func SendErr(err error, w http.ResponseWriter, code int) {
	log.Println(err)
	http.Error(w, http.StatusText(code), code)
}
