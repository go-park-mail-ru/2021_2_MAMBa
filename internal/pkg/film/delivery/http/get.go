package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"2021_2_MAMBa/internal/pkg/utils/queryChecker"
	"encoding/json"
	"errors"
	"net/http"
)

const (
	defaultLimit = 10
	defaultSkip  = 0
	minYear      = 1800
	maxYear      = 2100
	minMonth     = 0
	maxMonth     = 13
)

func (handler *FilmHandler) GetFilm(w http.ResponseWriter, r *http.Request) {
	filmID, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	skipReview, err := queryChecker.CheckIsIn(w, r, "skip_reviews", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	limitReview, err := queryChecker.CheckIsIn(w, r, "limit_reviews", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrLimitMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	skipRecom, err := queryChecker.CheckIsIn(w, r, "skip_recommend", defaultSkip, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	limitRecom, err := queryChecker.CheckIsIn(w, r, "limit_recommend", defaultLimit, customErrors.ErrorLimit)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrLimitMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	rq := cast.CookieToRq(r, 0)
	authIDMessage, err := handler.AuthClient.CheckSession(r.Context(), &rq)
	var authId uint64 = 0
	if err != nil && err.Error() != customErrors.RPCErrUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(err.Error()), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	} else if err == nil {
		authId = authIDMessage.ID
	}

	filmPageInfo, err := handler.FilmUsecase.GetFilm(authId, filmID, skipReview, limitReview, skipRecom, limitRecom)
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

	if filmPageInfo.FilmMain.Id == 0 {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrNotFoundMsg), Status: http.StatusNotFound}
		resp.Write(w)
		return
	}

	x, err := filmPageInfo.CustomEasyJSON()
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) PostRating(w http.ResponseWriter, r *http.Request) {
	rq := cast.CookieToRq(r, 0)
	authIDMessage, err := handler.AuthClient.CheckSession(r.Context(), &rq)
	if err != nil && err.Error() == customErrors.RPCErrUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	authId := authIDMessage.ID
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	rating, err := queryChecker.CheckIsInFloat64(w, r, "rating", 0, errors.New(customErrors.ErrEncMsg))
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrRateMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	newRating, err := handler.FilmUsecase.PostRating(id, authId, rating)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	jsRate := cast.Float64toJSONp1f(newRating)
	x, err := json.Marshal(domain.NewRate{Rating: jsRate})
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) LoadMyRv(w http.ResponseWriter, r *http.Request) {
	rq := cast.CookieToRq(r, 0)
	authIDMessage, err := handler.AuthClient.CheckSession(r.Context(), &rq)
	if err != nil && err.Error() == customErrors.RPCErrUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	authId := authIDMessage.ID
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	review, err := handler.FilmUsecase.LoadMyReview(id, authId)
	if err == customErrors.ErrorNoReviewForFilm {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrNoReviewMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(review)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) loadFilmReviews(w http.ResponseWriter, r *http.Request) {
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
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

	reviews, err := handler.FilmUsecase.LoadFilmReviews(id, skip, limit)
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
	x, err := json.Marshal(reviews)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) GetFilmsByMonthYear(w http.ResponseWriter, r *http.Request) {
	month, err := queryChecker.CheckIsIn(w, r, "month", 0, customErrors.ErrorDate)
	if err != nil || !(month < maxMonth && month > minMonth) {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDateMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	year, err := queryChecker.CheckIsIn(w, r, "year", 0, customErrors.ErrorDate)
	if err != nil || !(year < maxYear && year > minYear) {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDateMsg), Status: http.StatusBadRequest}
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

	filmList, err := handler.FilmUsecase.GetFilmsByMonthYear(month, year, limit, skip)
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
	x, err := json.Marshal(filmList)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	if filmList.FilmTotal == 0 {
		resp.Status = http.StatusNotFound
	}
	resp.Write(w)
}

func (handler *FilmHandler) loadFilmRecommendations(w http.ResponseWriter, r *http.Request) {
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
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
	recommendations, err := handler.FilmUsecase.LoadFilmRecommendations(id, skip, limit)
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
	x, err := json.Marshal(recommendations)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) LoadUserBookmarks(w http.ResponseWriter, r *http.Request) {
	id, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
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

	bookmarks, err := handler.FilmUsecase.LoadUserBookmarks(id, skip, limit)
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
	x, err := json.Marshal(bookmarks)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) BookmarkFilm(w http.ResponseWriter, r *http.Request) {
	rq := cast.CookieToRq(r, 0)
	authIDMessage, err := handler.AuthClient.CheckSession(r.Context(), &rq)
	if err != nil && err.Error() == customErrors.RPCErrUserNotLoggedIn {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrorUserNotLoggedIn.Error()), Status: http.StatusUnauthorized}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	userID := authIDMessage.ID
	filmID, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorID)
	if err != nil || filmID == 0 {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	bookmarked, err := queryChecker.CheckIsInBool(w, r, "bookmarked", false, customErrors.ErrorBookmarked)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrBookmarkedMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	err = handler.FilmUsecase.BookmarkFilm(userID, filmID, bookmarked)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}

	bookmarkedResult := domain.PostBookmarkResult{
		FilmID:     filmID,
		Bookmarked: bookmarked,
	}
	x, err := json.Marshal(bookmarkedResult)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genreList, err := handler.FilmUsecase.GetGenres()
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(genreList)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) GetBanners(w http.ResponseWriter, r *http.Request) {
	bannersList, err := handler.FilmUsecase.GetBanners()
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(bannersList)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) GetPopularFilms(w http.ResponseWriter, r *http.Request) {
	filmsList, err := handler.FilmUsecase.GetPopularFilms()
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(filmsList)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) GetFilmsByGenre(w http.ResponseWriter, r *http.Request) {
	genreID, err := queryChecker.CheckIsIn64(w, r, "id", 0, customErrors.ErrorSkip)
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

	genreFilmList, err := handler.FilmUsecase.GetFilmsByGenre(genreID, limit, skip)
	if err == customErrors.ErrorSkip {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrSkipMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	if err == customErrors.ErrNotFound {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrNotFoundMsg), Status: http.StatusNotFound}
		resp.Write(w)
		return
	}
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(genreFilmList)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}

func (handler *FilmHandler) GetRandomFilms(w http.ResponseWriter, r *http.Request) {
	genreID1, err := queryChecker.CheckIsIn64(w, r, "id1", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	genreID2, err := queryChecker.CheckIsIn64(w, r, "id2", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	genreID3, err := queryChecker.CheckIsIn64(w, r, "id3", 0, customErrors.ErrorSkip)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrIdMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	yearStart, err := queryChecker.CheckIsIn(w, r, "year_start", 0, customErrors.ErrorDate)
	if err != nil || !(yearStart < maxYear && yearStart > minYear) {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDateMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}
	yearEnd, err := queryChecker.CheckIsIn(w, r, "year_end", 0, customErrors.ErrorDate)
	if err != nil || !(yearEnd < maxYear && yearEnd > minYear) {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDateMsg), Status: http.StatusBadRequest}
		resp.Write(w)
		return
	}

	list, err := handler.FilmUsecase.GetRandomFilms(genreID1, genreID2, genreID3, yearStart, yearEnd)
	if err != nil {
		resp := domain.Response{Body: cast.ErrorToJson(customErrors.ErrDBMsg), Status: http.StatusInternalServerError}
		resp.Write(w)
		return
	}
	x, err := json.Marshal(list)
	resp := domain.Response{
		Body:   x,
		Status: http.StatusOK,
	}
	resp.Write(w)
}
