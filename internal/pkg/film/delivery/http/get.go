package http

import (
	"2021_2_MAMBa/internal/pkg/film"
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
	skipString, isIn = r.URL.Query()["skip_recommnd"]
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