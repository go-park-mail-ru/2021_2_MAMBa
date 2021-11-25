package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	mock2 "2021_2_MAMBa/internal/pkg/film/usecase/mock"
	authRPC "2021_2_MAMBa/internal/pkg/sessions/delivery/grpc"
	mockSessions "2021_2_MAMBa/internal/pkg/sessions/mock"
	userDel "2021_2_MAMBa/internal/pkg/user/delivery/http"
	mock3 "2021_2_MAMBa/internal/pkg/user/usecase/mock"
	"encoding/json"
	"errors"
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

var testTableGetFilmSuccess = [...]testRow{
	{
		inQuery: "id=8&skip_reviews=0&limit_reviews=10&skip_recommend=0&limit_recommend=10",
		out:     `{"film":{"id":8,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":8.5,"description":"третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер,Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","total_revenue":"$795,634,069.00","poster_url":"server/images/harry3.webp","trailer_url":"trailer","content_type":"фильм","release_year":2004,"duration":142,"origin_countries":["Великобритания","США"],"director":{"id":21,"name_rus":"Крис Коламбус","career":[""]},"screenwriter":{"id":26,"name_rus":"Альфонсо Куарон","career":[""]}},"reviews":{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"},{"id":13,"film_id":8,"author_name":"Максим Дудник","author_picture_url":"/pic/1.jpg","review_text":"ffff","review_type":1,"stars":0,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10},"recommendations":{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":0.0,"poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":0.0,"poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10},"my_review":{"id":0,"film_id":0,"review_type":0,"stars":0,"date":""},"bookmarked":false}` + "\n",
		status:  http.StatusOK,
		name:    `full works`,
		skip:    0,
		limit:   10,
		skip1:   0,
		limit1:  10,
	},
	{
		inQuery: "id=8",
		out:     `{"film":{"id":8,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":8.5,"description":"третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер,Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","total_revenue":"$795,634,069.00","poster_url":"server/images/harry3.webp","trailer_url":"trailer","content_type":"фильм","release_year":2004,"duration":142,"origin_countries":["Великобритания","США"],"director":{"id":21,"name_rus":"Крис Коламбус","career":[""]},"screenwriter":{"id":26,"name_rus":"Альфонсо Куарон","career":[""]}},"reviews":{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"},{"id":13,"film_id":8,"author_name":"Максим Дудник","author_picture_url":"/pic/1.jpg","review_text":"ffff","review_type":1,"stars":0,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10},"recommendations":{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":0.0,"poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":0.0,"poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10},"my_review":{"id":0,"film_id":0,"review_type":0,"stars":0,"date":""},"bookmarked":false}` + "\n",
		status:  http.StatusOK,
		name:    `empty works`,
		skip:    0,
		limit:   10,
		skip1:   0,
		limit1:  10,
	},
}
var testTableGetFilmFailure = [...]testRow{
	{
		inQuery: "id=8&skip_reviews=-1&limit_reviews=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
		skip1:   0,
		limit1:  10,
	},
	{
		inQuery: "id=8&skip_reviews=11&limit_reviews=-2",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    11,
		limit:   -2,
		skip1:   0,
		limit1:  10,
	},
	{
		inQuery: "id=8&skip_reviews=14&limit_reviews=1",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `skip overshoot`,
		skip:    14,
		limit:   1,
		skip1:   0,
		limit1:  10,
	},
	{
		inQuery: "id=8&skip_recommend=-1&limit_recommend=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip1:   -1,
		limit1:  10,
		skip:    0,
		limit:   10,
	},
	{
		inQuery: "id=8&skip_recommend=11&limit_recommend=-2",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip1:   11,
		limit1:  -2,
		skip:    0,
		limit:   10,
	},
	{
		inQuery: "id=8&skip_recommend=14&limit_recommend=1",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `skip overshoot`,
		skip1:   14,
		limit1:  1,
		skip:    0,
		limit:   10,
	},
}

func TestGetFilmSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getFilm?"
	for _, test := range testTableGetFilmSuccess {
		var cl domain.FilmPageInfo
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mockSessions := mockSessions.NewMockSessionRPCClient(ctrl)
		mock.EXPECT().GetFilm(uint64(0), uint64(8), test.skip, test.limit, test.skip1, test.limit1).Times(1).Return(cl, nil)
		handler := FilmHandler{FilmUsecase: mock, AuthClient: mockSessions}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, errors.New(customErrors.RPCErrUserNotLoggedIn)).Times(1)
		handler.GetFilm(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetFilmFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getFilm?"
	for i, test := range testTableGetFilmFailure {
		var cl domain.FilmPageInfo
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mockSessions := mockSessions.NewMockSessionRPCClient(ctrl)
		if i == 2 || i == 5 {
			mock.EXPECT().GetFilm(uint64(0), uint64(8), test.skip, test.limit, test.skip1, test.limit1).Times(1).Return(domain.FilmPageInfo{}, customErrors.ErrorSkip)
		}
		handler := FilmHandler{FilmUsecase: mock, AuthClient: mockSessions}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		mockSessions.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, errors.New(customErrors.RPCErrUserNotLoggedIn)).AnyTimes()
		handler.GetFilm(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTablePostRatingSuccess = [...]testRow{
	{
		inQuery:    "id=2&rating=10",
		out:        `{"rating":10.0}` + "\n",
		bodyString: `{"film_id":7,"review_text":"и так тоже","review_type":2}`,
		status:     http.StatusOK,
		name:       `normal`,
	},
}

func TestPostReviewSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/postRating?"
	test1 := testRow{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "123456"}`,
		out:        `{"id":2,"first_name":"Ivan","surname":"Ivanov","email":"iva21@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusOK,
		name:       "log in user",
	}
	for _, test := range testTablePostRatingSuccess {
		mock := mock3.NewMockUserUsecase(ctrl)
		var cl domain.UserToLogin
		var ret domain.User
		_ = json.Unmarshal([]byte(test1.bodyString), &cl)
		_ = json.Unmarshal([]byte(test1.out), &ret)
		mockSessions1 := mockSessions.NewMockSessionRPCClient(ctrl)
		handler := userDel.UserHandler{UserUsecase: mock, AuthClient: mockSessions1}
		mock.EXPECT().Login(&cl).Times(1).Return(ret, nil)
		bodyReader := strings.NewReader(test1.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/login"+test1.inQuery, bodyReader)
		mockSessions1.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, errors.New(customErrors.RPCErrUserNotLoggedIn)).Times(1)
		mockSessions1.EXPECT().StartSession(r.Context(), &authRPC.Request{ID: 2}).Return(&authRPC.Session{Name: "session-name", Value: "aaa"}, nil)
		handler.Login(w, r)

		_ = json.Unmarshal([]byte(test.bodyString), &cl)
		_ = json.Unmarshal([]byte(test.out), &ret)
		mock2 := mock2.NewMockFilmUsecase(ctrl)
		mockSessions2 := mockSessions.NewMockSessionRPCClient(ctrl)
		mock2.EXPECT().PostRating(uint64(2), uint64(2), 10.0).Times(1).Return(10.0, nil)
		handler2 := FilmHandler{FilmUsecase: mock2, AuthClient: mockSessions2}
		r = httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
		w = httptest.NewRecorder()
		mockSessions2.EXPECT().CheckSession(r.Context(), &authRPC.Request{Name: "session-name", Value: "aaa"}).Return(&authRPC.ID{ID: 2}, nil).Times(1)
		handler2.PostRating(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

var testTableGetMySuccess = [...]testRow{
	{
		inQuery:    "id=2",
		out:        `{"id":2,"film_id":2,"film_title_ru":"С любовью, Рози","film_title_original":"Love, Rosie","film_picture_url":"server/images/love-rosie.webp","author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":")","review_type":2,"stars":10,"date":"2021-10-31T00:00:00Z"}` + "\n",
		bodyString: ``,
		status:     http.StatusOK,
		name:       `normal`,
	},
}

func TestGetMySuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/loadMyReviewForFilm?"
	test1 := testRow{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "123456"}`,
		out:        `{"id":2,"first_name":"Ivan","surname":"Ivanov","email":"iva21@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusOK,
		name:       "log in user",
	}
	for _, test := range testTableGetMySuccess {
		mock := mock3.NewMockUserUsecase(ctrl)
		var cl domain.UserToLogin
		var ret domain.User
		_ = json.Unmarshal([]byte(test1.bodyString), &cl)
		_ = json.Unmarshal([]byte(test1.out), &ret)
		mockSessions1 := mockSessions.NewMockSessionRPCClient(ctrl)
		handler := userDel.UserHandler{UserUsecase: mock, AuthClient: mockSessions1}
		mock.EXPECT().Login(&cl).Times(1).Return(ret, nil)
		bodyReader := strings.NewReader(test1.bodyString)
		w := httptest.NewRecorder()
		r := httptest.NewRequest("POST", "/api/user/login"+test1.inQuery, bodyReader)
		mockSessions1.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, errors.New(customErrors.RPCErrUserNotLoggedIn)).Times(1)
		mockSessions1.EXPECT().StartSession(r.Context(), &authRPC.Request{ID: 2}).Return(&authRPC.Session{Name: "session-name", Value: "aaa"}, nil)
		handler.Login(w, r)

		var review domain.Review
		_ = json.Unmarshal([]byte(test.out), &review)
		mock2 := mock2.NewMockFilmUsecase(ctrl)
		mockSessions2 := mockSessions.NewMockSessionRPCClient(ctrl)
		mock2.EXPECT().LoadMyReview(uint64(2), uint64(2)).Return(review, nil)
		handler2 := FilmHandler{FilmUsecase: mock2, AuthClient: mockSessions2}
		r = httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
		w = httptest.NewRecorder()
		mockSessions2.EXPECT().CheckSession(r.Context(), &authRPC.Request{Name: "session-name", Value: "aaa"}).Return(&authRPC.ID{ID: 2}, nil).Times(1)
		handler2.LoadMyRv(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

var testTableGetReviewsSuccess = [...]testRow{
	{
		inQuery: "id=8&skips=0&limits=10",
		out:     `{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"},{"id":13,"film_id":8,"author_name":"Максим Дудник","author_picture_url":"/pic/1.jpg","review_text":"ffff","review_type":1,"stars":0,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `full works`,
		skip:    0,
		limit:   10,
	},
	{
		inQuery: "id=8",
		out:     `{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"},{"id":13,"film_id":8,"author_name":"Максим Дудник","author_picture_url":"/pic/1.jpg","review_text":"ffff","review_type":1,"stars":0,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `empty works`,
		skip:    0,
		limit:   10,
	},
}
var testTableGetReviewsFailure = [...]testRow{
	{
		inQuery: "id=8&skip=-1&limit=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
	},
	{
		inQuery: "id=8&skip_reviews=11&limit=-2",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    11,
		limit:   -2,
	},
	{
		inQuery: "id=8&skip=14&limit=1",
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
	apiPath := "/api/film/loadFilmReviews?"
	for _, test := range testTableGetReviewsSuccess {
		var cl domain.FilmReviews
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mock.EXPECT().LoadFilmReviews(uint64(8), test.skip, test.limit).Return(cl, nil)
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.loadFilmReviews(w, r)
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
		mock := mock2.NewMockFilmUsecase(ctrl)
		if i == 2 {
			mock.EXPECT().LoadFilmReviews(uint64(8), test.skip, test.limit).Return(domain.FilmReviews{}, customErrors.ErrorSkip)
		}
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.loadFilmReviews(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTableGetRecomSuccess = [...]testRow{
	{
		inQuery: "id=8&skips=0&limits=10",
		out:     `{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":0.0,"poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":0.0,"poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `full works`,
		skip:    0,
		limit:   10,
	},
	{
		inQuery: "id=8",
		out:     `{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":0.0,"poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":0.0,"poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `empty works`,
		skip:    0,
		limit:   10,
	},
}
var testTableGetRecomsFailure = [...]testRow{
	{
		inQuery: "id=8&skip=-1&limit=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
	},
	{
		inQuery: "id=8&skip_reviews=11&limit=-2",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    11,
		limit:   -2,
	},
	{
		inQuery: "id=8&skip=14&limit=1",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `skip overshoot`,
		skip:    14,
		limit:   1,
	},
}

func TestGetRecomSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/loadFilmRecommendations?"
	for _, test := range testTableGetRecomSuccess {
		var cl domain.FilmRecommendations
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mock.EXPECT().LoadFilmRecommendations(uint64(8), test.skip, test.limit).Return(cl, nil)
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.loadFilmRecommendations(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetRecomFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/loadFilmRecommendations?"
	for i, test := range testTableGetRecomsFailure {
		mock := mock2.NewMockFilmUsecase(ctrl)
		if i == 2 {
			mock.EXPECT().LoadFilmRecommendations(uint64(8), test.skip, test.limit).Return(domain.FilmRecommendations{}, customErrors.ErrorSkip)
		}
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.loadFilmRecommendations(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}
