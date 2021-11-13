package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	mock2 "2021_2_MAMBa/internal/pkg/user/usecase/mock"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		out:        customErrors.ErrIdMsg,
		status:     http.StatusNotFound,
		name:       "out of index",
	},
	{
		inQuery:    "a",
		bodyString: ``,
		out:        customErrors.ErrorBadInput.Error(),
		status:     http.StatusBadRequest,
		name:       "no uinteger",
	},
	{
		inQuery:    "",
		bodyString: ``,
		out:        customErrors.ErrorBadInput.Error(),
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
		result := `{"body":` + test.out + `,"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
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
		if i == 0 {
			mock.EXPECT().GetBasicInfo(uint64(3)).Times(1).Return(domain.User{}, customErrors.ErrorNoUser)
		}
		bodyReader := strings.NewReader(test.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/api/user/get/"+test.inQuery, bodyReader)
		vars := map[string]string{
			"id": test.inQuery,
		}
		r = mux.SetURLVars(r, vars)

		handler.GetBasicInfo(w, r)
		result := `{"body":{"error":"` + test.out + `"},"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var LoginSuccess = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"email": "ivan@mail.ru","password": "1234abcd"}`,
		out:        `{"id":1,"first_name":"Иван","surname":"Иванов","email":"ivan@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusOK,
		name:       "log in user",
	},
}

var testTableGetProfileSuccess = [...]testRow{
	{
		inQuery:    "id=2",
		bodyString: ``,
		out:        `{"id":2,"first_name":"Алексей","surname":"Самойлов","picture_url":"/pic/1.jpg","email":"lexa@mail.ru","gender":"male","register_date":"2021-10-31T16:32:26.284085Z","sub_count":0,"bookmark_count":0,"am_subscribed":false}`,
		status:     http.StatusOK,
		name:       "find user",
	},
}

func TestGetProfileInfoSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	test := LoginSuccess[0]
	mock := mock2.NewMockUserUsecase(ctrl)
	var cl domain.UserToLogin
	var ret domain.User
	_ = json.Unmarshal([]byte(test.bodyString), &cl)
	_ = json.Unmarshal([]byte(test.out), &ret)
	handler := UserHandler{UserUsecase: mock}
	mock.EXPECT().Login(&cl).Times(1).Return(ret, nil)
	bodyReader := strings.NewReader(test.bodyString)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/login", bodyReader)
	handler.Login(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	for _, testCase := range testTableGetProfileSuccess {
		r = httptest.NewRequest("GET", "/api/user/getProfile?"+testCase.inQuery, bodyReader)
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
		var in domain.User
		var ret1 domain.Profile
		_ = json.Unmarshal([]byte(testCase.bodyString), &in)
		_ = json.Unmarshal([]byte(testCase.out), &ret1)
		mock.EXPECT().GetProfileById(ret.ID, ret1.ID).Return(ret1, nil)
		w = httptest.NewRecorder()
		handler.GetProfile(w, r)
		result := `{"body":` + testCase.out + `,"status":` + fmt.Sprint(testCase.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, testCase.status, w.Code, "Test: "+test.name)
	}
}

var testTableUpdateProfileSuccess = [...]testRow{
	{
		inQuery:    "id=2",
		bodyString: `{"id":1,"first_name":"Алексей","surname":"Самойлов","picture_url":"/pic/1.jpg","email":"lexa@mail.ru","gender":"male" }`,
		out:        `{"id":1,"first_name":"Алексей","surname":"Самойлов","picture_url":"/pic/1.jpg","email":"lexa@mail.ru","gender":"male","register_date":"2021-10-31T16:32:26.284085Z","sub_count":0,"bookmark_count":0,"am_subscribed":false}`,
		status:     http.StatusOK,
		name:       "up profile",
	},
}

var testTableUpdateProfileFailure = [...]testRow{
	{
		inQuery:    "id=2",
		bodyString: `{"id":0,"first_name":"Алексей","surname":"Самойлов","picture_url":"/pic/1.jpg","email":"lexa@mail.ru","gender":"male" }`,
		out:        customErrors.ErrorInternalServer.Error(),
		status:     http.StatusInternalServerError,
		name:       "up profile fail",
	},
}

func TestUpdateProfileInfoSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	test := LoginSuccess[0]
	mock := mock2.NewMockUserUsecase(ctrl)
	var cl domain.UserToLogin
	var ret domain.User
	_ = json.Unmarshal([]byte(test.bodyString), &cl)
	_ = json.Unmarshal([]byte(test.out), &ret)
	handler := UserHandler{UserUsecase: mock}
	mock.EXPECT().Login(&cl).Times(1).Return(ret, nil)
	bodyReader := strings.NewReader(test.bodyString)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/login", bodyReader)
	handler.Login(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	for _, testCase := range testTableUpdateProfileSuccess {
		bodyReader := strings.NewReader(testCase.bodyString)
		r = httptest.NewRequest("GET", "/api/user/changeProfile?"+testCase.inQuery, bodyReader)
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
		var in domain.Profile
		var ret1 domain.Profile
		_ = json.Unmarshal([]byte(testCase.bodyString), &in)
		_ = json.Unmarshal([]byte(testCase.out), &ret1)
		mock.EXPECT().UpdateProfile(in).Return(ret1, nil)
		w = httptest.NewRecorder()
		handler.UpdateProfile(w, r)
		result := `{"body":` + testCase.out + `,"status":` + fmt.Sprint(testCase.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, testCase.status, w.Code, "Test: "+testCase.name)
	}
}

func TestUpdateProfileInfoFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mock2.NewMockUserUsecase(ctrl)
	handler := UserHandler{UserUsecase: mock}
	for _, testCase := range testTableUpdateProfileFailure {
		bodyReader := strings.NewReader(testCase.bodyString)
		r := httptest.NewRequest("GET", "/api/user/changeProfile?"+testCase.inQuery, bodyReader)
		w := httptest.NewRecorder()
		w = httptest.NewRecorder()
		handler.UpdateProfile(w, r)
		result := `{"body":{"error":"` + testCase.out + `"},"status":` + fmt.Sprint(testCase.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+testCase.name)
	}
}

var testTableSubscribeSuccess = [...]testRow{
	{
		inQuery:    "",
		bodyString: `{"id":2}`,
		out:        `{"id":2,"first_name":"Алексей","surname":"Самойлов","picture_url":"/pic/1.jpg","email":"lexa@mail.ru","gender":"male","register_date":"2021-10-31T16:32:26.284085Z","sub_count":0,"bookmark_count":0,"am_subscribed":true}`,
		status:     http.StatusOK,
		name:       "sub",
	},
}

var testTableSubscribeFailure = [...]testRow{
	{
		inQuery:    "id=2",
		bodyString: `{"id":2}`,
		out:        customErrors.ErrorInternalServer.Error(),
		status:     http.StatusInternalServerError,
		name:       "up profile fail",
	},
}

func TestSubscribeSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	test := LoginSuccess[0]
	mock := mock2.NewMockUserUsecase(ctrl)
	var cl domain.UserToLogin
	var ret domain.User
	_ = json.Unmarshal([]byte(test.bodyString), &cl)
	_ = json.Unmarshal([]byte(test.out), &ret)
	handler := UserHandler{UserUsecase: mock}
	mock.EXPECT().Login(&cl).Times(1).Return(ret, nil)
	bodyReader := strings.NewReader(test.bodyString)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/api/login", bodyReader)
	handler.Login(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	for _, testCase := range testTableSubscribeSuccess {
		bodyReader := strings.NewReader(testCase.bodyString)
		r = httptest.NewRequest("GET", "/api/user/subscribeTo?"+testCase.inQuery, bodyReader)
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
		var in domain.Profile
		var ret1 domain.Profile
		_ = json.Unmarshal([]byte(testCase.bodyString), &in)
		_ = json.Unmarshal([]byte(testCase.out), &ret1)
		mock.EXPECT().CreateSubscription(ret.ID, ret1.ID).Return(ret1, nil)
		w = httptest.NewRecorder()
		handler.CreateSubscription(w, r)
		result := `{"body":` + testCase.out + `,"status":` + fmt.Sprint(testCase.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, testCase.status, w.Code, "Test: "+testCase.name)
	}
}

func TestSubscribeFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mock2.NewMockUserUsecase(ctrl)
	handler := UserHandler{UserUsecase: mock}
	for _, testCase := range testTableSubscribeFailure {
		bodyReader := strings.NewReader(testCase.bodyString)
		r := httptest.NewRequest("GET", "/api/user/subscribeTo?"+testCase.inQuery, bodyReader)
		w := httptest.NewRecorder()
		w = httptest.NewRecorder()
		handler.CreateSubscription(w, r)
		result := `{"body":{"error":"` + testCase.out + `"},"status":` + fmt.Sprint(testCase.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+testCase.name)
	}
}

var testTableGetReviewsSuccess = [...]testRow{
	{
		inQuery: "id=1&skips=0&limits=10",
		out:     `{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `full works`,
		skip:    0,
		limit:   10,
	},
	{
		inQuery: "id=1",
		out:     `{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `empty works`,
		skip:    0,
		limit:   10,
	},
}
var testTableGetReviewsFailure = [...]testRow{
	{
		inQuery: "id=1&skip=-1&limit=10",
		out:     customErrors.ErrSkipMsg,
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
	},
	{
		inQuery: "id=1&skip_reviews=11&limit=-2",
		out:     customErrors.ErrLimitMsg,
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    11,
		limit:   -2,
	},
	{
		inQuery: "id=1&skip=14&limit=1",
		out:     customErrors.ErrSkipMsg,
		status:  http.StatusBadRequest,
		name:    `skip overshoot`,
		skip:    14,
		limit:   1,
	},
}

func TestGetReviewsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/user/getReviewsAndStars?"
	for _, test := range testTableGetReviewsSuccess {
		var cl domain.FilmReviews
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockUserUsecase(ctrl)
		mock.EXPECT().LoadUserReviews(uint64(1), test.skip, test.limit).Return(cl, nil)
		handler := UserHandler{UserUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.LoadUserReviews(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetReviewsFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/loadFilmReviews?"
	for i, test := range testTableGetReviewsFailure {
		mock := mock2.NewMockUserUsecase(ctrl)
		if i == 2 {
			mock.EXPECT().LoadUserReviews(uint64(1), test.skip, test.limit).Return(domain.FilmReviews{}, customErrors.ErrorSkip)
		}
		handler := UserHandler{UserUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.LoadUserReviews(w, r)
		result := `{"body":{"error":"` + test.out + `"},"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}
