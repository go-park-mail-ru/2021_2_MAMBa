package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	mock2 "2021_2_MAMBa/internal/pkg/reviews/usecase/mock"
	authRPC "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	mockSessions "2021_2_MAMBa/internal/pkg/sessions/mock"
	userDel "2021_2_MAMBa/internal/pkg/user/delivery/http"
	mock3 "2021_2_MAMBa/internal/pkg/user/usecase/mock"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
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

var testTableGetReviewSuccess = [...]testRow{
	{
		inQuery: "id=8",
		out:     `{"id":8,"film_id":8,"film_title_ru":"Гарри Поттер и узник Азкабана","film_title_original":"Harry Potter and the Prisoner of Azkaban","film_picture_url":"server/images/harry3.webp","author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"}` + "\n",
		status:  http.StatusOK,
		name:    `empty works`,
	},
}

var testTableGetReviewFailure = [...]testRow{
	{
		inQuery: "id=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `overshoot`,
	},
	{
		inQuery: "id=-2",
		out:     customErrors.ErrIdMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `neg skip`,
	},
}

func TestGetReviewSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getReview?"
	for _, test := range testTableGetReviewSuccess {
		var cl domain.Review
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockReviewUsecase(ctrl)
		mock.EXPECT().GetReview(uint64(8)).Times(1).Return(cl, nil)
		handler := ReviewHandler{ReiviewUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetReview(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetReviewFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getReview?"
	for i, test := range testTableGetReviewFailure {
		var cl domain.Review
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockReviewUsecase(ctrl)
		if i == 0 {
			mock.EXPECT().GetReview(uint64(10)).Times(1).Return(domain.Review{}, customErrors.ErrorSkip)
		}
		handler := ReviewHandler{ReiviewUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetReview(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTablePostRatingSuccess = [...]testRow{
	{
		inQuery:    "",
		out:        `{"id":7,"film_id":7,"film_title_ru":"Гарри Поттер и Тайная комната","film_title_original":"Harry Potter and the Chamber of Secrets","film_picture_url":"server/images/harry2.webp","author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"и так тоже","review_type":2,"stars":8,"date":"2021-10-31T00:00:00Z"}` + "\n",
		bodyString: `{"film_id":7,"review_text":"и так тоже","review_type":2}`,
		status:     http.StatusOK,
		name:       `normal`,
	},
}

func TestPostReviewSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/postReview?"
	test1 := testRow{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "123456"}`,
		out:        `{"id":2,"first_name":"Ivan","surname":"Ivanov","email":"iva21@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusOK,
		name:       "log in user",
	}
	for _, test := range testTablePostRatingSuccess {
		mock := mock3.NewMockUserUsecase(ctrl)
		mockSessions1 := mockSessions.NewMockSessionRPCClient(ctrl)
		var cl domain.UserToLogin
		var ret domain.User
		_ = json.Unmarshal([]byte(test1.bodyString), &cl)
		_ = json.Unmarshal([]byte(test1.out), &ret)
		handler := userDel.UserHandler{UserUsecase: mock, AuthClient: mockSessions1}
		mock.EXPECT().Login(&cl).Times(1).Return(ret, nil)
		bodyReader := strings.NewReader(test1.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/login"+test1.inQuery, bodyReader)
		mockSessions1.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).Times(1)
		mockSessions1.EXPECT().StartSession(r.Context(), &authRPC.Request{ID: 2}).Return(&authRPC.Session{Name: "session-name", Value: "aaa"}, nil)
		handler.Login(w, r)

		var toPost domain.Review
		var res domain.Review
		_ = json.Unmarshal([]byte(test.bodyString), &toPost)
		_ = json.Unmarshal([]byte(test.out), &res)
		mock2 := mock2.NewMockReviewUsecase(ctrl)
		mockSessions2 := mockSessions.NewMockSessionRPCClient(ctrl)
		toPost.AuthorId = uint64(2)
		mock2.EXPECT().PostReview(toPost).Times(1).Return(res, nil)
		handler2 := ReviewHandler{ReiviewUsecase: mock2, AuthClient: mockSessions2}
		bodyReader = strings.NewReader(test.bodyString)
		r = httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
		w = httptest.NewRecorder()
		mockSessions2.EXPECT().CheckSession(r.Context(), &authRPC.Request{Name: "session-name", Value: "aaa"}).Return(&authRPC.ID{ID: 2}, nil).Times(1)
		handler2.PostReview(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

var testTableGetReviewsSuccess = [...]testRow{
	{
		inQuery: "id=13&film_id=8&skips=0&limits=10",
		out:     `{"review_list":[{"id":8,"film_id":8,"film_title_ru":"Гарри Поттер и узник Азкабана","film_title_original":"Harry Potter and the Prisoner of Azkaban","film_picture_url":"server/images/harry3.webp","author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `full works`,
		skip:    0,
		limit:   10,
	},
	{
		inQuery: "id=13&film_id=8",
		out:     `{"review_list":[{"id":8,"film_id":8,"film_title_ru":"Гарри Поттер и узник Азкабана","film_title_original":"Harry Potter and the Prisoner of Azkaban","film_picture_url":"server/images/harry3.webp","author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `empty works`,
		skip:    0,
		limit:   10,
	},
}
var testTableGetReviewsFailure = [...]testRow{
	{
		inQuery: "id=13&film_id=8&skip=-1&limit=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
	},
	{
		inQuery: "id=13&film_id=8&skip_reviews=11&limit=-2",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    11,
		limit:   -2,
	},
	{
		inQuery: "id=13&film_id=8&skip=14&limit=1",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `skip overshoot`,
		skip:    14,
		limit:   1,
	},
}

func TestGetReviewsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/loadReviewsExcept?"
	for _, test := range testTableGetReviewsSuccess {
		var cl domain.FilmReviews
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockReviewUsecase(ctrl)
		mock.EXPECT().LoadReviewsExcept(uint64(13), uint64(8), test.skip, test.limit).Return(cl, nil)
		handler := ReviewHandler{ReiviewUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.LoadExcept(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetReviewsFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/loadReviewsExcept?"
	for i, test := range testTableGetReviewsFailure {
		mock := mock2.NewMockReviewUsecase(ctrl)
		if i == 2 {
			mock.EXPECT().LoadReviewsExcept(uint64(13), uint64(8), test.skip, test.limit).Return(domain.FilmReviews{}, customErrors.ErrorSkip)
		}
		handler := ReviewHandler{ReiviewUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.LoadExcept(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}
