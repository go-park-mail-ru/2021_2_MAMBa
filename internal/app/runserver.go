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
var allowedOrigins = map[string]struct{}{
	"": {},

	"http://localhost":           {},
	"http://localhost:3001":      {},
	"http://localhost:8080":      {},
	"http://89.208.198.137":      {},
	"http://89.208.198.137:3001": {},
	"http://film4u.club":         {},
	"http://film4u.club:3001":    {},

	"https://localhost":           {},
	"https://localhost:3001":      {},
	"https://localhost:8080":      {},
	"https://89.208.198.137":      {},
	"https://89.208.198.137:3001": {},
	"https://film4u.club":         {},
	"https://film4u.club:3001":    {},
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI, "\nHeader: ", r.Header, "\n-------------------------")
		next.ServeHTTP(w, r)
	})
}

func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Recovered from panic with err: " + err.(string))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func CORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		_, isIn := allowedOrigins[origin]
		if isIn {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			// TODO -  на nginx настроить cors и раскомментить
			fmt.Println("unknown origin", `"`+origin+`"`)
			// http.Error(w, `Access denied`, http.StatusForbidden)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH")
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

	// middleware
	api.Use(PanicRecovery)
	api.Use(Logger)
	api.Use(CORS)

	// Users
	api.HandleFunc("/user/{id:[0-9]+}", user.GetBasicInfo).Methods("GET", "OPTIONS")
	api.HandleFunc("/user/register", user.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/user/login", user.Login).Methods("POST", "OPTIONS")
	api.HandleFunc("/user/logout", user.Logout).Methods("GET", "OPTIONS")
	api.HandleFunc("/user/checkAuth", user.CheckAuth).Methods("GET", "OPTIONS")

	// Collections
	api.HandleFunc("/collections/getCollections", collections.GetCollections).Methods("GET", "OPTIONS")

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
