package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"2021_2_MAMBa/internal/pkg/utils/xss"
	"encoding/json"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip  = 0
)

func (handler *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	var err error
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	review, err := handler.ReiviewUsecase.GetReview(id)
	if err == customErrors.ErrorSkip {
		resp := domain.Response{Error: cast.StringToJson(err.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(review)
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *ReviewHandler) PostReview(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reviewForm := domain.Review{}
	err := json.NewDecoder(r.Body).Decode(&reviewForm)
	if err != nil {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	xss.SanitizeReview(&reviewForm)

	authId, err := sessions.CheckSession(r)
	if err != nil {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrUserNotLoggedIn.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}
	reviewForm.AuthorId = authId
	review, err := handler.ReiviewUsecase.PostReview(reviewForm)
	x, err := json.Marshal(review)
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *ReviewHandler) LoadExcept(w http.ResponseWriter, r *http.Request) {
	var err error
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	filmId, err := queryChecker.CheckIsIn64(w, r, "film_id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	skip, err := queryChecker.CheckIsIn(w, r, "skip", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	limit, err := queryChecker.CheckIsIn(w, r, "limit", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrLimitMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	review, err := handler.ReiviewUsecase.LoadReviewsExcept(id, filmId, skip, limit)
	if err == customErrors.ErrorSkip {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Error: cast.StringToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(review)
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
