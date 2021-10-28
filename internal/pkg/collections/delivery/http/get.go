package http

import (
	"2021_2_MAMBa/internal/pkg/collections"
	"encoding/json"
	"net/http"
	"strconv"
)

func (handler *CollectionsHandler) GetCollections(w http.ResponseWriter, r *http.Request) {
	var err error
	// default
	limit, skip := 10, 0
	skipString, isIn := r.URL.Query()["skip"]
	if isIn {
		skip, err = strconv.Atoi(skipString[0])
		if err != nil || skip < 0 {
			http.Error(w, collections.ErrSkipMsg, http.StatusBadRequest)
			return
		}
	}
	limitString, isIn := r.URL.Query()["limit"]
	if isIn {
		limit, err = strconv.Atoi(limitString[0])
		if err != nil || limit <= 0 {
			http.Error(w, collections.ErrLimitMsg, http.StatusBadRequest)
			return
		}
	}

	collectionsList, err := handler.CollectionsUsecase.GetCollections(skip, limit)
	if err == collections.ErrorSkip {
		http.Error(w, collections.ErrSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, collections.ErrDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(collectionsList)
	if err != nil {
		http.Error(w, collections.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}
