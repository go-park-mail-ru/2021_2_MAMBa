package server

import (
	"2021_2_MAMBa/internal/pkg/collections"
	"2021_2_MAMBa/internal/pkg/user"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)
// TODO - add all desirable origins
var allowedOrigins = map[string]struct{} {
	"https://localhost":{},
	"http://localhost":{},
}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		_, isIn := allowedOrigins[origin]
		if isIn {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			http.Error(w, `Access denied`, http.StatusForbidden)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-language, Content-Type, Content-Language, Content-Encoding")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}


func RunServer(addr string) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.Use(CORS)
	// Users
	api.HandleFunc("/user/{id:[0-9]+}", user.GetBasicInfo).Methods("GET")
	api.HandleFunc("/user/register", user.Register).Methods("POST")
	api.HandleFunc("/user/login", user.Login).Methods("POST")
	api.HandleFunc("/user/logout", user.Logout).Methods("GET")
	api.HandleFunc("/user/checkAuth", user.CheckAuth).Methods("GET")

	// Collections
	api.HandleFunc("/collections/getCollections", collections.GetCollections).Methods("GET")

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	fmt.Println("Starting web-server at", addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
