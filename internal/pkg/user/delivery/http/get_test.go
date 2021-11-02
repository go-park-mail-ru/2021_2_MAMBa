package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/user"
	mock2 "2021_2_MAMBa/internal/pkg/user/usecase/mock"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testRow struct {
	inQuery    string
	bodyString string
	out        string
	status     int
	name       string
	skip       int
	limit      int
	skip1      int
	limit1     int
}

var testTableGetSuccess = [...]testRow{
	{
		inQuery:    "2",
		bodyString: ``,
		out:        `{"id":2,"first_name":"Ivan","surname":"Ivanov","email":"iva21@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusOK,
		name:       "find user",
	},
}

var testTableGetFailure = [...]testRow{
	{
		inQuery:    "3",
		bodyString: ``,
		out:        user.ErrorBadInput.Error() + "\n",
		status:     http.StatusNotFound,
		name:       "out of index",
	},
	{
		inQuery:    "a",
		bodyString: ``,
		out:        user.ErrorBadInput.Error() + "\n",
		status:     http.StatusBadRequest,
		name:       "no uinteger",
	},
	{
		inQuery:    "",
		bodyString: ``,
		out:        user.ErrorBadInput.Error() + "\n",
		status:     http.StatusBadRequest,
		name:       "empty",
	},
}

func TestGetBasicInfoSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range testTableGetSuccess {
		mock := mock2.NewMockUserUsecase(ctrl)
		var cl domain.User
		_ = json.Unmarshal([]byte(test.out), &cl)
		handler := UserHandler{UserUsecase: mock}
		mock.EXPECT().GetBasicInfo(uint64(2)).Times(1).Return(cl, nil)
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/user/get/"+test.inQuery, bodyReader)
		vars := map[string]string{
			"id": test.inQuery,
		}
		r = mux.SetURLVars(r, vars)

		handler.GetBasicInfo(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetBasicInfoFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for i, test := range testTableGetFailure {
		mock := mock2.NewMockUserUsecase(ctrl)
		var cl domain.User
		_ = json.Unmarshal([]byte(test.out), &cl)
		handler := UserHandler{UserUsecase: mock}
		if i==0 {
			mock.EXPECT().GetBasicInfo(uint64(3)).Times(1).Return(domain.User{}, user.ErrorNoUser)
		}
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/user/get/"+test.inQuery, bodyReader)
		vars := map[string]string{
			"id": test.inQuery,
		}
		r = mux.SetURLVars(r, vars)

		handler.GetBasicInfo(w, r)
		assert.Equal(t, test.out, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}