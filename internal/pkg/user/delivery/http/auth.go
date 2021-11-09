package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
	"2021_2_MAMBa/internal/pkg/utils/xss"
	"encoding/json"
	"net/http"
)

func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userForm := new(domain.User)
	err := json.NewDecoder(r.Body).Decode(&userForm)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	xss.SanitizeUser(userForm)

	us, err := handler.UserUsecase.Register(userForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sessions.StartSession(w, r, us.ID)
	if err != nil && us.ID != 0 {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(b)
}

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userForm := new(domain.UserToLogin)
	err := json.NewDecoder(r.Body).Decode(&userForm)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}

	_, err = sessions.CheckSession(r)
	if err != customErrors.ErrUserNotLoggedIn {
		http.Error(w, customErrors.ErrorAlreadyIn.Error(), http.StatusBadRequest)
		return
	}

	us, err := handler.UserUsecase.Login(userForm)
	if err == customErrors.ErrorBadInput {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err == customErrors.ErrorBadCredentials {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	err = sessions.StartSession(w, r, us.ID)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
}

func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	id, err := sessions.CheckSession(r)
	if err == customErrors.ErrUserNotLoggedIn {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusForbidden)
		return
	}
	err = sessions.EndSession(w, r, id)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	userID, err := sessions.CheckSession(r)
	if err == customErrors.ErrUserNotLoggedIn {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	us, err := handler.UserUsecase.CheckAuth(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
}
