package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"net/http"
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
	router.HandleFunc("/user/getProfile", handler.GetProfile).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/changeProfile", handler.UpdateProfile).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/subscribeTo", handler.CreateSubscription).Methods("POST", "OPTIONS")
	router.HandleFunc("/user/getReviewsAndStars", handler.LoadUserReviews).Methods("GET", "OPTIONS")
	router.HandleFunc("/user/avatar", handler.UploadAvatar).Methods("POST", "OPTIONS")

	router.HandleFunc("/csrf", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		w.WriteHeader(http.StatusNoContent)
	}).Methods("GET", "OPTIONS")
}
