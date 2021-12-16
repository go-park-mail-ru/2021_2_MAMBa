// Code generated by MockGen. DO NOT EDIT.
// Source: 2021_2_MAMBa/internal/pkg/domain (interfaces: FilmRepository)

// Package mock is a generated GoMock package.
package mock

import (
	domain "2021_2_MAMBa/internal/pkg/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockFilmRepository is a mock of FilmRepository interface.
type MockFilmRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFilmRepositoryMockRecorder
}

// MockFilmRepositoryMockRecorder is the mock recorder for MockFilmRepository.
type MockFilmRepositoryMockRecorder struct {
	mock *MockFilmRepository
}

// NewMockFilmRepository creates a new mock instance.
func NewMockFilmRepository(ctrl *gomock.Controller) *MockFilmRepository {
	mock := &MockFilmRepository{ctrl: ctrl}
	mock.recorder = &MockFilmRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilmRepository) EXPECT() *MockFilmRepositoryMockRecorder {
	return m.recorder
}

// BookmarkFilm mocks base method.
func (m *MockFilmRepository) BookmarkFilm(arg0, arg1 uint64, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BookmarkFilm", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// BookmarkFilm indicates an expected call of BookmarkFilm.
func (mr *MockFilmRepositoryMockRecorder) BookmarkFilm(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BookmarkFilm", reflect.TypeOf((*MockFilmRepository)(nil).BookmarkFilm), arg0, arg1, arg2)
}

// CheckFilmBookmarked mocks base method.
func (m *MockFilmRepository) CheckFilmBookmarked(arg0, arg1 uint64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckFilmBookmarked", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckFilmBookmarked indicates an expected call of CheckFilmBookmarked.
func (mr *MockFilmRepositoryMockRecorder) CheckFilmBookmarked(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckFilmBookmarked", reflect.TypeOf((*MockFilmRepository)(nil).CheckFilmBookmarked), arg0, arg1)
}

// CountBookmarks mocks base method.
func (m *MockFilmRepository) CountBookmarks(arg0 uint64) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountBookmarks", arg0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountBookmarks indicates an expected call of CountBookmarks.
func (mr *MockFilmRepositoryMockRecorder) CountBookmarks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountBookmarks", reflect.TypeOf((*MockFilmRepository)(nil).CountBookmarks), arg0)
}

// GetBanners mocks base method.
func (m *MockFilmRepository) GetBanners() (domain.BannersList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBanners")
	ret0, _ := ret[0].(domain.BannersList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBanners indicates an expected call of GetBanners.
func (mr *MockFilmRepositoryMockRecorder) GetBanners() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBanners", reflect.TypeOf((*MockFilmRepository)(nil).GetBanners))
}

// GetFilm mocks base method.
func (m *MockFilmRepository) GetFilm(arg0 uint64) (domain.Film, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilm", arg0)
	ret0, _ := ret[0].(domain.Film)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilm indicates an expected call of GetFilm.
func (mr *MockFilmRepositoryMockRecorder) GetFilm(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilm", reflect.TypeOf((*MockFilmRepository)(nil).GetFilm), arg0)
}

// GetFilmRecommendations mocks base method.
func (m *MockFilmRepository) GetFilmRecommendations(arg0 uint64, arg1, arg2 int) (domain.FilmRecommendations, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmRecommendations", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.FilmRecommendations)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmRecommendations indicates an expected call of GetFilmRecommendations.
func (mr *MockFilmRepositoryMockRecorder) GetFilmRecommendations(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmRecommendations", reflect.TypeOf((*MockFilmRepository)(nil).GetFilmRecommendations), arg0, arg1, arg2)
}

// GetFilmReviews mocks base method.
func (m *MockFilmRepository) GetFilmReviews(arg0 uint64, arg1, arg2 int) (domain.FilmReviews, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmReviews", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.FilmReviews)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmReviews indicates an expected call of GetFilmReviews.
func (mr *MockFilmRepositoryMockRecorder) GetFilmReviews(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmReviews", reflect.TypeOf((*MockFilmRepository)(nil).GetFilmReviews), arg0, arg1, arg2)
}

// GetFilmsByGenre mocks base method.
func (m *MockFilmRepository) GetFilmsByGenre(arg0 uint64, arg1, arg2 int) (domain.GenreFilmList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmsByGenre", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.GenreFilmList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmsByGenre indicates an expected call of GetFilmsByGenre.
func (mr *MockFilmRepositoryMockRecorder) GetFilmsByGenre(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmsByGenre", reflect.TypeOf((*MockFilmRepository)(nil).GetFilmsByGenre), arg0, arg1, arg2)
}

// GetFilmsByMonthYear mocks base method.
func (m *MockFilmRepository) GetFilmsByMonthYear(arg0, arg1, arg2, arg3 int) (domain.FilmList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilmsByMonthYear", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(domain.FilmList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilmsByMonthYear indicates an expected call of GetFilmsByMonthYear.
func (mr *MockFilmRepositoryMockRecorder) GetFilmsByMonthYear(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilmsByMonthYear", reflect.TypeOf((*MockFilmRepository)(nil).GetFilmsByMonthYear), arg0, arg1, arg2, arg3)
}

// GetGenres mocks base method.
func (m *MockFilmRepository) GetGenres() (domain.GenresList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGenres")
	ret0, _ := ret[0].(domain.GenresList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGenres indicates an expected call of GetGenres.
func (mr *MockFilmRepositoryMockRecorder) GetGenres() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGenres", reflect.TypeOf((*MockFilmRepository)(nil).GetGenres))
}

// GetMyReview mocks base method.
func (m *MockFilmRepository) GetMyReview(arg0, arg1 uint64) (domain.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMyReview", arg0, arg1)
	ret0, _ := ret[0].(domain.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyReview indicates an expected call of GetMyReview.
func (mr *MockFilmRepositoryMockRecorder) GetMyReview(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyReview", reflect.TypeOf((*MockFilmRepository)(nil).GetMyReview), arg0, arg1)
}

// GetPopularFilms mocks base method.
func (m *MockFilmRepository) GetPopularFilms() (domain.FilmList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPopularFilms")
	ret0, _ := ret[0].(domain.FilmList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPopularFilms indicates an expected call of GetPopularFilms.
func (mr *MockFilmRepositoryMockRecorder) GetPopularFilms() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPopularFilms", reflect.TypeOf((*MockFilmRepository)(nil).GetPopularFilms))
}

// LoadUserBookmarkedFilmsID mocks base method.
func (m *MockFilmRepository) LoadUserBookmarkedFilmsID(arg0 uint64, arg1, arg2 int) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadUserBookmarkedFilmsID", arg0, arg1, arg2)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadUserBookmarkedFilmsID indicates an expected call of LoadUserBookmarkedFilmsID.
func (mr *MockFilmRepositoryMockRecorder) LoadUserBookmarkedFilmsID(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadUserBookmarkedFilmsID", reflect.TypeOf((*MockFilmRepository)(nil).LoadUserBookmarkedFilmsID), arg0, arg1, arg2)
}

// PostRating mocks base method.
func (m *MockFilmRepository) PostRating(arg0, arg1 uint64, arg2 float64) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostRating", arg0, arg1, arg2)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostRating indicates an expected call of PostRating.
func (mr *MockFilmRepositoryMockRecorder) PostRating(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostRating", reflect.TypeOf((*MockFilmRepository)(nil).PostRating), arg0, arg1, arg2)
}