package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"github.com/gorilla/mux"
)

type CollectionsHandler struct {
	CollectionsUsecase domain.CollectionsUsecase
}

func NewHandlers(router *mux.Router, uc domain.CollectionsUsecase) {
	handler := &CollectionsHandler{
		CollectionsUsecase: uc,
	}

	router.HandleFunc("/collections/getCollections", handler.GetCollections).Methods("GET", "OPTIONS")
}
