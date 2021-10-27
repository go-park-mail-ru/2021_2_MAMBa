package server

import (
	collectionsDelivery "2021_2_MAMBa/internal/pkg/collections/delivery/http"
	collectionsRepository "2021_2_MAMBa/internal/pkg/collections/repository"
	collectionsUsecase "2021_2_MAMBa/internal/pkg/collections/usecase"
	"2021_2_MAMBa/internal/pkg/middlewares"
	userDelivery "2021_2_MAMBa/internal/pkg/user/delivery/http"
	userRepository "2021_2_MAMBa/internal/pkg/user/repository"
	userUsecase "2021_2_MAMBa/internal/pkg/user/usecase"
	"2021_2_MAMBa/internal/pkg/utils/log"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func RunServer(addr string) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	// middleware
	api.Use(middlewares.PanicRecovery)
	api.Use(middlewares.Logger)
	api.Use(middlewares.CORS)

	userRepo := userRepository.NewUserRepository()
	collectionsRepo := collectionsRepository.NewCollectionsRepository()

	userUsecase := userUsecase.NewUserUsecase(userRepo)
	collectionsUsecase := collectionsUsecase.NewCollectionsUsecase(collectionsRepo)

	userDelivery.NewHandlers(api, userUsecase)
	collectionsDelivery.NewHandlers(api, collectionsUsecase)

	// Static files
	fileRouter := r.PathPrefix("/static").Subrouter()
	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))
	fileRouter.PathPrefix("/media/").Handler(fileServer)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	fmt.Println("Starting web-server at", addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
