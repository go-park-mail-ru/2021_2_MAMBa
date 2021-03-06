package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	authRPC "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	mockSessions "2021_2_MAMBa/internal/pkg/sessions/mock"
	mock2 "2021_2_MAMBa/internal/pkg/user/usecase/mock"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testTableRegisterSuccess = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"first_name": "Ivan","surname": "Ivanov","email": "ivan1@mail.ru","password": "123456","password_repeat": "123456"}`,
		out:        `{"id":1,"first_name":"Ivan","surname":"Ivanov","email":"ivan1@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusCreated,
		name:       "register one",
	},
}
var testTableRegisterFailure = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"first_name": "Ivan","surname": "Ivanov","email": "ivan1@mail.ru","password": "123456","password_repeat": "123456"}`,
		out:        customErrors.ErrorAlreadyExists.Error(),
		status:     http.StatusConflict,
		name:       "already in",
	},
	{
		inQuery:    "",
		bodyString: `{"first_nme": "Ivan",}`,
		out:        customErrors.ErrorBadInput.Error(),
		status:     http.StatusBadRequest,
		name:       "bad fields",
	},
	{
		inQuery:    "",
		bodyString: `{"first_name": "Ivan",}`,
		out:        customErrors.ErrorBadInput.Error(),
		status:     http.StatusBadRequest,
		name:       "empty fields",
	},
	{
		inQuery:    "",
		bodyString: `{"first_name": "Ivan12","surname": "Ivanov","email": "ivan131@mail.ru","password": "123455","password_repeat": "123456"}`,
		out:        customErrors.ErrorBadInput.Error(),
		status:     http.StatusBadRequest,
		name:       "unmatching passwords",
	},
}

func TestRegisterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range testTableRegisterSuccess {
		mock := mock2.NewMockUserUsecase(ctrl)
		mockSessions := mockSessions.NewMockSessionRPCClient(ctrl)
		var cl, ret domain.User
		_ = json.Unmarshal([]byte(test.bodyString), &cl)
		_ = json.Unmarshal([]byte(test.out), &ret)
		handler := UserHandler{UserUsecase: mock, AuthClient: mockSessions}
		mock.EXPECT().Register(&cl).Times(1).Return(ret, nil)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/register"+test.inQuery, bodyReader)
		mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, nil).AnyTimes()
		mockSessions.EXPECT().StartSession(r.Context(), &authRPC.Request{ID: 1}).Return(&authRPC.Session{Name: "session-name"}, nil)
		handler.Register(w, r)
		result := `{"body":` + test.out + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, 201, w.Code, "Test: "+test.name)
	}
}
func TestRegisterFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for i, test := range testTableRegisterFailure {
		mock := mock2.NewMockUserUsecase(ctrl)
		var cl, ret domain.User
		_ = json.Unmarshal([]byte(test.bodyString), &cl)
		_ = json.Unmarshal([]byte(test.out), &ret)
		handler := UserHandler{UserUsecase: mock}
		if i == 0 {
			mock.EXPECT().Register(&cl).Times(1).Return(ret, customErrors.ErrorAlreadyExists)
		}
		if i == 3 {
			mock.EXPECT().Register(&cl).Times(1).Return(domain.User{}, customErrors.ErrorBadInput)
		}
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/register"+test.inQuery, bodyReader)
		handler.Register(w, r)
		result := `{"body":{"error":"` + test.out + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTableLoginSuccess = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "123456"}`,
		out:        `{"id":2,"first_name":"Ivan","surname":"Ivanov","email":"iva21@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusOK,
		name:       "log in user",
	},
}

var testTableLoginFailure = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"email": "raddom@mail.su","password": "123456"}`,
		out:        customErrors.ErrorBadCredentials.Error(),
		status:     http.StatusUnauthorized,
		name:       "user not in base",
	},
	{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "122456"}`,
		out:        customErrors.ErrorBadCredentials.Error(),
		status:     http.StatusUnauthorized,
		name:       "wrong pass",
	},
	{
		inQuery:    "",
		bodyString: `{"password": "122456"}`,
		out:        customErrors.ErrorBadInput.Error(),
		status:     http.StatusBadRequest,
		name:       "no email",
	},
	{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru"}`,
		out:        customErrors.ErrorBadInput.Error(),
		status:     http.StatusBadRequest,
		name:       "no pass",
	},
	{
		inQuery:    "",
		bodyString: `{"emal": "iva21@mail.ru","password": "123456"}`,
		out:        customErrors.ErrorBadInput.Error(),
		status:     http.StatusBadRequest,
		name:       "wrong json",
	},
}

