package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	authRPC "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	"github.com/gorilla/mux"
)

type FilmHandler struct {
	FilmUsecase domain.FilmUsecase
	AuthClient  authRPC.SessionRPCClient
}

func NewHandlers(router *mux.Router, uc domain.FilmUsecase, auth authRPC.SessionRPCClient) {
	handler := &FilmHandler{
		FilmUsecase: uc,
		AuthClient:  auth,
	}

	router.HandleFunc("/film/getFilm", handler.GetFilm).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/postRating", handler.PostRating).Methods("POST", "OPTIONS")
	router.HandleFunc("/film/loadFilmReviews", handler.loadFilmReviews).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/loadFilmRecommendations", handler.loadFilmRecommendations).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/loadMyReviewForFilm", handler.LoadMyRv).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/postBookmark", handler.BookmarkFilm).Methods("POST", "OPTIONS")
	router.HandleFunc("/film/calendar", handler.GetFilmsByMonthYear).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/genres", handler.GetGenres).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/genreFilms", handler.GetFilmsByGenre).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/calendar", handler.GetFilmsByMonthYear).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/banners", handler.GetBanners).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/popularFilms", handler.GetPopularFilms).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/getBookmarks", handler.LoadUserBookmarks).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/getRandom", handler.GetRandomFilms).Methods("GET", "OPTIONS")
}
