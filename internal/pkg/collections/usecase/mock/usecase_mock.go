// Code generated by MockGen. DO NOT EDIT.
// Source: 2021_2_MAMBa/internal/pkg/domain (interfaces: CollectionsUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	domain "2021_2_MAMBa/internal/pkg/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCollectionsUsecase is a mock of CollectionsUsecase interface.
type MockCollectionsUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockCollectionsUsecaseMockRecorder
}

// MockCollectionsUsecaseMockRecorder is the mock recorder for MockCollectionsUsecase.
type MockCollectionsUsecaseMockRecorder struct {
	mock *MockCollectionsUsecase
}

// NewMockCollectionsUsecase creates a new mock instance.
func NewMockCollectionsUsecase(ctrl *gomock.Controller) *MockCollectionsUsecase {
	mock := &MockCollectionsUsecase{ctrl: ctrl}
	mock.recorder = &MockCollectionsUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCollectionsUsecase) EXPECT() *MockCollectionsUsecaseMockRecorder {
	return m.recorder
}

// GetCollectionPage mocks base method.
func (m *MockCollectionsUsecase) GetCollectionPage(arg0 uint64) (domain.CollectionPage, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollectionPage", arg0)
	ret0, _ := ret[0].(domain.CollectionPage)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollectionPage indicates an expected call of GetCollectionPage.
func (mr *MockCollectionsUsecaseMockRecorder) GetCollectionPage(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollectionPage", reflect.TypeOf((*MockCollectionsUsecase)(nil).GetCollectionPage), arg0)
}

// GetCollections mocks base method.
func (m *MockCollectionsUsecase) GetCollections(arg0, arg1 int) (domain.Collections, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCollections", arg0, arg1)
	ret0, _ := ret[0].(domain.Collections)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCollections indicates an expected call of GetCollections.
func (mr *MockCollectionsUsecaseMockRecorder) GetCollections(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCollections", reflect.TypeOf((*MockCollectionsUsecase)(nil).GetCollections), arg0, arg1)
}
