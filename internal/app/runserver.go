package server

import (
	"2021_2_MAMBa/internal/pkg/collections"
	"2021_2_MAMBa/internal/pkg/user"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func RunServer(addr string) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

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
