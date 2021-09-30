package collections

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"2021_2_MAMBa/internal/pkg/database"
)

type collections struct {
	CollArray       []database.CollectionPreview `json:"collections_list"`
	MoreAvailable   bool                         `json:"more_available"`
	CollectionTotal int                          `json:"collection_total"`
	CurrentSort     string                       `json:"current_sort"`
	CurrentLimit    int                          `json:"current_limit"`
	CurrentSkip     int                          `json:"current_skip"`
}

var (
	errSkipMsg  = "incorrect skip"
	errLimitMsg = "incorrect limit"
	errDBMsg    = "DB error"
	errEncMsg   = "Encoding error"
	errorSkip   = errors.New(errSkipMsg)
	errorLimit  = errors.New(errLimitMsg)
	db          = database.CollectionsMockDatabase{Previews: database.PreviewMock}
)

// БД и хэндлер отдельно
func getCollectionsDB(skip int, limit int) (collections, error) {
	db.RLock()
	dbSize := len(db.Previews)
	db.RUnlock()
	if skip >= dbSize {
		return collections{}, errorSkip
	}
	moreAvailable := skip+limit < dbSize
	next := skip + limit
	if !moreAvailable {
		next = dbSize
	}
	db.RLock()
	previews := db.Previews[skip:next]
	db.RUnlock()

	collect := collections{
		CollArray:       previews,
		MoreAvailable:   moreAvailable,
		CollectionTotal: dbSize,
		CurrentLimit:    limit,
		CurrentSkip:     skip + limit,
	}
	return collect, nil
}

func GetCollections(w http.ResponseWriter, r *http.Request) {
	var err error
	// default
	limit, skip := 10, 0
	skipString, isIn := r.URL.Query()["skip"]
	if isIn {
		skip, err = strconv.Atoi(skipString[0])
		if err != nil || skip < 0 {
			http.Error(w, errSkipMsg, http.StatusBadRequest)
			return
		}
	}
	limitString, isIn := r.URL.Query()["limit"]
	if isIn {
		limit, err = strconv.Atoi(limitString[0])
		if err != nil || limit <= 0 {
			http.Error(w, errLimitMsg, http.StatusBadRequest)
			return
		}
	}
	collectionsList, err := getCollectionsDB(skip, limit)
	if err == errorSkip {
		http.Error(w, errSkipMsg, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, errDBMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(collectionsList)
	if err != nil {
		http.Error(w, errEncMsg, http.StatusInternalServerError)
		return
	}
}
