package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/film"
	"2021_2_MAMBa/internal/pkg/sessions"
	"encoding/json"
	"net/http"
	"strconv"
)

func (handler *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	var err error
	// default
	id, limitReview, skipReview, limitRecom, skipRecom := 0, 10, 0, 10, 0
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.Atoi(idString[0])
		if err != nil || id < 0 {
			http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	skipString, isIn := r.URL.Query()["skip_reviews"]
	if isIn {
		skipReview, err = strconv.Atoi(skipString[0])
		if err != nil || skipReview < 0 {
			http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	limitString, isIn := r.URL.Query()["limit_reviews"]
	if isIn {
		limitReview, err = strconv.Atoi(limitString[0])
		if err != nil || limitReview <= 0 {
			http.Error(w, film.ErrLimitMsg, http.StatusBadRequest)
			return
		}
	}
	skipString, isIn = r.URL.Query()["skip_recommend"]
	if isIn {
		skipRecom, err = strconv.Atoi(skipString[0])
		if err != nil || skipRecom < 0 {
			http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	limitString, isIn = r.URL.Query()["limit_recommend"]
	if isIn {
		limitRecom, err = strconv.Atoi(limitString[0])
		if err != nil || limitRecom <= 0 {
			http.Error(w, film.ErrLimitMsg, http.StatusBadRequest)
			return
		}
	}

	filmList, err := handler.FilmUsecase.GetFilm(uint64(id),skipReview,limitReview,skipRecom,limitRecom)
	if err == film.ErrorSkip {
		http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, film.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(filmList)
	if err != nil {
		http.Error(w, film.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *FilmHandler) PostRating (w http.ResponseWriter, r *http.Request) {
	var id uint64
	var rating float64
	authId, err := sessions.CheckSession(r)
	if err != nil {
		http.Error(w, film.ErrDBMsg, http.StatusBadRequest)
		return
	}
	rating = -1
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	ratingString, isIn := r.URL.Query()["rating"]
	if isIn {
		rating, err = strconv.ParseFloat(ratingString[0], 64)
		if err != nil || rating <= 0 {
			http.Error(w, film.ErrEncMsg, http.StatusBadRequest)
			return
		}
	}

	newRating, err := handler.FilmUsecase.PostRating(id, authId, rating)
	if err != nil {
		http.Error(w, film.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(domain.NewRate{Rating: newRating})
	if err != nil {
		http.Error(w, film.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *FilmHandler) LoadMyRv (w http.ResponseWriter, r *http.Request) {
	var id uint64
	authId, err := sessions.CheckSession(r)
	if err != nil {
		http.Error(w, film.ErrDBMsg, http.StatusBadRequest)
		return
	}
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id <= 0 {
			http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}

	review, err := handler.FilmUsecase.LoadMyReview(id, authId)
	if err == film.ErrorNoReviewForFilm {
		http.Error(w, film.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	if err != nil {
		http.Error(w, film.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(review)
	if err != nil {
		http.Error(w, film.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *FilmHandler) loadFilmReviews (w http.ResponseWriter, r *http.Request) {
	limit, skip := 10, 0
	var err error
	var id uint64
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id < 0 {
			http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
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
	reviews, err := handler.FilmUsecase.LoadFilmReviews(id, skip, limit)
	if err == film.ErrorSkip {
		http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, film.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(reviews)
	if err != nil {
		http.Error(w, film.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

func (handler *FilmHandler) loadFilmRecommendations (w http.ResponseWriter, r *http.Request) {
	limit, skip := 10, 0
	var err error
	var id uint64
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id < 0 {
			http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
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
	recommendations, err := handler.FilmUsecase.LoadFilmRecommendations(id, skip, limit)
	if err == film.ErrorSkip {
		http.Error(w, film.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, film.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(recommendations)
	if err != nil {
		http.Error(w, film.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

