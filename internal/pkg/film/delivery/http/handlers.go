package http
import (
	"2021_2_MAMBa/internal/pkg/domain"
	"github.com/gorilla/mux"
)

type FilmHandler struct {
	FilmUsecase domain.FilmUsecase
}

func NewHandlers(router *mux.Router, uc domain.FilmUsecase) {
	handler := &FilmHandler{
		FilmUsecase: uc,
	}

	router.HandleFunc("/film/getFilm", handler.GetFilm).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/postRating", handler.PostRating).Methods("POST", "OPTIONS")
	router.HandleFunc("/film/loadFilmReviews", handler.loadFilmReviews).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/loadFilmRecommendations", handler.loadFilmRecommendations).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/loadMyReviewForFilm", handler.LoadMyRv).Methods("GET", "OPTIONS")
}