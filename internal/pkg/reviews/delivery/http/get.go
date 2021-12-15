package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"2021_2_MAMBa/internal/pkg/utils/xss"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip  = 0
)

func (handler *ReviewHandler) GetReview(w http.ResponseWriter, r *http.Request) {
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	review, err := handler.ReiviewUsecase.GetReview(id)
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
	x, err := review.MarshalJSON()
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *ReviewHandler) PostReview(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	reviewForm := domain.Review{}
	var p []byte
	r.Body.Read(p)
	err := reviewForm.UnmarshalJSON(p)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorBadInput.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	xss.SanitizeReview(&reviewForm)

	rq := cast.CookieToRq(r, 0)
	authIDMessage, err := handler.AuthClient.CheckSession(r.Context(), &rq)
	if err != nil && err.Error() == customErrors.RPCErrUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}
	authId := authIDMessage.ID

	reviewForm.AuthorId = authId
	review, err := handler.ReiviewUsecase.PostReview(reviewForm)
	x, err := review.MarshalJSON()
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *ReviewHandler) LoadExcept(w http.ResponseWriter, r *http.Request) {
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	filmId, err := queryChecker.CheckIsIn64(w, r, "film_id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
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
	review, err := handler.ReiviewUsecase.LoadReviewsExcept(id, filmId, skip, limit)
	if err == customErrors.ErrorSkip {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := review.MarshalJSON()
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
