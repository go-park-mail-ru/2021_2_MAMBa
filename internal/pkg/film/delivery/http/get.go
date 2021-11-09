package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip  = 0
)

func (handler *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	var err error
	// default
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	skipReview, err := queryChecker.CheckIsIn(w, r, "skip_reviews", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	limitReview, err := queryChecker.CheckIsIn(w, r, "limit_reviews", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrLimitMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	skipRecom, err := queryChecker.CheckIsIn(w, r, "skip_recommend", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	limitRecom, err := queryChecker.CheckIsIn(w, r, "limit_recommend", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrLimitMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	filmPageInfo, err := handler.FilmUsecase.GetFilm(uint64(id), skipReview, limitReview, skipRecom, limitRecom)
	if err == customErrors.ErrorSkip {
		resp := domain.Response{Body: cast.ErrorToJson(err.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(filmPageInfo)
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusOK,
  }
	resp.Write(w)
}

func (handler *FilmHandler) PostRating(w http.ResponseWriter, r *http.Request) {
	authId, err := sessions.CheckSession(r)
	if err == customErrors.ErrorUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	rating, err := queryChecker.CheckIsInFloat64(w, r, "rating", 0, errors.New(customErrors.ErrEncMsg))
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrRateMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	newRating, err := handler.FilmUsecase.PostRating(id, authId, rating)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	jsRate := json.Number(fmt.Sprintf("%.1f", newRating))
	x, err := json.Marshal(domain.NewRate{Rating: jsRate})
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusOK,
  }
	resp.Write(w)
}

func (handler *FilmHandler) LoadMyRv(w http.ResponseWriter, r *http.Request) {
	authId, err := sessions.CheckSession(r)
	if err == customErrors.ErrorUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	review, err := handler.FilmUsecase.LoadMyReview(id, authId)
	if err == customErrors.ErrorNoReviewForFilm {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrNoReviewMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
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

func (handler *FilmHandler) loadFilmReviews(w http.ResponseWriter, r *http.Request) {
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

	reviews, err := handler.FilmUsecase.LoadFilmReviews(id, skip, limit)
	if err == customErrors.ErrorSkip {
		resp := domain.Response{Body: cast.ErrorToJson(err.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(reviews)
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) loadFilmRecommendations(w http.ResponseWriter, r *http.Request) {
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
	recommendations, err := handler.FilmUsecase.LoadFilmRecommendations(id, skip, limit)
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
	x, err := json.Marshal(recommendations)
	resp:= domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
