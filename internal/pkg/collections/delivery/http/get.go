package http

import (
	"2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"encoding/json"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip  = 0
)

func (handler *CollectionsHandler) GetCollections(w http.ResponseWriter, r *http.Request) {
	var err error
	skip, err := queryChecker.CheckIsIn(w, r, "skip", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		return
	}
	limit, err := queryChecker.CheckIsIn(w, r, "limit", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		return
	}

	collectionsList, err := handler.CollectionsUsecase.GetCollections(skip, limit)
	if err == customErrors.ErrorSkip {
		http.Error(w, customErrors.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, customErrors.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(collectionsList)
	if err != nil {
		http.Error(w, customErrors.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}
