package http

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	mock2 "2021_2_MAMBa/internal/pkg/search/usecase/mock"
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
	resultJSON string
	status     int
	name       string
	skip       int
	limit      int
	skip1      int
	limit1     int
}

var testTableGetMYSuccess = [...]testRow{
	{
		inQuery:    "query=%D1%80%D0%B8",
		out:        `{"films":{"film_list":[{"id":1,"title":"Гарри Поттер и философский камень","title_original":"Harry Potter and the Sorcerer's Stone","rating":6.8,"description":"Жизнь десятилетнего Гарри Поттера нельзя назвать сладкой: родители умерли, едва ему исполнился год, а от дяди и тёти, взявших сироту на воспитание, достаются лишь тычки да подзатыльники. Но в одиннадцатый день рождения Гарри всё меняется. Странный гость, неожиданно появившийся на пороге, приносит письмо, из которого мальчик узнаёт, что на самом деле он - волшебник и зачислен в школу магии под названием Хогвартс. А уже через пару недель Гарри будет мчаться в поезде Хогвартс-экспресс навстречу новой жизни, где его ждут невероятные приключения, верные друзья и самое главное — ключ к разгадке тайны смерти его родителей.","total_revenue":"$0.00","poster_url":"/static/media/img/films/1.webp","trailer_url":"https://www.youtube.com/watch?v=ppnEwu-z9kU","content_type":"film","release_year":2001,"duration":152,"premiere_ru":"2021-11-23","origin_countries":["Великобритания","США"],"cast":[{"id":3,"name_en":"Daniel Radcliffe","name_rus":"Дэниэл Рэдклифф","picture_url":"/static/media/img/persons/3.jpg","career":["Актер, Продюсер"]},{"id":4,"name_en":"Rupert Grint","name_rus":"Руперт Гринт","picture_url":"/static/media/img/persons/4.jpg","career":["Актер, Продюсер"]},{"id":5,"name_en":"Emma Watson","name_rus":"Эмма Уотсон","picture_url":"/static/media/img/persons/5.jpg","career":["Актриса"]}],"director":{"id":2,"name_en":"Chris Columbus","name_rus":"Крис Коламбус","picture_url":"/static/media/img/persons/2.jpg","career":["Продюсер, Режиссер, Сценарист, Актер"]},"screenwriter":{"id":1,"name_en":"Steven Kloves","name_rus":"Стивен Кловз","picture_url":"/static/media/img/persons/1.webp","career":["Сценарист, Продюсер, Режиссер"]},"genres":[{"id":9,"name":"Фантастика"},{"id":10,"name":"Приключения"},{"id":11,"name":"Семейный"}]},{"id":3,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":2.0,"description":"В третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер, Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","total_revenue":"$0.00","poster_url":"/static/media/img/films/3.webp","trailer_url":"https://www.youtube.com/watch?v=ofHrcJFd8hA","content_type":"film","release_year":2004,"duration":142,"premiere_ru":"2021-11-02","origin_countries":["Великобритания","США"],"cast":[{"id":3,"name_en":"Daniel Radcliffe","name_rus":"Дэниэл Рэдклифф","picture_url":"/static/media/img/persons/3.jpg","career":["Актер, Продюсер"]},{"id":4,"name_en":"Rupert Grint","name_rus":"Руперт Гринт","picture_url":"/static/media/img/persons/4.jpg","career":["Актер, Продюсер"]},{"id":5,"name_en":"Emma Watson","name_rus":"Эмма Уотсон","picture_url":"/static/media/img/persons/5.jpg","career":["Актриса"]}],"director":{"id":6,"name_en":"Alfonso Cuarón","name_rus":"Альфонсо Куарон","picture_url":"/static/media/img/persons/6.webp","career":["Продюсер, Режиссер, Сценарист, Оператор, Монтажер, Актер"]},"screenwriter":{"id":1,"name_en":"Steven Kloves","name_rus":"Стивен Кловз","picture_url":"/static/media/img/persons/1.webp","career":["Сценарист, Продюсер, Режиссер"]},"genres":[{"id":9,"name":"Фантастика"},{"id":10,"name":"Приключения"},{"id":11,"name":"Семейный"}]},{"id":2,"title":"Гарри Поттер и Тайная комната","title_original":"Harry Potter and the Chamber of Secrets","rating":3.5,"description":"Гарри Поттер переходит на второй курс Школы чародейства и волшебства Хогвартс. Эльф Добби предупреждает Гарри об опасности, которая поджидает его там, и просит больше не возвращаться в школу. Юный волшебник не следует совету эльфа и становится свидетелем таинственных событий, разворачивающихся в Хогвартсе. Вскоре Гарри и его друзья узнают о существовании Тайной Комнаты и сталкиваются с новыми приключениями, пытаясь победить темные силы.","total_revenue":"$0.00","poster_url":"/static/media/img/films/2.webp","trailer_url":"https://www.youtube.com/watch?v=fvfHr2SB1mg","content_type":"film","release_year":2002,"duration":161,"premiere_ru":"2022-01-22","origin_countries":["Великобритания","США","Германия"],"cast":[{"id":3,"name_en":"Daniel Radcliffe","name_rus":"Дэниэл Рэдклифф","picture_url":"/static/media/img/persons/3.jpg","career":["Актер, Продюсер"]},{"id":4,"name_en":"Rupert Grint","name_rus":"Руперт Гринт","picture_url":"/static/media/img/persons/4.jpg","career":["Актер, Продюсер"]},{"id":5,"name_en":"Emma Watson","name_rus":"Эмма Уотсон","picture_url":"/static/media/img/persons/5.jpg","career":["Актриса"]}],"director":{"id":2,"name_en":"Chris Columbus","name_rus":"Крис Коламбус","picture_url":"/static/media/img/persons/2.jpg","career":["Продюсер, Режиссер, Сценарист, Актер"]},"screenwriter":{"id":1,"name_en":"Steven Kloves","name_rus":"Стивен Кловз","picture_url":"/static/media/img/persons/1.webp","career":["Сценарист, Продюсер, Режиссер"]},"genres":[{"id":9,"name":"Фантастика"},{"id":10,"name":"Приключения"},{"id":11,"name":"Семейный"}]}],"more_available":true,"film_total":3,"current_limit":10,"current_skip":10},"persons":{"person_list":[{"id":2,"name_en":"Chris Columbus","name_rus":"Крис Коламбус","picture_url":"/static/media/img/persons/2.jpg","career":["Продюсер"," Режиссер"," Сценарист"," Актер"],"height":174,"age":63,"birthday":"10.09.1958","birth_place":"Спанглер, Пенсильвания, США","gender":"male","family_status":"Женат","film_number":3},{"id":4,"name_en":"Rupert Grint","name_rus":"Руперт Гринт","picture_url":"/static/media/img/persons/4.jpg","career":["Актер"," Продюсер"],"height":173,"age":33,"birthday":"24.08.1988","birth_place":"Стивенэйдж, Хартфордшир, Англия, Великобритания","gender":"male","family_status":"Холост","film_number":3},{"id":11,"name_en":"Christian Ditter","name_rus":"Кристиан Диттер","picture_url":"/static/media/img/persons/11.jpg","career":["Режиссер"," Сценарист"," Продюсер"," Монтажер"],"height":180,"age":44,"birthday":"01.01.1977","birth_place":"Гисен, Германия","gender":"male","family_status":"Холост","film_number":1},{"id":14,"name_en":"Adrian Martinez","name_rus":"Адриан Мартинес","picture_url":"/static/media/img/persons/14.jpg","career":["Актер"," Продюсер"," Режиссер"," Сценарист"],"height":173,"age":49,"birthday":"20.01.1972","birth_place":"Нью-Йорк, США","gender":"male","family_status":"Холост","film_number":1},{"id":16,"name_en":"Chris Pratt","name_rus":"Крис Пратт","picture_url":"/static/media/img/persons/16.jpg","career":["Актер"," Продюсер"],"height":188,"age":42,"birthday":"21.06.1979","birth_place":"Вирджиния, Миннесота, США","gender":"male","family_status":"Женат","film_number":1},{"id":19,"name_en":"Alison Brie","name_rus":"Элисон Бри","picture_url":"/static/media/img/persons/19.jpg","career":["Актриса"," Продюсер"," Сценарист"," Режиссер"],"height":160,"age":38,"birthday":"29.12.1982","birth_place":"Лос-Анджелес, Калифорния, США","gender":"female","family_status":"Замужем","film_number":1},{"id":20,"name_en":"Eric Roth","name_rus":"Эрик Рот","picture_url":"/static/media/img/persons/20.jpg","career":["Сценарист"," Продюсер"," Актер"],"height":178,"age":76,"birthday":"22.03.1945","birth_place":"Нью-Йорк, США","gender":"male","family_status":"Холост","film_number":1}],"more_available":true,"person_total":7,"current_limit":10,"current_skip":10}}` + "\n",
		resultJSON: `{"films":{"film_list":[{"id":1,"title":"Гарри Поттер и философский камень","title_original":"Harry Potter and the Sorcerer's Stone","rating":"6.8","description":"Жизнь десятилетнего Гарри Поттера нельзя назвать сладкой: родители умерли, едва ему исполнился год, а от дяди и тёти, взявших сироту на воспитание, достаются лишь тычки да подзатыльники. Но в одиннадцатый день рождения Гарри всё меняется. Странный гость, неожиданно появившийся на пороге, приносит письмо, из которого мальчик узнаёт, что на самом деле он - волшебник и зачислен в школу магии под названием Хогвартс. А уже через пару недель Гарри будет мчаться в поезде Хогвартс-экспресс навстречу новой жизни, где его ждут невероятные приключения, верные друзья и самое главное — ключ к разгадке тайны смерти его родителей.","total_revenue":"$0.00","poster_url":"/static/media/img/films/1.webp","trailer_url":"https://www.youtube.com/watch?v=ppnEwu-z9kU","content_type":"film","release_year":2001,"duration":152,"premiere_ru":"2021-11-23","origin_countries":["Великобритания","США"],"cast":[{"id":3,"name_en":"Daniel Radcliffe","name_rus":"Дэниэл Рэдклифф","picture_url":"/static/media/img/persons/3.jpg","career":["Актер, Продюсер"]},{"id":4,"name_en":"Rupert Grint","name_rus":"Руперт Гринт","picture_url":"/static/media/img/persons/4.jpg","career":["Актер, Продюсер"]},{"id":5,"name_en":"Emma Watson","name_rus":"Эмма Уотсон","picture_url":"/static/media/img/persons/5.jpg","career":["Актриса"]}],"director":{"id":2,"name_en":"Chris Columbus","name_rus":"Крис Коламбус","picture_url":"/static/media/img/persons/2.jpg","career":["Продюсер, Режиссер, Сценарист, Актер"]},"screenwriter":{"id":1,"name_en":"Steven Kloves","name_rus":"Стивен Кловз","picture_url":"/static/media/img/persons/1.webp","career":["Сценарист, Продюсер, Режиссер"]},"genres":[{"id":9,"name":"Фантастика"},{"id":10,"name":"Приключения"},{"id":11,"name":"Семейный"}]},{"id":3,"title":"Гарри Поттер и узник Азкабана","title_original":"Harry Potter and the Prisoner of Azkaban","rating":"2.0","description":"В третьей части истории о юном волшебнике полюбившиеся всем герои — Гарри Поттер, Рон и Гермиона — возвращаются уже на третий курс школы чародейства и волшебства Хогвартс. На этот раз они должны раскрыть тайну узника, сбежавшего из зловещей тюрьмы Азкабан, чье пребывание на воле создает для Гарри смертельную опасность...","total_revenue":"$0.00","poster_url":"/static/media/img/films/3.webp","trailer_url":"https://www.youtube.com/watch?v=ofHrcJFd8hA","content_type":"film","release_year":2004,"duration":142,"premiere_ru":"2021-11-02","origin_countries":["Великобритания","США"],"cast":[{"id":3,"name_en":"Daniel Radcliffe","name_rus":"Дэниэл Рэдклифф","picture_url":"/static/media/img/persons/3.jpg","career":["Актер, Продюсер"]},{"id":4,"name_en":"Rupert Grint","name_rus":"Руперт Гринт","picture_url":"/static/media/img/persons/4.jpg","career":["Актер, Продюсер"]},{"id":5,"name_en":"Emma Watson","name_rus":"Эмма Уотсон","picture_url":"/static/media/img/persons/5.jpg","career":["Актриса"]}],"director":{"id":6,"name_en":"Alfonso Cuarón","name_rus":"Альфонсо Куарон","picture_url":"/static/media/img/persons/6.webp","career":["Продюсер, Режиссер, Сценарист, Оператор, Монтажер, Актер"]},"screenwriter":{"id":1,"name_en":"Steven Kloves","name_rus":"Стивен Кловз","picture_url":"/static/media/img/persons/1.webp","career":["Сценарист, Продюсер, Режиссер"]},"genres":[{"id":9,"name":"Фантастика"},{"id":10,"name":"Приключения"},{"id":11,"name":"Семейный"}]},{"id":2,"title":"Гарри Поттер и Тайная комната","title_original":"Harry Potter and the Chamber of Secrets","rating":"3.5","description":"Гарри Поттер переходит на второй курс Школы чародейства и волшебства Хогвартс. Эльф Добби предупреждает Гарри об опасности, которая поджидает его там, и просит больше не возвращаться в школу. Юный волшебник не следует совету эльфа и становится свидетелем таинственных событий, разворачивающихся в Хогвартсе. Вскоре Гарри и его друзья узнают о существовании Тайной Комнаты и сталкиваются с новыми приключениями, пытаясь победить темные силы.","total_revenue":"$0.00","poster_url":"/static/media/img/films/2.webp","trailer_url":"https://www.youtube.com/watch?v=fvfHr2SB1mg","content_type":"film","release_year":2002,"duration":161,"premiere_ru":"2022-01-22","origin_countries":["Великобритания","США","Германия"],"cast":[{"id":3,"name_en":"Daniel Radcliffe","name_rus":"Дэниэл Рэдклифф","picture_url":"/static/media/img/persons/3.jpg","career":["Актер, Продюсер"]},{"id":4,"name_en":"Rupert Grint","name_rus":"Руперт Гринт","picture_url":"/static/media/img/persons/4.jpg","career":["Актер, Продюсер"]},{"id":5,"name_en":"Emma Watson","name_rus":"Эмма Уотсон","picture_url":"/static/media/img/persons/5.jpg","career":["Актриса"]}],"director":{"id":2,"name_en":"Chris Columbus","name_rus":"Крис Коламбус","picture_url":"/static/media/img/persons/2.jpg","career":["Продюсер, Режиссер, Сценарист, Актер"]},"screenwriter":{"id":1,"name_en":"Steven Kloves","name_rus":"Стивен Кловз","picture_url":"/static/media/img/persons/1.webp","career":["Сценарист, Продюсер, Режиссер"]},"genres":[{"id":9,"name":"Фантастика"},{"id":10,"name":"Приключения"},{"id":11,"name":"Семейный"}]}],"more_available":true,"film_total":3,"current_limit":10,"current_skip":10},"persons":{"person_list":[{"id":2,"name_en":"Chris Columbus","name_rus":"Крис Коламбус","picture_url":"/static/media/img/persons/2.jpg","career":["Продюсер"," Режиссер"," Сценарист"," Актер"],"height":174,"age":63,"birthday":"10.09.1958","birth_place":"Спанглер, Пенсильвания, США","gender":"male","family_status":"Женат","film_number":3},{"id":4,"name_en":"Rupert Grint","name_rus":"Руперт Гринт","picture_url":"/static/media/img/persons/4.jpg","career":["Актер"," Продюсер"],"height":173,"age":33,"birthday":"24.08.1988","birth_place":"Стивенэйдж, Хартфордшир, Англия, Великобритания","gender":"male","family_status":"Холост","film_number":3},{"id":11,"name_en":"Christian Ditter","name_rus":"Кристиан Диттер","picture_url":"/static/media/img/persons/11.jpg","career":["Режиссер"," Сценарист"," Продюсер"," Монтажер"],"height":180,"age":44,"birthday":"01.01.1977","birth_place":"Гисен, Германия","gender":"male","family_status":"Холост","film_number":1},{"id":14,"name_en":"Adrian Martinez","name_rus":"Адриан Мартинес","picture_url":"/static/media/img/persons/14.jpg","career":["Актер"," Продюсер"," Режиссер"," Сценарист"],"height":173,"age":49,"birthday":"20.01.1972","birth_place":"Нью-Йорк, США","gender":"male","family_status":"Холост","film_number":1},{"id":16,"name_en":"Chris Pratt","name_rus":"Крис Пратт","picture_url":"/static/media/img/persons/16.jpg","career":["Актер"," Продюсер"],"height":188,"age":42,"birthday":"21.06.1979","birth_place":"Вирджиния, Миннесота, США","gender":"male","family_status":"Женат","film_number":1},{"id":19,"name_en":"Alison Brie","name_rus":"Элисон Бри","picture_url":"/static/media/img/persons/19.jpg","career":["Актриса"," Продюсер"," Сценарист"," Режиссер"],"height":160,"age":38,"birthday":"29.12.1982","birth_place":"Лос-Анджелес, Калифорния, США","gender":"female","family_status":"Замужем","film_number":1},{"id":20,"name_en":"Eric Roth","name_rus":"Эрик Рот","picture_url":"/static/media/img/persons/20.jpg","career":["Сценарист"," Продюсер"," Актер"],"height":178,"age":76,"birthday":"22.03.1945","birth_place":"Нью-Йорк, США","gender":"male","family_status":"Холост","film_number":1}],"more_available":true,"person_total":7,"current_limit":10,"current_skip":10}}` + "\n",
		status:     http.StatusOK,
		name:       `full works`,
		skip:       0,
		limit:      10,
	},
}

