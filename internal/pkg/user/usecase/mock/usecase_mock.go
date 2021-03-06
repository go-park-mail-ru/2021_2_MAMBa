// Code generated by MockGen. DO NOT EDIT.
// Source: 2021_2_MAMBa/internal/pkg/domain (interfaces: UserUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	domain "2021_2_MAMBa/internal/pkg/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserUsecase is a mock of UserUsecase interface.
type MockUserUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUserUsecaseMockRecorder
}

// MockUserUsecaseMockRecorder is the mock recorder for MockUserUsecase.
type MockUserUsecaseMockRecorder struct {
	mock *MockUserUsecase
}

// NewMockUserUsecase creates a new mock instance.
func NewMockUserUsecase(ctrl *gomock.Controller) *MockUserUsecase {
	mock := &MockUserUsecase{ctrl: ctrl}
	mock.recorder = &MockUserUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserUsecase) EXPECT() *MockUserUsecaseMockRecorder {
	return m.recorder
}

// CheckAuth mocks base method.
func (m *MockUserUsecase) CheckAuth(arg0 uint64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckAuth", arg0)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckAuth indicates an expected call of CheckAuth.
func (mr *MockUserUsecaseMockRecorder) CheckAuth(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckAuth", reflect.TypeOf((*MockUserUsecase)(nil).CheckAuth), arg0)
}

// CreateSubscription mocks base method.
func (m *MockUserUsecase) CreateSubscription(arg0, arg1 uint64) (domain.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSubscription", arg0, arg1)
	ret0, _ := ret[0].(domain.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSubscription indicates an expected call of CreateSubscription.
func (mr *MockUserUsecaseMockRecorder) CreateSubscription(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubscription", reflect.TypeOf((*MockUserUsecase)(nil).CreateSubscription), arg0, arg1)
}

// GetBasicInfo mocks base method.
func (m *MockUserUsecase) GetBasicInfo(arg0 uint64) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBasicInfo", arg0)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBasicInfo indicates an expected call of GetBasicInfo.
func (mr *MockUserUsecaseMockRecorder) GetBasicInfo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBasicInfo", reflect.TypeOf((*MockUserUsecase)(nil).GetBasicInfo), arg0)
}

// GetProfileById mocks base method.
func (m *MockUserUsecase) GetProfileById(arg0, arg1 uint64) (domain.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProfileById", arg0, arg1)
	ret0, _ := ret[0].(domain.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProfileById indicates an expected call of GetProfileById.
func (mr *MockUserUsecaseMockRecorder) GetProfileById(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProfileById", reflect.TypeOf((*MockUserUsecase)(nil).GetProfileById), arg0, arg1)
}

// LoadUserReviews mocks base method.
func (m *MockUserUsecase) LoadUserReviews(arg0 uint64, arg1, arg2 int) (domain.FilmReviews, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadUserReviews", arg0, arg1, arg2)
	ret0, _ := ret[0].(domain.FilmReviews)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoadUserReviews indicates an expected call of LoadUserReviews.
func (mr *MockUserUsecaseMockRecorder) LoadUserReviews(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadUserReviews", reflect.TypeOf((*MockUserUsecase)(nil).LoadUserReviews), arg0, arg1, arg2)
}

// Login mocks base method.
func (m *MockUserUsecase) Login(arg0 *domain.UserToLogin) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserUsecaseMockRecorder) Login(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserUsecase)(nil).Login), arg0)
}

// Register mocks base method.
func (m *MockUserUsecase) Register(arg0 *domain.User) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockUserUsecaseMockRecorder) Register(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserUsecase)(nil).Register), arg0)
}

// UpdateAvatar mocks base method.
func (m *MockUserUsecase) UpdateAvatar(arg0 uint64, arg1 string) (domain.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAvatar", arg0, arg1)
	ret0, _ := ret[0].(domain.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAvatar indicates an expected call of UpdateAvatar.
func (mr *MockUserUsecaseMockRecorder) UpdateAvatar(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAvatar", reflect.TypeOf((*MockUserUsecase)(nil).UpdateAvatar), arg0, arg1)
}

// UpdateProfile mocks base method.
func (m *MockUserUsecase) UpdateProfile(arg0 domain.Profile) (domain.Profile, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProfile", arg0)
	ret0, _ := ret[0].(domain.Profile)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateProfile indicates an expected call of UpdateProfile.
func (mr *MockUserUsecaseMockRecorder) UpdateProfile(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProfile", reflect.TypeOf((*MockUserUsecase)(nil).UpdateProfile), arg0)
}
