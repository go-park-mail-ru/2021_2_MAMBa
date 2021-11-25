package http

import (
	"2021_2_MAMBa/internal/pkg/collections/delivery/grpc"
	"github.com/gorilla/mux"
)

type CollectionsHandler struct {
	CollectionsClient grpc.CollectionsRPCClient
}

func NewHandlers(router *mux.Router, cl grpc.CollectionsRPCClient) {
	handler := &CollectionsHandler{
		CollectionsClient: cl,
	}

	router.HandleFunc("/collections/getCollections", handler.GetCollections).Methods("GET", "OPTIONS")
	router.HandleFunc("/collections/getCollectionFilms", handler.GetCollectionFilms).Methods("GET", "OPTIONS")
}
