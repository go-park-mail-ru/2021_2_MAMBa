package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/film"
	"2021_2_MAMBa/internal/pkg/reviews"
	"2021_2_MAMBa/internal/pkg/sessions"
	"encoding/json"
	"net/http"
	"strconv"
)

func (handler *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	var id uint64
	var err error
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, reviews.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	review, err := handler.ReiviewUsecase.GetReview(id)
	if err == reviews.ErrorSkip {
		http.Error(w, reviews.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, reviews.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(review)
	if err != nil {
		http.Error(w, reviews.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *ReviewHandler) PostReview(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reviewForm := domain.Review{}
	err := json.NewDecoder(r.Body).Decode(&reviewForm)
	if err != nil {
		http.Error(w, reviews.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	authId, err := sessions.CheckSession(r)
	if err != nil {
		http.Error(w, reviews.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	reviewForm.AuthorId = authId
	review, err := handler.ReiviewUsecase.PostReview(reviewForm)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(review)
	if err != nil {
		http.Error(w, reviews.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *ReviewHandler) LoadExcept (w http.ResponseWriter, r *http.Request) {
	var id, filmId uint64
	skip, limit := 0, 10
	var err error
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, reviews.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	idString, isIn = r.URL.Query()["film_id"]
	if isIn {
		filmId, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, reviews.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	skipString, isIn := r.URL.Query()["skip"]
	if isIn {
		skip, err = strconv.Atoi(skipString[0])
		if err != nil || skip  < 0 {
			http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	limitString, isIn := r.URL.Query()["limit"]
	if isIn {
		limit, err = strconv.Atoi(limitString[0])
		if err != nil || limit <= 0 {
			http.Error(w, film.ErrLimitMsg, http.StatusBadRequest)
			return
		}
	}
	review, err := handler.ReiviewUsecase.LoadReviewsExcept(id, filmId, skip, limit)
	if err == reviews.ErrorSkip {
		http.Error(w, reviews.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, reviews.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(review)
	if err != nil {
		http.Error(w, reviews.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}