var testTableGetMYFailure = [...]testRow{
	{
		inQuery: "query=\"a\"&skip_films=-1&limit_films=10&skip_persons=10&limit_persons=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    -1,
		limit:   10,
		skip1:   10,
		limit1:  10,
	},
	{
		inQuery: "query=\"a\"&skip_films=10&limit_films=-2&skip_persons=10&limit_persons=10",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    10,
		limit:   -2,
		skip1:   10,
		limit1:  10,
	},
	{
		inQuery: "query=\"a\"&skip_films=10&limit_films=10&skip_persons=-1&limit_persons=10",
		out:     customErrors.ErrSkipMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative skip`,
		skip:    10,
		limit:   10,
		skip1:   -1,
		limit1:  10,
	},
	{
		inQuery: "query=\"a\"&skip_films=10&limit_films=-2&skip_persons=10&limit_persons=-1",
		out:     customErrors.ErrLimitMsg + "\n",
		status:  http.StatusBadRequest,
		name:    `negative limit`,
		skip:    10,
		limit:   10,
		skip1:   10,
		limit1:  -1,
	},
}

func TestGetMYSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/calendar?"
	for _, test := range testTableGetMYSuccess {
		var cl domain.SearchResult
		_ = json.Unmarshal([]byte(test.out[:len(test.out)-1]), &cl)
		mock := mock2.NewMockSearchUsecase(ctrl)
		mock.EXPECT().GetSearch("ри", 0, 10, 0, 10).Return(cl, nil)
		handler := SearchHandler{
			SearchUsecase: mock,
		}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetSearch(w, r)
		result := `{"body":` + test.resultJSON[:len(test.resultJSON)-1] + `,"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
		assert.Equal(t, test.status, w.Code, "Test: "+test.name)
	}
}

func TestGetMYFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	apiPath := "/api/film/calendar?"
	for _, test := range testTableGetMYFailure {
		mock := mock2.NewMockSearchUsecase(ctrl)
		handler := SearchHandler{
			SearchUsecase: mock,
		}
		bodyReader := strings.NewReader("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", apiPath+test.inQuery, bodyReader)
		handler.GetSearch(w, r)
		result := `{"body":{"error":"` + test.out[:len(test.out)-1] + `"},"status":` + fmt.Sprint(test.status) + "}"
		assert.Equal(t, result, w.Body.String(), "Test: "+test.name)
	}
}
