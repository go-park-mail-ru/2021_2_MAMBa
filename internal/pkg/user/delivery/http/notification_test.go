package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	mock2 "2021_2_MAMBa/internal/pkg/user/usecase/mock"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testTablePostFailure = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"token":"eDZKmx6C-QtDqqF6ALphK9:APA91bFqOmRnGjbBNukcyOS-zGBME2c7rPYJrTqzMeO8WZDMjfYp2MLrnl8YGgMdCZWK-aZtmxrqdgjCEv0lLBLHqTnPm7Wg6wnt7yS4RbhOqewY8_sMqMbBQmqgBXMcl7k5ayblZWRl" }`,
		out:        `{"error":"error - bad input"}`,
		status:     http.StatusBadRequest,
		name:       "subscribe topic",
	},
}

func TestAddUserToNotificationTopicFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range testTablePostFailure {
		mock := mock2.NewMockUserUsecase(ctrl)
		var cl domain.User
		_ = json.Unmarshal([]byte(test.out), &cl)
		handler := UserHandler{UserUsecase: mock}
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/subscribePush/"+test.inQuery, bodyReader)

		handler.AddUserToNotificationTopic(w, r)
		result := `{"body":` + test.out + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

var testTablePostFailure2 = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"title":"aaa","description": "bbb" }`,
		out:        `"error - bad input"`,
		status:     http.StatusBadRequest,
		name:       "send push topic",
	},
}

func TestSendPushToAllFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range testTablePostFailure2 {
		mock := mock2.NewMockUserUsecase(ctrl)
		var cl domain.User
		_ = json.Unmarshal([]byte(test.out), &cl)
		handler := UserHandler{UserUsecase: mock}
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/user/get/"+test.inQuery, bodyReader)

		handler.SendPushToAll(w, r)
		result := `{"body":{"error":` + test.out + `},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}
