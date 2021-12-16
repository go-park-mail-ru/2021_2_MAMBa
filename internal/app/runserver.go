package server

import (
	"2021_2_MAMBa/internal/app/config"
	grpcCollection "2021_2_MAMBa/internal/pkg/collections/delivery/grpc"
	collectionsDelivery "2021_2_MAMBa/internal/pkg/collections/delivery/http"
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
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
	searchDelivery "2021_2_MAMBa/internal/pkg/search/delivery/http"
	searchRepository "2021_2_MAMBa/internal/pkg/search/repository"
	searchUsecase "2021_2_MAMBa/internal/pkg/search/usecase"
	grpcAuth "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	userDelivery "2021_2_MAMBa/internal/pkg/user/delivery/http"
	userRepository "2021_2_MAMBa/internal/pkg/user/repository"
	userUsecase "2021_2_MAMBa/internal/pkg/user/usecase"
	"2021_2_MAMBa/internal/pkg/utils/log"
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"math"
	"net/http"
	"sync"
	"time"
)

func RunServer(configPath string) {
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	cfg := config.ParseMain(configPath)

	// middleware
	middlewares.RegisterMetrics()
	metrics := middlewares.InitMetrics()
	api.Use(middlewares.PanicRecovery)
	api.Use(metrics.Metrics)
	api.Use(middlewares.Logger)
	if cfg.CORS {
		api.Use(middlewares.CORS)
	}
	if cfg.CSRF {
		api.Use(middlewares.CSRF)
	}

	// database
	db := database.Connect(cfg.Db)
	defer db.Disconnect()

	userRepo := userRepository.NewUserRepository(db)
	filmRepo := filmRepository.NewFilmRepository(db)
	personRepo := personRepository.NewPersonRepository(db)
	reviewRepo := reviewsRepository.NewReviewRepository(db)
	searchRepo := searchRepository.NewSearchRepository(db)

	collConn, err := grpc.Dial("localhost:"+cfg.CollectPort, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer collConn.Close()
	clientCollections := grpcCollection.NewCollectionsRPCClient(collConn)
	authConn, err := grpc.Dial("localhost:"+cfg.AuthPort, grpc.WithInsecure())
	if err != nil {
		return
	}
	defer authConn.Close()
	clientAuth := grpcAuth.NewSessionRPCClient(authConn)

	usUsecase := userUsecase.NewUserUsecase(userRepo)
	filUsecase := filmUsecase.NewFilmUsecase(filmRepo)
	persUsecase := personUsecase.NewPersonUsecase(personRepo)
	revUsecase := reviewsUsecase.NewReviewUsecase(reviewRepo)
	searUsecase := searchUsecase.NewSearchUsecase(searchRepo, personRepo, filmRepo)

	userDelivery.NewHandlers(api, usUsecase, clientAuth)
	collectionsDelivery.NewHandlers(api, clientCollections)
	filmDelivery.NewHandlers(api, filUsecase, clientAuth)
	personDelivery.NewHandlers(api, persUsecase)
	reviewsDelivery.NewHandlers(api, revUsecase, clientAuth)
	searchDelivery.NewHandlers(api, searUsecase, clientAuth)
	r.Handle("/metrics", promhttp.Handler())

	go notificationWorker(filmRepo)

	// Static files
	fileRouter := r.PathPrefix("/static").Subrouter()
	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))
	fileRouter.PathPrefix("/media/").Handler(fileServer)

	server := http.Server{
		Addr:    ":" + cfg.ListenPort,
		Handler: r,
	}

	log.Info(fmt.Sprintf("Starting web-server at %s", cfg.ListenPort))
	err = server.ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}

func notificationWorker(filmRepo domain.FilmRepository) {
	// storage for films released today
	comingFilms := struct {
		films []domain.Film
		sync.RWMutex
	}{}

	// daemon to update coming this month films
	go func(filmRepo domain.FilmRepository) {
		for {
			comingFilms.Lock()
			comingFilms.films = []domain.Film{}
			comingFilms.Unlock()

			year, month, _ := time.Now().Date()
			filmsBuffer, err := filmRepo.GetFilmsByMonthYear(int(month), year, math.MaxInt, 0)
			if err == nil {
				for _, v := range filmsBuffer.FilmList {
					if time.Now().Format("2006-01-02") == v.PremiereRu {
						comingFilms.Lock()
						comingFilms.films = append(comingFilms.films, v)
						comingFilms.Unlock()
					}
				}
				log.Info(fmt.Sprintf("Found %d films released today", len(comingFilms.films)))
			}
			if err != nil {
				// retry delay
				log.Warn("Error in updating coming this month films")
				time.Sleep(1 * time.Minute)
				continue
			}
			time.Sleep(24 * time.Hour)
		}
	}(filmRepo)

	// preparing firebase messages sender
	opt := option.WithCredentialsFile("firebasePrivateKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Warn(fmt.Sprintf("error initializing Firebase app: %v\n", err))
	}
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Warn(fmt.Sprintf("error getting Firebase Messaging client: %v\n", err))
	}

	// then we send notifications for all films that were released today
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			hr, min, _ := time.Now().Clock()
			if hr == 12 && min == 0 {
				comingFilms.RLock()
				for _, v := range comingFilms.films {
					message := &messaging.Message{
						Notification: &messaging.Notification{
							Title: "Сегодня вышел в прокат фильм",
							Body:  v.Title,
						},
						Topic: "all",
					}
					response, err := client.Send(ctx, message)
					if err != nil {
						log.Error(err)
						continue
					}
					log.Info(fmt.Sprintf("Successfully sent message: %v, for film id: %d", response, v.Id))
				}
				comingFilms.RUnlock()
			}
		}
	}

}
