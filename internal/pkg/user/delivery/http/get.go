package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

	targetIDString, isIn := r.URL.Query()["id"]
	var targetID uint64
	var err error
	if isIn {
		targetID, err = strconv.ParseUint(targetIDString[0], 10, 64)
		if err != nil {
			http.Error(w, customErrors.ErrSkipMsg, http.StatusBadRequest)
			return
		}
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
	var id uint64
	skip, limit := 0, 10
	var err error
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, customErrors.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	skipString, isIn := r.URL.Query()["skip"]
	if isIn {
		skip, err = strconv.Atoi(skipString[0])
		if err != nil || skip  < 0 {
			http.Error(w, customErrors.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	limitString, isIn := r.URL.Query()["limit"]
	if isIn {
		limit, err = strconv.Atoi(limitString[0])
		if err != nil || limit <= 0 {
			http.Error(w, customErrors.ErrLimitMsg, http.StatusBadRequest)
			return
		}
	}
	review, err := handler.UserUsecase.LoadUserReviews(id, skip, limit)
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
