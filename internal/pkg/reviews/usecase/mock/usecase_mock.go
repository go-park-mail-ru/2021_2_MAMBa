// Code generated by MockGen. DO NOT EDIT.
// Source: 2021_2_MAMBa/internal/pkg/domain (interfaces: ReviewUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	domain "2021_2_MAMBa/internal/pkg/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockReviewUsecase is a mock of ReviewUsecase interface.
type MockReviewUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockReviewUsecaseMockRecorder
}

// MockReviewUsecaseMockRecorder is the mock recorder for MockReviewUsecase.
type MockReviewUsecaseMockRecorder struct {
	mock *MockReviewUsecase
}

// NewMockReviewUsecase creates a new mock instance.
func NewMockReviewUsecase(ctrl *gomock.Controller) *MockReviewUsecase {
	mock := &MockReviewUsecase{ctrl: ctrl}
	mock.recorder = &MockReviewUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReviewUsecase) EXPECT() *MockReviewUsecaseMockRecorder {
	return m.recorder
}

// GetReview mocks base method.
func (m *MockReviewUsecase) GetReview(arg0 uint64) (domain.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReview", arg0)
	ret0, _ := ret[0].(domain.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReview indicates an expected call of GetReview.
func (mr *MockReviewUsecaseMockRecorder) GetReview(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReview", reflect.TypeOf((*MockReviewUsecase)(nil).GetReview), arg0)
}

// LoadReviewsExcept mocks base method.
func (m *MockReviewUsecase) LoadReviewsExcept(arg0, arg1 uint64, arg2, arg3 int) (domain.FilmReviews, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadReviewsExcept", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(domain.FilmReviews)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadReviewsExcept indicates an expected call of LoadReviewsExcept.
func (mr *MockReviewUsecaseMockRecorder) LoadReviewsExcept(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadReviewsExcept", reflect.TypeOf((*MockReviewUsecase)(nil).LoadReviewsExcept), arg0, arg1, arg2, arg3)
}

// PostReview mocks base method.
func (m *MockReviewUsecase) PostReview(arg0 domain.Review) (domain.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PostReview", arg0)
	ret0, _ := ret[0].(domain.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PostReview indicates an expected call of PostReview.
func (mr *MockReviewUsecaseMockRecorder) PostReview(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PostReview", reflect.TypeOf((*MockReviewUsecase)(nil).PostReview), arg0)
}
