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
	resultJSON string
	status     int
	name       string
	skip       int
	limit      int
	skip1      int
	limit1     int
}

var testTableGetFilmSuccess = [...]testRow{
	{
		inQuery:    "id=8&skip_reviews=0&limit_reviews=10&skip_recommend=0&limit_recommend=10",
		out:        `{"film":{"id":8,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":8.5,"description":"третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер,Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","total_revenue":"$795,634,069.00","poster_url":"server/images/harry3.webp","trailer_url":"trailer","content_type":"фильм","release_year":2004,"duration":142,"origin_countries":["Великобритания","США"],"director":{"id":21,"name_rus":"Крис Коламбус","career":[""]},"screenwriter":{"id":26,"name_rus":"Альфонсо Куарон","career":[""]}},"reviews":{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"},{"id":13,"film_id":8,"author_name":"Максим Дудник","author_picture_url":"/pic/1.jpg","review_text":"ffff","review_type":1,"stars":0,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10},"recommendations":{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":0.0,"poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":0.0,"poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10},"my_review":{"id":0,"film_id":0,"review_type":0,"stars":0,"date":""},"bookmarked":false}` + "\n",
		resultJSON: `{"film":{"id":8,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":"8.5","description":"третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер,Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","total_revenue":"$795,634,069.00","poster_url":"server/images/harry3.webp","trailer_url":"trailer","content_type":"фильм","release_year":2004,"duration":142,"origin_countries":["Великобритания","США"],"director":{"id":21,"name_rus":"Крис Коламбус","career":[""]},"screenwriter":{"id":26,"name_rus":"Альфонсо Куарон","career":[""]}},"reviews":{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"},{"id":13,"film_id":8,"author_name":"Максим Дудник","author_picture_url":"/pic/1.jpg","review_text":"ffff","review_type":1,"stars":0,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10},"recommendations":{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":"0.0","poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":"0.0","poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10},"my_review":{"id":0,"film_id":0,"review_type":0,"stars":0,"date":""},"bookmarked":false}` + "\n",
		status:     http.StatusOK,
		name:       `full works`,
		skip:       0,
		limit:      10,
		skip1:      0,
		limit1:     10,
	},
	{
		inQuery:    "id=8",
		out:        `{"film":{"id":8,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":8.5,"description":"третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер,Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","total_revenue":"$795,634,069.00","poster_url":"server/images/harry3.webp","trailer_url":"trailer","content_type":"фильм","release_year":2004,"duration":142,"origin_countries":["Великобритания","США"],"director":{"id":21,"name_rus":"Крис Коламбус","career":[""]},"screenwriter":{"id":26,"name_rus":"Альфонсо Куарон","career":[""]}},"reviews":{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"},{"id":13,"film_id":8,"author_name":"Максим Дудник","author_picture_url":"/pic/1.jpg","review_text":"ffff","review_type":1,"stars":0,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10},"recommendations":{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":0.0,"poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":0.0,"poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10},"my_review":{"id":0,"film_id":0,"review_type":0,"stars":0,"date":""},"bookmarked":false}` + "\n",
		resultJSON: `{"film":{"id":8,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":"8.5","description":"третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер,Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","total_revenue":"$795,634,069.00","poster_url":"server/images/harry3.webp","trailer_url":"trailer","content_type":"фильм","release_year":2004,"duration":142,"origin_countries":["Великобритания","США"],"director":{"id":21,"name_rus":"Крис Коламбус","career":[""]},"screenwriter":{"id":26,"name_rus":"Альфонсо Куарон","career":[""]}},"reviews":{"review_list":[{"id":8,"film_id":8,"author_name":"Иван Иванов","author_picture_url":"/pic/1.jpg","review_text":"отвал башки","review_type":3,"stars":10,"date":"2021-10-31T00:00:00Z"},{"id":13,"film_id":8,"author_name":"Максим Дудник","author_picture_url":"/pic/1.jpg","review_text":"ffff","review_type":1,"stars":0,"date":"2021-10-31T00:00:00Z"}],"more_available":false,"review_total":2,"current_sort":"","current_limit":10,"current_skip":10},"recommendations":{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":"0.0","poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":"0.0","poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10},"my_review":{"id":0,"film_id":0,"review_type":0,"stars":0,"date":""},"bookmarked":false}` + "\n",
		status:     http.StatusOK,
		name:       `empty works`,
		skip:       0,
		limit:      10,
		skip1:      0,
		limit1:     10,
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
		result := `{"body":` + test.resultJSON[:len(test.resultJSON)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
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
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTablePostRatingSuccess = [...]testRow{
	{
		inQuery:    "id=2&rating=10",
		out:        `{"rating":"10.0"}` + "\n",
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
		mockSessions1.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).Times(1)
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
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
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
		mockSessions1.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).Times(1)
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
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
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
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
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
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTableGetRecomSuccess = [...]testRow{
	{
		inQuery:    "id=8&skips=0&limits=10",
		out:        `{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":0.0,"poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":0.0,"poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10}` + "\n",
		resultJSON: `{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":"0.0","poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":"0.0","poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10}` + "\n",
		status:     http.StatusOK,
		name:       `full works`,
		skip:       0,
		limit:      10,
	},
	{
		inQuery:    "id=8",
		out:        `{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":0.0,"poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":0.0,"poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10}` + "\n",
		resultJSON: `{"recommendation_list":[{"id":6,"title":"Гарри Поттер и философский камень","rating":"0.0","poster_url":"server/images/harry1.webp","director":{},"screenwriter":{}},{"id":7,"title":"Гарри Поттер и Тайная комната","rating":"0.0","poster_url":"server/images/harry2.webp","director":{},"screenwriter":{}}],"more_available":false,"recommendation_total":2,"current_limit":10,"current_skip":10}` + "\n",
		status:     http.StatusOK,
		name:       `empty works`,
		skip:       0,
		limit:      10,
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
		result := `{"body":` + test.resultJSON[:len(test.resultJSON)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
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
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTableGetMYSuccess = [...]testRow{
	{
		inQuery:    "id=8&skips=0&limits=10&month=2&year=2010",
		out:        `{"film_list":[{"id":1,"title":"Еще по одной","rating":0.0,"description":"В ресторане собираются учитель истории, психологии, музыки и физрук, чтобы отметить 40-летие одного из них. И решают проверить научную теорию о том, что c самого рождения человек страдает от нехватки алкоголя в крови, а чтобы стать по-настоящему счастливым, нужно быть немного нетрезвым. Друзья договариваются наблюдать, как возлияния скажутся на их работе и личной жизни, и устанавливают правила: не пить вечером и по выходным. Казалось бы, что может пойти не так?","poster_url":"server/images/one-more-drink.webp","release_year":2020,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10}` + "\n",
		resultJSON: `{"film_list":[{"id":1,"title":"Еще по одной","rating":"0.0","description":"В ресторане собираются учитель истории, психологии, музыки и физрук, чтобы отметить 40-летие одного из них. И решают проверить научную теорию о том, что c самого рождения человек страдает от нехватки алкоголя в крови, а чтобы стать по-настоящему счастливым, нужно быть немного нетрезвым. Друзья договариваются наблюдать, как возлияния скажутся на их работе и личной жизни, и устанавливают правила: не пить вечером и по выходным. Казалось бы, что может пойти не так?","poster_url":"server/images/one-more-drink.webp","release_year":2020,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10}` + "\n",
		status:     http.StatusOK,
		name:       `full works`,
		skip:       0,
		limit:      10,
	},
	{
		inQuery:    "id=8&month=2&year=2010",
		out:        `{"film_list":[{"id":1,"title":"Еще по одной","rating":0.0,"description":"В ресторане собираются учитель истории, психологии, музыки и физрук, чтобы отметить 40-летие одного из них. И решают проверить научную теорию о том, что c самого рождения человек страдает от нехватки алкоголя в крови, а чтобы стать по-настоящему счастливым, нужно быть немного нетрезвым. Друзья договариваются наблюдать, как возлияния скажутся на их работе и личной жизни, и устанавливают правила: не пить вечером и по выходным. Казалось бы, что может пойти не так?","poster_url":"server/images/one-more-drink.webp","release_year":2020,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10}` + "\n",
		resultJSON: `{"film_list":[{"id":1,"title":"Еще по одной","rating":"0.0","description":"В ресторане собираются учитель истории, психологии, музыки и физрук, чтобы отметить 40-летие одного из них. И решают проверить научную теорию о том, что c самого рождения человек страдает от нехватки алкоголя в крови, а чтобы стать по-настоящему счастливым, нужно быть немного нетрезвым. Друзья договариваются наблюдать, как возлияния скажутся на их работе и личной жизни, и устанавливают правила: не пить вечером и по выходным. Казалось бы, что может пойти не так?","poster_url":"server/images/one-more-drink.webp","release_year":2020,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10}` + "\n",
		status:     http.StatusOK,
		name:       `empty works`,
		skip:       0,
		limit:      10,
	},
}

var testTableGetMYFailure = [...]testRow{
	{
		inQuery: "id=8&skip=-1&limit=10&month=2&year=2010",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
	},
	{
		inQuery: "id=8&skip=11&limit=-2&month=2&year=2010",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    11,
		limit:   -2,
	},
	{
		inQuery: "id=8&skip=14&limit=1&month=2&year=2010",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `skip overshoot`,
		skip:    14,
		limit:   1,
	},
	{
		inQuery: "id=8&skip=1&limit=10&month=-2&year=2010",
		out:     customErrors.ErrDateMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
	},
	{
		inQuery: "id=8&skip=11&limit=2&month=2&year=10",
		out:     customErrors.ErrDateMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    11,
		limit:   -2,
	},
}

func TestGetMYSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/calendar?"
	for _, test := range testTableGetMYSuccess {
		var cl domain.FilmList
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mock.EXPECT().GetFilmsByMonthYear(2, 2010, 10, 0).Return(cl, nil)
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetFilmsByMonthYear(w, r)
		result := `{"body":` + test.resultJSON[:len(test.resultJSON)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetMYFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/calendar?"
	for i, test := range testTableGetMYFailure {
		mock := mock2.NewMockFilmUsecase(ctrl)
		if i == 2 {
			mock.EXPECT().GetFilmsByMonthYear(2, 2010, 1, 14).Return(domain.FilmList{}, customErrors.ErrorSkip)
		}
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetFilmsByMonthYear(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTableGetBMSuccess = [...]testRow{
	{
		inQuery:    "id=8&skips=0&limits=10",
		out:        `{"bookmarks_list":null,"more_available":false,"films_total":0,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		resultJSON: `{"bookmarks_list":null,"more_available":false,"films_total":0,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		status:     http.StatusOK,
		name:       `full works`,
		skip:       0,
		limit:      10,
	},
	{
		inQuery:    "id=8",
		out:        `{"bookmarks_list":null,"more_available":false,"films_total":0,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		resultJSON: `{"bookmarks_list":null,"more_available":false,"films_total":0,"current_sort":"","current_limit":10,"current_skip":10}` + "\n",
		status:     http.StatusOK,
		name:       `empty works`,
		skip:       0,
		limit:      10,
	},
}
var testTableGetBMFailure = [...]testRow{
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

func TestGetBMSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getBookmarks?"
	for _, test := range testTableGetBMSuccess {
		var cl domain.FilmBookmarks
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mock.EXPECT().LoadUserBookmarks(uint64(8), test.skip, test.limit).Return(cl, nil)
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.LoadUserBookmarks(w, r)
		result := `{"body":` + test.resultJSON[:len(test.resultJSON)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetBMFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/loadFilmRecommendations?"
	for i, test := range testTableGetRecomsFailure {
		mock := mock2.NewMockFilmUsecase(ctrl)
		if i == 2 {
			mock.EXPECT().LoadUserBookmarks(uint64(8), test.skip, test.limit).Return(domain.FilmBookmarks{}, customErrors.ErrorSkip)
		}
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.LoadUserBookmarks(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTablePostBMSuccess = [...]testRow{
	{
		inQuery:    "id=2&bookmarked=false",
		out:        `{"film_id":2,"bookmarked":false}` + "\n",
		bodyString: `{"film_id":7,"review_text":"и так тоже","review_type":2}`,
		status:     http.StatusOK,
		name:       `normal`,
	},
}

func TestBookmarkSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/Bookmark?"
	test1 := testRow{
		inQuery:    "",
		bodyString: `{"email": "iva21@mail.ru","password": "123456"}`,
		out:        `{"id":2,"first_name":"Ivan","surname":"Ivanov","email":"iva21@mail.ru","profile_pic":"/pic/1.jpg"}`,
		status:     http.StatusOK,
		name:       "log in user",
	}
	for _, test := range testTablePostBMSuccess {
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
		mockSessions1.EXPECT().CheckSession(r.Context(), &authRPC.Request{ID: 0}).Return(&authRPC.ID{ID: 0}, customErrors.ErrorUserNotLoggedIn).Times(1)
		mockSessions1.EXPECT().StartSession(r.Context(), &authRPC.Request{ID: 2}).Return(&authRPC.Session{Name: "session-name", Value: "aaa"}, nil)
		handler.Login(w, r)

		_ = json.Unmarshal([]byte(test.bodyString), &cl)
		_ = json.Unmarshal([]byte(test.out), &ret)
		mock2 := mock2.NewMockFilmUsecase(ctrl)
		mockSessions2 := mockSessions.NewMockSessionRPCClient(ctrl)
		mock2.EXPECT().BookmarkFilm(uint64(2), uint64(2), false).Times(1).Return(nil)
		handler2 := FilmHandler{FilmUsecase: mock2, AuthClient: mockSessions2}
		r = httptest.NewRequest("POST", apiPath+test.inQuery, bodyReader)
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			r.AddCookie(cookie)
		}
		w = httptest.NewRecorder()
		mockSessions2.EXPECT().CheckSession(r.Context(), &authRPC.Request{Name: "session-name", Value: "aaa"}).Return(&authRPC.ID{ID: 2}, nil).Times(1)
		handler2.BookmarkFilm(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

var testTableGetGenresSuccess = [...]testRow{
	{
		inQuery: "id=2&bookmarked=false",
		out:     `{"genres_list":[{"id":1,"name":"Боевик","picture_url":"/static/media/img/genres/action.svg"},{"id":2,"name":"Драма","picture_url":"/static/media/img/genres/drama.svg"},{"id":3,"name":"Комедия","picture_url":"/static/media/img/genres/comedy.svg"},{"id":4,"name":"Ужасы","picture_url":"/static/media/img/genres/horror.svg"},{"id":5,"name":"Мелодрама","picture_url":"/static/media/img/genres/romantic.svg"},{"id":6,"name":"Триллер","picture_url":"/static/media/img/genres/triller.svg"},{"id":7,"name":"Вестерн","picture_url":"/static/media/img/genres/western.svg"},{"id":8,"name":"Документальный","picture_url":"/static/media/img/genres/documentary.svg"},{"id":9,"name":"Фантастика","picture_url":"/static/media/img/genres/fantasy.svg"},{"id":10,"name":"Приключения","picture_url":"/static/media/img/genres/adventures.svg"},{"id":11,"name":"Семейный","picture_url":"/static/media/img/genres/family.svg"},{"id":12,"name":"Криминал","picture_url":"/static/media/img/genres/criminal.svg"},{"id":13,"name":"Исторический","picture_url":"/static/media/img/genres/historic.svg"},{"id":14,"name":"Аниме","picture_url":"/static/media/img/genres/anime.svg"},{"id":15,"name":"Детектив","picture_url":"/static/media/img/genres/detective.svg"}]}` + "\n",
		status:  http.StatusOK,
		name:    `normal`,
	},
}

func TestGetGenresSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getGenres?"
	for _, test := range testTableGetGenresSuccess {
		var cl domain.GenresList
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mock.EXPECT().GetGenres().Return(cl, nil)
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetGenres(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

var testTableGetBannersSuccess = [...]testRow{
	{
		inQuery: "id=2&bookmarked=false",
		out:     `{"banners_list":[{"id":1,"title":"Фильм «Гарри Поттер и философский камень»","description":"Жизнь десятилетнего Гарри Поттера нельзя назвать сладкой: родители умерли, едва ему исполнился год, а от дяди и тёти, взявших сироту на воспитание, достаются лишь тычки да подзатыльники. Но в одиннадцатый день рождения Гарри всё меняется.","picture_url":"/static/media/img/films/1.webp","link":"/films/1"},{"id":2,"title":"Подборка «Фильмы на вечер»","description":"Если на сегодняшний вечер у вас нет особых планов, а провести его хочется с приятностью, мы подобрали несколько интересных фильмов на самые разные вкусы.","picture_url":"/static/media/img/collections/7.webp","link":"/collections/7"},{"id":3,"title":"Фильм «Фокус»","description":"История об опытном мошеннике, который влюбляется в девушку, делающую первые шаги на поприще нелегального отъема средств у граждан. Отношения становятся для них проблемой, когда обнаруживается, что романтика мешает их нечестному бизнесу.","picture_url":"/static/media/img/films/4.jpg","link":"/films/4"}]}` + "\n",
		status:  http.StatusOK,
		name:    `normal`,
	},
}

func TestGetBannersSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getGenres?"
	for _, test := range testTableGetBannersSuccess {
		var cl domain.BannersList
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mock.EXPECT().GetBanners().Return(cl, nil)
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetBanners(w, r)
		result := `{"body":` + test.out[:len(test.out)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

var testTableGetPopularSuccess = [...]testRow{
	{
		inQuery:    "id=2&bookmarked=false",
		out:        `{"film_list":[{"id":1,"title":"Еще по одной","rating":0.0,"description":"В ресторане собираются учитель истории, психологии, музыки и физрук, чтобы отметить 40-летие одного из них. И решают проверить научную теорию о том, что c самого рождения человек страдает от нехватки алкоголя в крови, а чтобы стать по-настоящему счастливым, нужно быть немного нетрезвым. Друзья договариваются наблюдать, как возлияния скажутся на их работе и личной жизни, и устанавливают правила: не пить вечером и по выходным. Казалось бы, что может пойти не так?","poster_url":"server/images/one-more-drink.webp","release_year":2020,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10}` + "\n",
		resultJSON: `{"film_list":[{"id":1,"title":"Еще по одной","rating":"0.0","description":"В ресторане собираются учитель истории, психологии, музыки и физрук, чтобы отметить 40-летие одного из них. И решают проверить научную теорию о том, что c самого рождения человек страдает от нехватки алкоголя в крови, а чтобы стать по-настоящему счастливым, нужно быть немного нетрезвым. Друзья договариваются наблюдать, как возлияния скажутся на их работе и личной жизни, и устанавливают правила: не пить вечером и по выходным. Казалось бы, что может пойти не так?","poster_url":"server/images/one-more-drink.webp","release_year":2020,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10}` + "\n",
		status:     http.StatusOK,
		name:       `normal`,
	},
}

func TestGetPopularSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getPopular?"
	for _, test := range testTableGetPopularSuccess {
		var cl domain.FilmList
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mock.EXPECT().GetPopularFilms().Return(cl, nil)
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetPopularFilms(w, r)
		result := `{"body":` + test.resultJSON[:len(test.resultJSON)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

var testTableGetGLSuccess = [...]testRow{
	{
		inQuery:    "id=8&skips=0&limits=10",
		out:        `{"id":9,"name":"Фантастика","films":{"film_list":[{"id":1,"title":"Гарри Поттер и философский камень","title_original":"Harry Potter and the Sorcerer's Stone","rating":6.8,"description":"Жизнь десятилетнего Гарри Поттера нельзя назвать сладкой: родители умерли, едва ему исполнился год, а от дяди и тёти, взявших сироту на воспитание, достаются лишь тычки да подзатыльники. Но в одиннадцатый день рождения Гарри всё меняется. Странный гость, неожиданно появившийся на пороге, приносит письмо, из которого мальчик узнаёт, что на самом деле он - волшебник и зачислен в школу магии под названием Хогвартс. А уже через пару недель Гарри будет мчаться в поезде Хогвартс-экспресс навстречу новой жизни, где его ждут невероятные приключения, верные друзья и самое главное — ключ к разгадке тайны смерти его родителей.","poster_url":"/static/media/img/films/1.webp","release_year":2001,"premiere_ru":"2021-11-23","director":{},"screenwriter":{}},{"id":2,"title":"Гарри Поттер и Тайная комната","title_original":"Harry Potter and the Chamber of Secrets","rating":3.5,"description":"Гарри Поттер переходит на второй курс Школы чародейства и волшебства Хогвартс. Эльф Добби предупреждает Гарри об опасности, которая поджидает его там, и просит больше не возвращаться в школу. Юный волшебник не следует совету эльфа и становится свидетелем таинственных событий, разворачивающихся в Хогвартсе. Вскоре Гарри и его друзья узнают о существовании Тайной Комнаты и сталкиваются с новыми приключениями, пытаясь победить темные силы.","poster_url":"/static/media/img/films/2.webp","release_year":2002,"premiere_ru":"2022-01-22","director":{},"screenwriter":{}},{"id":3,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":2.0,"description":"В третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер, Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","poster_url":"/static/media/img/films/3.webp","release_year":2004,"premiere_ru":"2021-11-02","director":{},"screenwriter":{}}],"more_available":false,"film_total":3,"current_limit":5,"current_skip":5}}` + "\n",
		resultJSON: `{"id":9,"name":"Фантастика","films":{"film_list":[{"id":1,"title":"Гарри Поттер и философский камень","title_original":"Harry Potter and the Sorcerer's Stone","rating":"6.8","description":"Жизнь десятилетнего Гарри Поттера нельзя назвать сладкой: родители умерли, едва ему исполнился год, а от дяди и тёти, взявших сироту на воспитание, достаются лишь тычки да подзатыльники. Но в одиннадцатый день рождения Гарри всё меняется. Странный гость, неожиданно появившийся на пороге, приносит письмо, из которого мальчик узнаёт, что на самом деле он - волшебник и зачислен в школу магии под названием Хогвартс. А уже через пару недель Гарри будет мчаться в поезде Хогвартс-экспресс навстречу новой жизни, где его ждут невероятные приключения, верные друзья и самое главное — ключ к разгадке тайны смерти его родителей.","poster_url":"/static/media/img/films/1.webp","release_year":2001,"premiere_ru":"2021-11-23","director":{},"screenwriter":{}},{"id":2,"title":"Гарри Поттер и Тайная комната","title_original":"Harry Potter and the Chamber of Secrets","rating":"3.5","description":"Гарри Поттер переходит на второй курс Школы чародейства и волшебства Хогвартс. Эльф Добби предупреждает Гарри об опасности, которая поджидает его там, и просит больше не возвращаться в школу. Юный волшебник не следует совету эльфа и становится свидетелем таинственных событий, разворачивающихся в Хогвартсе. Вскоре Гарри и его друзья узнают о существовании Тайной Комнаты и сталкиваются с новыми приключениями, пытаясь победить темные силы.","poster_url":"/static/media/img/films/2.webp","release_year":2002,"premiere_ru":"2022-01-22","director":{},"screenwriter":{}},{"id":3,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":"2.0","description":"В третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер, Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","poster_url":"/static/media/img/films/3.webp","release_year":2004,"premiere_ru":"2021-11-02","director":{},"screenwriter":{}}],"more_available":false,"film_total":3,"current_limit":5,"current_skip":5}}` + "\n",
		status:     http.StatusOK,
		name:       `full works`,
		skip:       0,
		limit:      10,
	},
}
var testTableGetGLFailure = [...]testRow{
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

func TestGetGLSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getGenres?"
	for _, test := range testTableGetGLSuccess {
		var cl domain.GenreFilmList
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockFilmUsecase(ctrl)
		mock.EXPECT().GetFilmsByGenre(uint64(8), test.limit, test.skip).Return(cl, nil)
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetFilmsByGenre(w, r)
		result := `{"body":` + test.resultJSON[:len(test.resultJSON)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetGLFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getGenres?"
	for i, test := range testTableGetGLFailure {
		mock := mock2.NewMockFilmUsecase(ctrl)
		if i == 2 {
			mock.EXPECT().GetFilmsByGenre(uint64(8), test.limit, test.skip).Return(domain.GenreFilmList{}, customErrors.ErrorSkip)
		}
		handler := FilmHandler{FilmUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetFilmsByGenre(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}
