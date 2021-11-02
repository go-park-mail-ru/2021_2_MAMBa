package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip = 0
)

func (handler *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	var err error
	// default
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	skipReview, err := queryChecker.CheckIsIn(w,r, "skip_reviews", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	limitReview, err := queryChecker.CheckIsIn(w, r, "limit_reviews", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		return
	}
	skipRecom, err := queryChecker.CheckIsIn(w,r, "skip_recommend", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	limitRecom, err := queryChecker.CheckIsIn(w, r, "limit_recommend", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		return
	}

	filmList, err := handler.FilmUsecase.GetFilm(uint64(id),skipReview,limitReview,skipRecom,limitRecom)
	if err == customErrors.ErrorSkip {
		http.Error(w, customErrors.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, customErrors.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(filmList)
	if err != nil {
		http.Error(w, customErrors.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *FilmHandler) PostRating (w http.ResponseWriter, r *http.Request) {
	authId, err := sessions.CheckSession(r)
	if err != nil {
		http.Error(w, customErrors.ErrDBMsg, http.StatusBadRequest)
		return
	}
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	rating, err := queryChecker.CheckIsInFloat64(w, r, "rating", 0, errors.New(customErrors.ErrEncMsg))
	if err != nil {
		return
	}

	newRating, err := handler.FilmUsecase.PostRating(id, authId, rating)
	if err != nil {
		http.Error(w, customErrors.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domain.NewRate{Rating: newRating})
	if err != nil {
		http.Error(w, customErrors.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *FilmHandler) LoadMyRv (w http.ResponseWriter, r *http.Request) {
	authId, err := sessions.CheckSession(r)
	if err != nil {
		http.Error(w, customErrors.ErrDBMsg, http.StatusBadRequest)
		return
	}
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		return
	}

	review, err := handler.FilmUsecase.LoadMyReview(id, authId)
	if err == customErrors.ErrorNoReviewForFilm {
		http.Error(w, customErrors.ErrDBMsg, http.StatusInternalServerError)
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

func (handler *FilmHandler) loadFilmReviews (w http.ResponseWriter, r *http.Request) {
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

	reviews, err := handler.FilmUsecase.LoadFilmReviews(id, skip, limit)
	if err == customErrors.ErrorSkip {
		http.Error(w, customErrors.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, customErrors.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		http.Error(w, customErrors.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *FilmHandler) loadFilmRecommendations (w http.ResponseWriter, r *http.Request) {
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
	recommendations, err := handler.FilmUsecase.LoadFilmRecommendations(id, skip, limit)
	if err == customErrors.ErrorSkip {
		http.Error(w, customErrors.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, customErrors.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(recommendations)
	if err != nil {
		http.Error(w, customErrors.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