func TestLoginSuccess(t *testing.T) {
	apiPath := "/api/user/login"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range testTableLoginSuccess {
		mock := mock2.NewMockUserUsecase(ctrl)
		mockSessions := mockSessions.NewMockSessionRPCClient(ctrl)
		var cl domain.UserToLogin
		var ret domain.User
		_ = json.Unmarshal([]byte(test.bodyString), &cl)
		_ = json.Unmarshal([]byte(test.out), &ret)
		handler := UserHandler{UserUsecase: mock, AuthClient: mockSessions}
		mock.EXPECT().Login(&cl).Times(1).Return(ret, nil)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).AnyTimes()
		mockSessions.EXPECT().StartSession(r.Context(), &authRPC.Request{ID: 2}).Return(&authRPC.Session{Name: "session-name"}, nil)
		handler.Login(w, r)
		result := `{"body":` + test.out + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestLoginFailure(t *testing.T) {
	apiPath := "/api/user/login"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for i, test := range testTableLoginFailure {
		mock := mock2.NewMockUserUsecase(ctrl)
		mockSessions := mockSessions.NewMockSessionRPCClient(ctrl)
		var cl domain.UserToLogin
		var ret domain.User
		_ = json.Unmarshal([]byte(test.bodyString), &cl)
		_ = json.Unmarshal([]byte(test.out), &ret)
		handler := UserHandler{UserUsecase: mock, AuthClient: mockSessions}
		if i <= 1 {
			mock.EXPECT().Login(&cl).Times(1).Return(domain.User{}, customErrors.ErrorBadCredentials)
		} else if i <= 4 {
			mock.EXPECT().Login(&cl).Times(1).Return(domain.User{}, customErrors.ErrorBadInput)
		}
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).AnyTimes()
		handler.Login(w, r)
		result := `{"body":{"error":"` + test.out + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTableLogoutFailure = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "123456"}`,
		out:        customErrors.ErrorUserNotLoggedIn.Error(),
		status:     http.StatusForbidden,
		name:       "logout not logged in",
	},
}

func TestLogoutFailure(t *testing.T) {
	apiPath := "/api/user/logout"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range testTableLogoutFailure {
		mock := mock2.NewMockUserUsecase(ctrl)
		mockSessions := mockSessions.NewMockSessionRPCClient(ctrl)
		handler := UserHandler{UserUsecase: mock, AuthClient: mockSessions}
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).AnyTimes()
		handler.Logout(w, r)
		result := `{"body":{"error":"` + test.out + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

func TestLogoutSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	test := testTableLoginSuccess[0]
	mock := mock2.NewMockUserUsecase(ctrl)
	mockSessions := mockSessions.NewMockSessionRPCClient(ctrl)
	var cl domain.UserToLogin
	var ret domain.User
	_ = json.Unmarshal([]byte(test.bodyString), &cl)
	_ = json.Unmarshal([]byte(test.out), &ret)
	handler := UserHandler{UserUsecase: mock, AuthClient: mockSessions}
	mock.EXPECT().Login(&cl).Times(1).Return(ret, nil)
	bodyReader := strings.NewReader(testTableLoginSuccess[0].bodyString)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/login", bodyReader)
	mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).Times(1)
	mockSessions.EXPECT().StartSession(r.Context(), &authRPC.Request{ID: 2}).Return(&authRPC.Session{Name: "session-name", Value: "aaa"}, nil)
	handler.Login(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	r = httptest.NewRequest("GET", "/api/user/logout", bodyReader)
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
	w = httptest.NewRecorder()
	mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{Name: "session-name", Value: "aaa"}).Return(&authRPC.ID{ID: 2}, nil).AnyTimes()
	mockSessions.EXPECT().EndSession(r.Context(), &authRPC.Request{Name: "session-name", Value: "aaa", ID: 2}).Return(&authRPC.Session{Name: "session-name"}, nil).AnyTimes()
	handler.Logout(w, r)
	assert.Equal(t, http.StatusOK, w.Code, "Test: Logout OK")
}

func TestCheckAuthSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	test := testTableLoginSuccess[0]
	mock := mock2.NewMockUserUsecase(ctrl)
	mockSessions := mockSessions.NewMockSessionRPCClient(ctrl)
	var cl domain.UserToLogin
	var ret domain.User
	_ = json.Unmarshal([]byte(test.bodyString), &cl)
	_ = json.Unmarshal([]byte(test.out), &ret)
	handler := UserHandler{UserUsecase: mock, AuthClient: mockSessions}
	mock.EXPECT().Login(&cl).Times(1).Return(ret, nil)
	bodyReader := strings.NewReader(testTableLoginSuccess[0].bodyString)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/login", bodyReader)
	mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).Times(1)
	mockSessions.EXPECT().StartSession(r.Context(), &authRPC.Request{ID: 2}).Return(&authRPC.Session{Name: "session-name", Value: "aaa"}, nil)
	handler.Login(w, r)
	require.Equal(t, http.StatusOK, w.Code)

	r = httptest.NewRequest("GET", "/api/user/checkAuth", bodyReader)
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		r.AddCookie(cookie)
	}
	w = httptest.NewRecorder()
	mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{Name: "session-name", Value: "aaa"}).Return(&authRPC.ID{ID: 2}, nil).AnyTimes()
	mock.EXPECT().CheckAuth(uint64(2)).Return(ret, nil)
	handler.CheckAuth(w, r)
	result := `{"body":` + test.out + `,"status":` + fmt.Sprint(test.status) + "}"
	assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	assert.Equal(t, http.StatusOK, w.Code, "Test: CheckAuth ok")
}

func TestCheckAuthFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mock2.NewMockUserUsecase(ctrl)
	mockSessions := mockSessions.NewMockSessionRPCClient(ctrl)
	test := testTableLogoutFailure[0]
	handler := UserHandler{UserUsecase: mock, AuthClient: mockSessions}
	bodyReader := strings.NewReader("")
	r := httptest.NewRequest("GET", "/api/user/checkAuth", bodyReader)
	w := httptest.NewRecorder()
	mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).AnyTimes()
	handler.CheckAuth(w, r)
	result := `{"body":{"error":"` + test.out + `"},"status":` + fmt.Sprint(test.status) + "}"
	assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
}
