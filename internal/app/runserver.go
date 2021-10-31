package server

import (
	collectionsDelivery "2021_2_MAMBa/internal/pkg/collections/delivery/http"
	collectionsRepository "2021_2_MAMBa/internal/pkg/collections/repository"
	collectionsUsecase "2021_2_MAMBa/internal/pkg/collections/usecase"
	"2021_2_MAMBa/internal/pkg/database"
	filmDelivery "2021_2_MAMBa/internal/pkg/film/delivery/http"
	filmRepository "2021_2_MAMBa/internal/pkg/film/repository"
	filmUsecase "2021_2_MAMBa/internal/pkg/film/usecase"
	"2021_2_MAMBa/internal/pkg/middlewares"
	reviewsDelivery "2021_2_MAMBa/internal/pkg/reviews/delivery/http"
	reviewsRepository "2021_2_MAMBa/internal/pkg/reviews/repository"
	reviewsUsecase "2021_2_MAMBa/internal/pkg/reviews/usecase"
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

	// database
	db := database.Connect()
	defer db.Disconnect()

	userRepo := userRepository.NewUserRepository(db)
	collectionsRepo := collectionsRepository.NewCollectionsRepository(db)
	filmRepo := filmRepository.NewFilmRepository(db)
	reviewRepo := reviewsRepository.NewReviewRepository(db)

	usUsecase := userUsecase.NewUserUsecase(userRepo)
	colUsecase := collectionsUsecase.NewCollectionsUsecase(collectionsRepo)
	filUsecase := filmUsecase.NewFilmUsecase(filmRepo)
	revUsecase := reviewsUsecase.NewReviewUsecase(reviewRepo)

	userDelivery.NewHandlers(api, usUsecase)
	collectionsDelivery.NewHandlers(api, colUsecase)
	filmDelivery.NewHandlers(api, filUsecase)
	reviewsDelivery.NewHandlers(api, revUsecase)

	// Static files
	fileRouter := r.PathPrefix("/static").Subrouter()
	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))
	fileRouter.PathPrefix("/media/").Handler(fileServer)

	server := http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Info(fmt.Sprintf("Starting web-server at %s", addr))
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
