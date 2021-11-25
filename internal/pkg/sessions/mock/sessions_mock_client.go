// Code generated by MockGen. DO NOT EDIT.
// Source: 2021_2_MAMBa/internal/pkg/sessions/delivery/grpc (interfaces: SessionRPCClient)

// Package sessions_mock is a generated GoMock package.
package sessions_mock

import (
	grpc "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	grpc0 "google.golang.org/grpc"
)

// MockSessionRPCClient is a mock of SessionRPCClient interface.
type MockSessionRPCClient struct {
	ctrl     *gomock.Controller
	recorder *MockSessionRPCClientMockRecorder
}

// MockSessionRPCClientMockRecorder is the mock recorder for MockSessionRPCClient.
type MockSessionRPCClientMockRecorder struct {
	mock *MockSessionRPCClient
}

// NewMockSessionRPCClient creates a new mock instance.
func NewMockSessionRPCClient(ctrl *gomock.Controller) *MockSessionRPCClient {
	mock := &MockSessionRPCClient{ctrl: ctrl}
	mock.recorder = &MockSessionRPCClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionRPCClient) EXPECT() *MockSessionRPCClientMockRecorder {
	return m.recorder
}

// CheckSession mocks base method.
func (m *MockSessionRPCClient) CheckSession(arg0 context.Context, arg1 *grpc.Request, arg2 ...grpc0.CallOption) (*grpc.ID, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "CheckSession", varargs...)
	ret0, _ := ret[0].(*grpc.ID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckSession indicates an expected call of CheckSession.
func (mr *MockSessionRPCClientMockRecorder) CheckSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckSession", reflect.TypeOf((*MockSessionRPCClient)(nil).CheckSession), varargs...)
}

// EndSession mocks base method.
func (m *MockSessionRPCClient) EndSession(arg0 context.Context, arg1 *grpc.Request, arg2 ...grpc0.CallOption) (*grpc.Session, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "EndSession", varargs...)
	ret0, _ := ret[0].(*grpc.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EndSession indicates an expected call of EndSession.
func (mr *MockSessionRPCClientMockRecorder) EndSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndSession", reflect.TypeOf((*MockSessionRPCClient)(nil).EndSession), varargs...)
}

// StartSession mocks base method.
func (m *MockSessionRPCClient) StartSession(arg0 context.Context, arg1 *grpc.Request, arg2 ...grpc0.CallOption) (*grpc.Session, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "StartSession", varargs...)
	ret0, _ := ret[0].(*grpc.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartSession indicates an expected call of StartSession.
func (mr *MockSessionRPCClientMockRecorder) StartSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartSession", reflect.TypeOf((*MockSessionRPCClient)(nil).StartSession), varargs...)
}