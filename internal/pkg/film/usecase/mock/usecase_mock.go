// Code generated by MockGen. DO NOT EDIT.
// Source: 2021_2_MAMBa/internal/pkg/domain (interfaces: FilmUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	domain "2021_2_MAMBa/internal/pkg/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFilmUsecase is a mock of FilmUsecase interface.
type MockFilmUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockFilmUsecaseMockRecorder
}

// MockFilmUsecaseMockRecorder is the mock recorder for MockFilmUsecase.
type MockFilmUsecaseMockRecorder struct {
	mock *MockFilmUsecase
}

// NewMockFilmUsecase creates a new mock instance.
func NewMockFilmUsecase(ctrl *gomock.Controller) *MockFilmUsecase {
	mock := &MockFilmUsecase{ctrl: ctrl}
	mock.recorder = &MockFilmUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilmUsecase) EXPECT() *MockFilmUsecaseMockRecorder {
	return m.recorder
}

// BookmarkFilm mocks base method.
func (m *MockFilmUsecase) BookmarkFilm(arg0, arg1 uint64, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BookmarkFilm", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// BookmarkFilm indicates an expected call of BookmarkFilm.
func (mr *MockFilmUsecaseMockRecorder) BookmarkFilm(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BookmarkFilm", reflect.TypeOf((*MockFilmUsecase)(nil).BookmarkFilm), arg0, arg1, arg2)
}

// GetBanners mocks base method.
func (m *MockFilmUsecase) GetBanners() (domain.BannersList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBanners")
	ret0, _ := ret[0].(domain.BannersList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBanners indicates an expected call of GetBanners.
func (mr *MockFilmUsecaseMockRecorder) GetBanners() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBanners", reflect.TypeOf((*MockFilmUsecase)(nil).GetBanners))
}

// GetFilm mocks base method.
func (m *MockFilmUsecase) GetFilm(arg0, arg1 uint64, arg2, arg3, arg4, arg5 int) (domain.FilmPageInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilm", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(domain.FilmPageInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilm indicates an expected call of GetFilm.
func (mr *MockFilmUsecaseMockRecorder) GetFilm(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilm", reflect.TypeOf((*MockFilmUsecase)(nil).GetFilm), arg0, arg1, arg2, arg3, arg4, arg5)
}

// GetFilmsByGenre mocks base method.
func (m *MockFilmUsecase) GetFilmsByGenre(arg0 uint64, arg1, arg2 int) (domain.GenreFilmList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmsByGenre", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.GenreFilmList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmsByGenre indicates an expected call of GetFilmsByGenre.
func (mr *MockFilmUsecaseMockRecorder) GetFilmsByGenre(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmsByGenre", reflect.TypeOf((*MockFilmUsecase)(nil).GetFilmsByGenre), arg0, arg1, arg2)
}

// GetFilmsByMonthYear mocks base method.
func (m *MockFilmUsecase) GetFilmsByMonthYear(arg0, arg1, arg2, arg3 int) (domain.FilmList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmsByMonthYear", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(domain.FilmList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmsByMonthYear indicates an expected call of GetFilmsByMonthYear.
func (mr *MockFilmUsecaseMockRecorder) GetFilmsByMonthYear(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmsByMonthYear", reflect.TypeOf((*MockFilmUsecase)(nil).GetFilmsByMonthYear), arg0, arg1, arg2, arg3)
}

// GetGenres mocks base method.
func (m *MockFilmUsecase) GetGenres() (domain.GenresList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres")
	ret0, _ := ret[0].(domain.GenresList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockFilmUsecaseMockRecorder) GetGenres() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockFilmUsecase)(nil).GetGenres))
}

// GetPopularFilms mocks base method.
func (m *MockFilmUsecase) GetPopularFilms() (domain.FilmList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPopularFilms")
	ret0, _ := ret[0].(domain.FilmList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPopularFilms indicates an expected call of GetPopularFilms.
func (mr *MockFilmUsecaseMockRecorder) GetPopularFilms() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPopularFilms", reflect.TypeOf((*MockFilmUsecase)(nil).GetPopularFilms))
}

// LoadFilmRecommendations mocks base method.
func (m *MockFilmUsecase) LoadFilmRecommendations(arg0 uint64, arg1, arg2 int) (domain.FilmRecommendations, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadFilmRecommendations", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.FilmRecommendations)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadFilmRecommendations indicates an expected call of LoadFilmRecommendations.
func (mr *MockFilmUsecaseMockRecorder) LoadFilmRecommendations(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadFilmRecommendations", reflect.TypeOf((*MockFilmUsecase)(nil).LoadFilmRecommendations), arg0, arg1, arg2)
}

// LoadFilmReviews mocks base method.
func (m *MockFilmUsecase) LoadFilmReviews(arg0 uint64, arg1, arg2 int) (domain.FilmReviews, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadFilmReviews", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.FilmReviews)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadFilmReviews indicates an expected call of LoadFilmReviews.
func (mr *MockFilmUsecaseMockRecorder) LoadFilmReviews(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadFilmReviews", reflect.TypeOf((*MockFilmUsecase)(nil).LoadFilmReviews), arg0, arg1, arg2)
}

// LoadMyReview mocks base method.
func (m *MockFilmUsecase) LoadMyReview(arg0, arg1 uint64) (domain.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadMyReview", arg0, arg1)
	ret0, _ := ret[0].(domain.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadMyReview indicates an expected call of LoadMyReview.
func (mr *MockFilmUsecaseMockRecorder) LoadMyReview(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadMyReview", reflect.TypeOf((*MockFilmUsecase)(nil).LoadMyReview), arg0, arg1)
}

// LoadUserBookmarks mocks base method.
func (m *MockFilmUsecase) LoadUserBookmarks(arg0 uint64, arg1, arg2 int) (domain.FilmBookmarks, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadUserBookmarks", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.FilmBookmarks)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadUserBookmarks indicates an expected call of LoadUserBookmarks.
func (mr *MockFilmUsecaseMockRecorder) LoadUserBookmarks(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadUserBookmarks", reflect.TypeOf((*MockFilmUsecase)(nil).LoadUserBookmarks), arg0, arg1, arg2)
}

// PostRating mocks base method.
func (m *MockFilmUsecase) PostRating(arg0, arg1 uint64, arg2 float64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostRating", arg0, arg1, arg2)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostRating indicates an expected call of PostRating.
func (mr *MockFilmUsecaseMockRecorder) PostRating(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostRating", reflect.TypeOf((*MockFilmUsecase)(nil).PostRating), arg0, arg1, arg2)
}
