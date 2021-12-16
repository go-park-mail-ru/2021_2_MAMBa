package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip  = 0
)

func (handler *SearchHandler) GetSearch(w http.ResponseWriter, r *http.Request) {
	// default
	queryList, isIn := r.URL.Query()["query"]
	query := queryList[0]
	if !isIn || query == "" {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSearchQuery), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	skipFilms, err := queryChecker.CheckIsIn(w, r, "skip_films", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	limitFilms, err := queryChecker.CheckIsIn(w, r, "limit_films", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrLimitMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	skipPersons, err := queryChecker.CheckIsIn(w, r, "skip_persons", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	limitPersons, err := queryChecker.CheckIsIn(w, r, "limit_persons", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrLimitMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	searchResult, err := handler.SearchUsecase.GetSearch(query, skipFilms, limitFilms, skipPersons, limitPersons)
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
	x, err := searchResult.MarshalJSON()
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
