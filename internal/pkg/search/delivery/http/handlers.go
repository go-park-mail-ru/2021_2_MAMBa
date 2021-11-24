package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	authRPC "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	"github.com/gorilla/mux"
)

type SearchHandler struct {
	SearchUsecase domain.SearchUsecase
	AuthClient authRPC.SessionRPCClient
}

func NewHandlers(router *mux.Router, uc domain.SearchUsecase, auth authRPC.SessionRPCClient) {
	handler := &SearchHandler{
		SearchUsecase: uc,
		AuthClient: auth,
	}

	router.HandleFunc("/search", handler.GetSearch).Methods("GET", "OPTIONS")

}
