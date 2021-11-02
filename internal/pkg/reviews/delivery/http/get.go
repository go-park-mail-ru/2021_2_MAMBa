package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/sessions"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"encoding/json"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip = 0
)

func (handler *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	var err error
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	review, err := handler.ReiviewUsecase.GetReview(id)
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

func (handler *ReviewHandler) PostReview(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reviewForm := domain.Review{}
	err := json.NewDecoder(r.Body).Decode(&reviewForm)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	authId, err := sessions.CheckSession(r)
	if err != nil {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	reviewForm.AuthorId = authId
	review, err := handler.ReiviewUsecase.PostReview(reviewForm)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(review)
	if err != nil {
		http.Error(w, customErrors.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *ReviewHandler) LoadExcept (w http.ResponseWriter, r *http.Request) {
	var err error
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	filmId, err := queryChecker.CheckIsIn64(w, r, "film_id", 0, customErrors.ErrorSkip)
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
	review, err := handler.ReiviewUsecase.LoadReviewsExcept(id, filmId, skip, limit)
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