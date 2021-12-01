package forumEnv

import (
	"log"
	"net/http"
)

func Log(in http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client := r.Header.Get("X-Forwarded-For")
		if client == "" {
			client = r.RemoteAddr
		}
		log.Printf("Request for %v from %v\n", r.URL.RequestURI(), client)
		in.ServeHTTP(w, r)
	})
}

func RedirectEmpty(redirPath string, handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, redirPath, http.StatusFound)
			return
		}

		handler.ServeHTTP(w, r)
	})
}
