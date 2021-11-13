package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"2021_2_MAMBa/internal/pkg/utils/xss"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	defaultLimit = 10
	defaultSkip  = 0
)

func (handler *UserHandler) GetBasicInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	u64, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	us, err := handler.UserUsecase.GetBasicInfo(u64)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusNotFound}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(us)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	targetID, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	clientID, err := sessions.CheckSession(r)
	if err != nil && err != customErrors.ErrorUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorInternalServer.Error()), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	us, err := handler.UserUsecase.GetProfileById(clientID, targetID)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusNotFound}
		resp.Write(w)
		return
	}

	x, err := json.Marshal(us)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	clientID, err := sessions.CheckSession(r)
	if err != nil || err == customErrors.ErrorUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorInternalServer.Error()), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	profileForm := domain.Profile{}
	err = json.NewDecoder(r.Body).Decode(&profileForm)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	xss.SanitizeProfile(&profileForm)
	profileForm.ID = clientID

	us, err := handler.UserUsecase.UpdateProfile(profileForm)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusNotFound}
		resp.Write(w)
		return
	}

	x, err := json.Marshal(us)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *UserHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	clientID, err := sessions.CheckSession(r)
	if err != nil || err == customErrors.ErrorUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorInternalServer.Error()), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	profileForm := domain.Profile{}
	err = json.NewDecoder(r.Body).Decode(&profileForm)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	us, err := handler.UserUsecase.CreateSubscription(clientID, profileForm.ID)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusNotFound}
		resp.Write(w)
		return
	}

	x, err := json.Marshal(us)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *UserHandler) LoadUserReviews(w http.ResponseWriter, r *http.Request) {
	var err error
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	skip, err := queryChecker.CheckIsIn(w, r, "skip", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	limit, err := queryChecker.CheckIsIn(w, r, "limit", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrLimitMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	review, err := handler.UserUsecase.LoadUserReviews(id, skip, limit)
	if err == customErrors.ErrorSkip {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(review)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
