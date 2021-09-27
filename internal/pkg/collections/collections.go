package collections

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type collectionPreview struct {
	Id          uint  `json:"id"`
	Title      string `json:"title"`
	PictureUrl string `json:"picture_url"`
}

type collections struct {
	CollArray       []collectionPreview `json:"collections_list"`
	MoreAvailable   bool                `json:"more_available"`
	CollectionTotal int                 `json:"collection_total"`
	CurrentSort     string `json:"current_sort"`
	CurrentLimit    int `json:"current_limit"`
	CurrentSkip     int `json:"current_skip"`
}

var (
	errSkip = `incorrect skip`
	errlimit = `incorrect limit`
	errDB = `DB error`
	errEnc =`Encoding error`
	errorsSkip = errors.New(errSkip)
)

var previewMock = []collectionPreview {
	{Id: 1, Title: "Для ценителей Хогвардса", PictureUrl: "server/images/collections1.png"},
	{Id: 2, Title: "Про настоящую любовь", PictureUrl: "server/images/collections2.png"},
	{Id: 3, Title: "Аферы века", PictureUrl: "server/images/collections3.png"},
	{Id: 4, Title: "Про Вторую Мировую", PictureUrl: "server/images/collections4.jpg"},
	{Id: 5, Title: "Осеннее настроение", PictureUrl: "server/images/collections5.png"},
	{Id: 6, Title: "Летняя атмосфера", PictureUrl: "server/images/collections6.png"},
	{Id: 7, Title: "Про дружбу", PictureUrl: "server/images/collections7.png"},
	{Id: 8, Title: "Романтические фильмы", PictureUrl: "server/images/collections8.jpg"},
	{Id: 9, Title: "Джунгди зовут", PictureUrl: "server/images/collections9.jpg"},
	{Id: 10, Title: "Фантастические фильмы", PictureUrl: "server/images/collections10.jpg"},
	{Id: 11, Title: "Про петлю времени", PictureUrl: "server/images/collections11.png"},
	{Id: 12, Title: "Классика на века", PictureUrl: "server/images/collections12.jpg"},
}

// separating handler from DB methods
func getCollectionsDB (skip int, limit int) (collections, error) {
	moreAvailable := skip+limit < len(previewMock)
	next := skip+limit
	if !moreAvailable {
		next = len(previewMock)
	}
	if skip >= len(previewMock) {
		return collections{}, errorsSkip
	}
	collect := collections {
		CollArray:       previewMock[skip:next],
		MoreAvailable:   moreAvailable,
		CollectionTotal: len(previewMock),
		CurrentLimit:    limit,
		CurrentSkip:     skip+limit,
	}
	return collect, nil
}

func GetCollections(w http.ResponseWriter, r *http.Request) {
	skipString, isIn := r.URL.Query()["skip"]
	var err error
	//default
	limit, skip := 10,0
	if isIn {
		skip, err = strconv.Atoi(skipString[0])
		if err != nil || skip < 0 {
			http.Error(w, errSkip, http.StatusBadRequest)
			return
		}
	}
	limitString, isIn := r.URL.Query()["limit"]
	if isIn {
		limit, err = strconv.Atoi(limitString[0])
		if err != nil || limit <= 0 {
			http.Error(w, errlimit, http.StatusBadRequest)
			return
		}
	}

	collectionsList, err := getCollectionsDB(skip, limit)
	if err == errorsSkip {
		http.Error(w, errSkip, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, errDB, http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(collectionsList)
	if err != nil {
		http.Error(w, errEnc, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}