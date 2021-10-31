package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	UserUsecase domain.UserUsecase
}

func NewHandlers(router *mux.Router, uc domain.UserUsecase) {
	handler := &UserHandler{
		UserUsecase: uc,
	}

	router.HandleFunc("/user/{id:[0-9]+}", handler.GetBasicInfo).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/register", handler.Register).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/login", handler.Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/logout", handler.Logout).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/checkAuth", handler.CheckAuth).Methods("GET", "OPTIONS")
}