package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"2021_2_MAMBa/internal/pkg/utils/xss"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	defaultLimit = 10
	defaultSkip = 0
)

func (handler *UserHandler) GetBasicInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	u64, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	us, err := handler.UserUsecase.GetBasicInfo(u64)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusNotFound)
		return
	}
	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	targetID, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	clientID, err := sessions.CheckSession(r)
	if err != nil && err != sessions.ErrUserNotLoggedIn {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	us, err := handler.UserUsecase.GetProfileById(clientID, targetID)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusNotFound)
		return
	}
	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	clientID, err := sessions.CheckSession(r)
	if err != nil || err == sessions.ErrUserNotLoggedIn {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	profileForm := domain.Profile{}
	err = json.NewDecoder(r.Body).Decode(&profileForm)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}

	xss.SanitizeProfile(&profileForm)
	profileForm.ID = clientID

	us, err := handler.UserUsecase.UpdateProfile(profileForm)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	clientID, err := sessions.CheckSession(r)
	if err != nil || err == sessions.ErrUserNotLoggedIn {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	profileForm := domain.Profile{}
	err = json.NewDecoder(r.Body).Decode(&profileForm)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}

	us, err := handler.UserUsecase.CreateSubscription(clientID, profileForm.ID)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *UserHandler) LoadUserReviews(w http.ResponseWriter, r *http.Request) {
	var err error
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	skip, err := queryChecker.CheckIsIn(w,r, "skip", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	limit, err := queryChecker.CheckIsIn(w, r, "limit", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		return

	}
	review, err := handler.UserUsecase.LoadUserReviews(id, skip, limit)
	if err == customErrors.ErrorSkip {
		http.Error(w, customErrors.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, customErrors.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(review)
	if err != nil {
		http.Error(w, customErrors.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}
