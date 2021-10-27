package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/sessions"
	"2021_2_MAMBa/internal/pkg/user"
	"encoding/json"
	"net/http"
)

func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userForm := new(domain.User)
	err := json.NewDecoder(r.Body).Decode(&userForm)
	if err != nil {
		http.Error(w, user.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}

	us, err := handler.UserUsecase.Register(userForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sessions.StartSession(w, r, us.ID)
	if err != nil && us.ID != 0 {
		http.Error(w, user.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, user.ErrorInternalServer.Error(), http.StatusInternalServerError)
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
		http.Error(w, user.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}

	us, err := handler.UserUsecase.Login(userForm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = sessions.CheckSession(r)
	if err != sessions.ErrUserNotLoggedIn {
		http.Error(w, user.ErrorAlreadyIn.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, user.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	err = sessions.StartSession(w, r, us.ID)
	if err != nil {
		http.Error(w, user.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
}

func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	id, err := sessions.CheckSession(r)
	if err == sessions.ErrUserNotLoggedIn {
		http.Error(w, user.ErrorBadInput.Error(), http.StatusForbidden)
		return
	}
	err = sessions.EndSession(w, r, id)
	if err != nil {
		http.Error(w, user.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	userID, err := sessions.CheckSession(r)
	if err == sessions.ErrUserNotLoggedIn {
		http.Error(w, user.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	us, err := handler.UserUsecase.CheckAuth(userID)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(us)
	if err != nil {
		http.Error(w, user.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
}
