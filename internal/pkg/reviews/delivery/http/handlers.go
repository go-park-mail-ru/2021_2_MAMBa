package http
import (
	"2021_2_MAMBa/internal/pkg/domain"
	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	ReiviewUsecase domain.ReviewUsecase
}

func NewHandlers(router *mux.Router, uc domain.ReviewUsecase) {
	handler := &ReviewHandler{
		ReiviewUsecase: uc,
	}

	router.HandleFunc("/film/getReview", handler.GetReview).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/loadReviewsExcept", handler.LoadExcept).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/postReview", handler.PostReview).Methods("POST", "OPTIONS")

}