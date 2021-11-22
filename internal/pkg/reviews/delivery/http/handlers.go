package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	authRPC "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	"github.com/gorilla/mux"
)

type ReviewHandler struct {
	ReiviewUsecase domain.ReviewUsecase
	AuthClient authRPC.SessionRPCClient
}

func NewHandlers(router *mux.Router, uc domain.ReviewUsecase, auth authRPC.SessionRPCClient) {
	handler := &ReviewHandler{
		ReiviewUsecase: uc,
		AuthClient: auth,
	}

	router.HandleFunc("/film/getReview", handler.GetReview).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/loadReviewsExcept", handler.LoadExcept).Methods("GET", "OPTIONS")
	router.HandleFunc("/film/postReview", handler.PostReview).Methods("POST", "OPTIONS")

}
