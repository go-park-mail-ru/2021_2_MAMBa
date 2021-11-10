package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/xss"
	"encoding/json"
	"net/http"
)

func (handler *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userForm := new(domain.User)
	err := json.NewDecoder(r.Body).Decode(&userForm)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	xss.SanitizeUser(userForm)

	us, err := handler.UserUsecase.Register(userForm)
	if err == customErrors.ErrorAlreadyExists {
		resp := domain.Response{Body: cast.ErrorToJson(err.Error()), Status: http.StatusConflict}
		resp.Write(w)
		return
	}
	if err == customErrors.ErrorBadInput {
		resp := domain.Response{Body: cast.ErrorToJson(err.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err !=nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}

	err = sessions.StartSession(w, r, us.ID)
	if err != nil && us.ID != 0 {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}

	x, err := json.Marshal(us)
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusCreated,
	}
	resp.Write(w)
}

func (handler *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userForm := new(domain.UserToLogin)
	err := json.NewDecoder(r.Body).Decode(&userForm)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	_, err = sessions.CheckSession(r)
	if err != customErrors.ErrorUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorAlreadyIn.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}

	us, err := handler.UserUsecase.Login(userForm)
	if err == customErrors.ErrorBadInput {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err == customErrors.ErrorBadCredentials {
		resp := domain.Response{Body: cast.ErrorToJson(err.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}

	err = sessions.StartSession(w, r, us.ID)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
	}

	x, err := json.Marshal(us)
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	id, err := sessions.CheckSession(r)
	if err == customErrors.ErrorUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusForbidden}
		resp.Write(w)
		return
	}
	err = sessions.EndSession(w, r, id)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) CheckAuth(w http.ResponseWriter, r *http.Request) {
	userID, err := sessions.CheckSession(r)
	if err == customErrors.ErrorUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusForbidden}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}

	us, err := handler.UserUsecase.CheckAuth(userID)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}

	x, err := json.Marshal(us)
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
