package middlewares

import (
	"log"
	"net/http"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, "\nHeader: ", r.Header, "\n-------------------------")
		next.ServeHTTP(w, r)
	})
}
