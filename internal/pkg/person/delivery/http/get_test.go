package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	mock2 "2021_2_MAMBa/internal/pkg/person/usecase/mock"
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

var testTableGetPersonSuccess = [...]testRow{
	{
		inQuery: "id=9",
		out:     `{"actor":{"id":9,"name_rus":"Сэм Клафлин","career":[""],"height":180,"age":20,"birthday":"0001-01-01 00:00:00 +0000 UTC","death":"0001-01-01 00:00:00 +0000 UTC","gender":"male"},"films":{"film_list":[{"id":2,"title":"С любовью, Рози","rating":0.0,"description":"Рози и Алекс были лучшими друзьями с детства, и теперь, по окончании школы, собираются вместе продолжить учёбу в университете.\\n\\nОднако в их судьбах происходит резкий поворот, когда после ночи со звездой школы Рози узнаёт, что у неё будет ребенок. Невзирая на то, что обстоятельства и жизненные ситуации разлучают героев, они и спустя годы продолжают помнить друг о друге и о том чувстве, что соединило их в юности…","poster_url":"server/images/love-rosie.webp","release_year":2014,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10},"popular_films":{"film_list":[{"id":2,"title":"С любовью, Рози","rating":0.0,"description":"Рози и Алекс были лучшими друзьями с детства, и теперь, по окончании школы, собираются вместе продолжить учёбу в университете.\\n\\nОднако в их судьбах происходит резкий поворот, когда после ночи со звездой школы Рози узнаёт, что у неё будет ребенок. Невзирая на то, что обстоятельства и жизненные ситуации разлучают героев, они и спустя годы продолжают помнить друг о друге и о том чувстве, что соединило их в юности…","poster_url":"server/images/love-rosie.webp","release_year":2014,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10}}` + "\n",
		status:  http.StatusOK,
		name:    `full works`,
		skip:    0,
		limit:   10,
		skip1:   0,
		limit1:  10,
	},
}
var testTableGetPersonFailure = [...]testRow{
	{
		inQuery: "",
		out:     customErrors.ErrIdMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `no id`,
	},
	{
		inQuery: "id=10",
		out:     customErrors.ErrDBMsg + "\n",
		status:  http.StatusInternalServerError,
		name:    `overshoot`,
	},
}

func TestGetPersonSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/person/getPerson?"
	for _, test := range testTableGetPersonSuccess {
		var cl domain.PersonPage
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockPersonUsecase(ctrl)
		mock.EXPECT().GetPerson(uint64(9)).Times(1).Return(cl, nil)
		handler := PersonHandler{PersonUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetPerson(w, r)
		result:= `{"body":`+test.out[:len(test.out)-1]+`,"status":`+fmt.Sprint(test.status)+"}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetPersonFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/person/getPerson?"
	for i, test := range testTableGetPersonFailure {
		mock := mock2.NewMockPersonUsecase(ctrl)
		if i == 0 {
			mock.EXPECT().GetPerson(uint64(0)).Times(1).Return(domain.PersonPage{}, customErrors.ErrorBadInput)
		}
		if i == 1 {
			mock.EXPECT().GetPerson(uint64(10)).Times(1).Return(domain.PersonPage{}, customErrors.ErrorInternalServer)
		}
		handler := PersonHandler{PersonUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetPerson(w, r)
		result:= `{"body":{"error":"`+test.out[:len(test.out)-1]+`"},"status":`+fmt.Sprint(test.status)+"}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}

var testTableGetPersonFilmsSuccess = [...]testRow{
	{
		inQuery: "id=5&skips=0&limits=10",
		out:     `{"film_list":[{"id":1,"title":"Еще по одной","rating":0.0,"description":"В ресторане собираются учитель истории, психологии, музыки и физрук, чтобы отметить 40-летие одного из них. И решают проверить научную теорию о том, что c самого рождения человек страдает от нехватки алкоголя в крови, а чтобы стать по-настоящему счастливым, нужно быть немного нетрезвым. Друзья договариваются наблюдать, как возлияния скажутся на их работе и личной жизни, и устанавливают правила: не пить вечером и по выходным. Казалось бы, что может пойти не так?","poster_url":"server/images/one-more-drink.webp","release_year":2020,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `full works`,
		skip:    0,
		limit:   10,
	},
	{
		inQuery: "id=5",
		out:     `{"film_list":[{"id":1,"title":"Еще по одной","rating":0.0,"description":"В ресторане собираются учитель истории, психологии, музыки и физрук, чтобы отметить 40-летие одного из них. И решают проверить научную теорию о том, что c самого рождения человек страдает от нехватки алкоголя в крови, а чтобы стать по-настоящему счастливым, нужно быть немного нетрезвым. Друзья договариваются наблюдать, как возлияния скажутся на их работе и личной жизни, и устанавливают правила: не пить вечером и по выходным. Казалось бы, что может пойти не так?","poster_url":"server/images/one-more-drink.webp","release_year":2020,"director":{},"screenwriter":{}}],"more_available":false,"film_total":1,"current_limit":10,"current_skip":10}` + "\n",
		status:  http.StatusOK,
		name:    `empty works`,
		skip:    0,
		limit:   10,
	},
}
var testTableGetPersonFilmsFailure = [...]testRow{
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
	apiPath := "/api/person/getPersonFilms?"
	for _, test := range testTableGetPersonFilmsSuccess {
		var cl domain.FilmList
		_ = json.Unmarshal([]byte(test.out), &cl)
		mock := mock2.NewMockPersonUsecase(ctrl)
		mock.EXPECT().GetFilms(uint64(5), test.skip, test.limit).Return(cl, nil)
		handler := PersonHandler{PersonUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetPersonFilms(w, r)
		result:= `{"body":`+test.out[:len(test.out)-1]+`,"status":`+fmt.Sprint(test.status)+"}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetRecomFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/getPersonFilms?"
	for i, test := range testTableGetPersonFilmsFailure {
		mock := mock2.NewMockPersonUsecase(ctrl)
		if i == 2 {
			mock.EXPECT().GetFilms(uint64(8), test.skip, test.limit).Return(domain.FilmList{}, customErrors.ErrorSkip)
		}
		handler := PersonHandler{PersonUsecase: mock}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetPersonFilms(w, r)
		result:= `{"body":{"error":"`+test.out[:len(test.out)-1]+`"},"status":`+fmt.Sprint(test.status)+"}\n"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}
