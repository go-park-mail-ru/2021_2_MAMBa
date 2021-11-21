package server

import (
	grpcCollection "2021_2_MAMBa/internal/pkg/collections/delivery/grpc"
	collectionsDelivery "2021_2_MAMBa/internal/pkg/collections/delivery/http"
	"2021_2_MAMBa/internal/pkg/database"
	filmDelivery "2021_2_MAMBa/internal/pkg/film/delivery/http"
	filmRepository "2021_2_MAMBa/internal/pkg/film/repository"
	filmUsecase "2021_2_MAMBa/internal/pkg/film/usecase"
	"2021_2_MAMBa/internal/pkg/middlewares"
	personDelivery "2021_2_MAMBa/internal/pkg/person/delivery/http"
	personRepository "2021_2_MAMBa/internal/pkg/person/repository"
	personUsecase "2021_2_MAMBa/internal/pkg/person/usecase"
	reviewsDelivery "2021_2_MAMBa/internal/pkg/reviews/delivery/http"
	reviewsRepository "2021_2_MAMBa/internal/pkg/reviews/repository"
	reviewsUsecase "2021_2_MAMBa/internal/pkg/reviews/usecase"
	userDelivery "2021_2_MAMBa/internal/pkg/user/delivery/http"
	userRepository "2021_2_MAMBa/internal/pkg/user/repository"
	userUsecase "2021_2_MAMBa/internal/pkg/user/usecase"
	"2021_2_MAMBa/internal/pkg/utils/log"
	"fmt"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"net/http"
)

func RunServer(addr string, collAddr string) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	// middleware
	api.Use(middlewares.PanicRecovery)
	api.Use(middlewares.Logger)
	api.Use(middlewares.CORS)
	//api.Use(middlewares.CSRF)

	// database
	db := database.Connect()
	defer db.Disconnect()

	userRepo := userRepository.NewUserRepository(db)
	filmRepo := filmRepository.NewFilmRepository(db)
	personRepo := personRepository.NewPersonRepository(db)
	reviewRepo := reviewsRepository.NewReviewRepository(db)

	conn, err := grpc.Dial("localhost:"+collAddr, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()
	clientCollections := grpcCollection.NewCollectionsRPCClient(conn)

	usUsecase := userUsecase.NewUserUsecase(userRepo)
	filUsecase := filmUsecase.NewFilmUsecase(filmRepo)
	persUsecase := personUsecase.NewPersonUsecase(personRepo)
	revUsecase := reviewsUsecase.NewReviewUsecase(reviewRepo)

	userDelivery.NewHandlers(api, usUsecase)
	collectionsDelivery.NewHandlers(api, clientCollections)
	filmDelivery.NewHandlers(api, filUsecase)
	personDelivery.NewHandlers(api, persUsecase)
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
	err = server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}
