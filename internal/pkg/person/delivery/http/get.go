package http

import (
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"encoding/json"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip = 0
)

func (handler *PersonHandler) GetPerson (w http.ResponseWriter, r *http.Request) {
	var err error
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	page, err := handler.PersonUsecase.GetPerson(id)
	if err == customErrors.ErrorBadInput {
		http.Error(w, customErrors.ErrorBadInput.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(page)
	if err != nil {
		http.Error(w, customErrors.ErrorInternalServer.Error(), http.StatusInternalServerError)
		return
	}
}

func (handler *PersonHandler) GetPersonFilms (w http.ResponseWriter, r *http.Request) {
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
	films, err := handler.PersonUsecase.GetFilms(id,skip, limit)
	if err == customErrors.ErrorSkip {
		http.Error(w, customErrors.ErrorSkip.Error(), http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, customErrors.ErrSkipMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(films)
	if err != nil {
		http.Error(w, customErrors.ErrSkipMsg, http.StatusInternalServerError)
		return
	}
}