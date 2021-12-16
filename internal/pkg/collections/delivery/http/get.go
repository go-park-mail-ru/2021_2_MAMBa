package http

import (
	"2021_2_MAMBa/internal/pkg/collections/delivery/grpc"
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"context"
	"encoding/json"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip  = 0
)

func (handler *CollectionsHandler) GetCollections(w http.ResponseWriter, r *http.Request) {
	skip, err := queryChecker.CheckIsIn(w, r, "skip", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	limit, err := queryChecker.CheckIsIn(w, r, "limit", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrLimitMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	collectionsList, err := handler.CollectionsClient.GetCollections(context.Background(), &grpc.SkipLimit{Skip: int64(skip), Limit: int64(limit)})
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
	x, err := json.Marshal(collectionsList)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *CollectionsHandler) GetCollectionFilms(w http.ResponseWriter, r *http.Request) {
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	collectionFilms, err := handler.CollectionsClient.GetCollectionPage(context.Background(), &grpc.ID{Id: id})
	if err != nil && err.Error() == customErrors.RPCErrNotFound {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrNotFoundMsg), Status: http.StatusNotFound}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}

	x, err := json.Marshal(collectionFilms)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
