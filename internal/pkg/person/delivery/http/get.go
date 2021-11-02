package http

import (
	"2021_2_MAMBa/internal/pkg/film"
	"2021_2_MAMBa/internal/pkg/person"
	"encoding/json"
	"net/http"
	"strconv"
)

func (handler *PersonHandler) GetPerson (w http.ResponseWriter, r *http.Request) {
	var err error
	var id uint64
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id < 0 {
			http.Error(w, person.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, person.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	page, err := handler.PersonUsecase.GetPerson(id)
	if err != nil {
		http.Error(w, person.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(page)
	if err != nil {
		http.Error(w, person.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *PersonHandler) GetPersonFilms (w http.ResponseWriter, r *http.Request) {
	var err error
	limit, skip := 10, 0
	var id uint64
	idString, isIn := r.URL.Query()["id"]
	if isIn {
		id, err = strconv.ParseUint(idString[0], 10, 64)
		if err != nil || id < 0 {
			http.Error(w, person.ErrSkipMsg, http.StatusBadRequest)
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
	films, err := handler.PersonUsecase.GetFilms(id,skip, limit)
	if err == person.ErrorSkip {
		http.Error(w, person.ErrorSkip.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, person.ErrSkipMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(films)
	if err != nil {
		http.Error(w, person.ErrSkipMsg, http.StatusInternalServerError)
		return
	}
}