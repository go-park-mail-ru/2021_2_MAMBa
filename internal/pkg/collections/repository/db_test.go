package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	mylog "2021_2_MAMBa/internal/pkg/utils/log"
	"encoding/binary"
	"errors"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"math"
	"regexp"
	"testing"
)

type testRow struct {
	inQuery    string
	bodyString string
	out        string
	status     int
	name       string
}

func MockDatabase() (*database.DBManager, pgxmock.PgxPoolIface, error) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		mylog.Error(errors.New("failed to create mock"))
	}
	return &database.DBManager{Pool: mock}, mock, err
}

var mockPersonPreview = domain.Person{
	Id:         1,
	NameEn:     "Miley",
	NameRus:    "Сайрус",
	PictureUrl: "/miley.webp",
	Career:     []string{"Актриса"},
}

func TestGetSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewCollectionsRepository(mdb)
	defer pool.Close()
	mockCollections := domain.Collections{
		CollArray: []domain.CollectionPreview{{
			Id:         1,
			Title:      "example",
			PictureUrl: "/111/a.webp",
		}},
		MoreAvailable:   false,
		CollectionTotal: 1,
		CurrentSort:     "",
		CurrentLimit:    10,
		CurrentSkip:     10,
	}

	countByte := make([]uint8, 8)
	binary.BigEndian.PutUint64(countByte, uint64(1))
	idByte := make([]uint8, 8)
	binary.BigEndian.PutUint64(idByte, mockCollections.CollArray[0].Id)
	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(countByte)
	rowsColl := pgxmock.NewRows([]string{"id", "name", "url"}).AddRow(idByte, []uint8(mockCollections.CollArray[0].Title), []uint8(mockCollections.CollArray[0].PictureUrl))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountCollections)).WithArgs().WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetCollections)).WithArgs(10, 0).WillReturnRows(rowsColl)
	pool.ExpectCommit()

	actual, err := repository.GetCollections(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, mockCollections, actual)
}

func TestGetCollFilmsSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewCollectionsRepository(mdb)
	defer pool.Close()

	FilmList := []domain.Film{{
		Id:        1,
		Title:     "Test_film",
		Rating:    2,
		PosterUrl: "/pic/TestPoster.webp"}}
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, FilmList[0].Id)
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))
	byteStars := make([]byte, 8)
	binary.BigEndian.PutUint64(byteStars, math.Float64bits(2))

	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	rowsRecom := pgxmock.NewRows([]string{"id", "title", "titor", "desc", "date", "poster_url"})
	rowsRate := pgxmock.NewRows([]string{"stars"}).AddRow(byteStars)
	byteRelease := make([]byte, 4)
	binary.BigEndian.PutUint32(byteRelease, uint32(FilmList[0].ReleaseYear))
	rowsRecom.AddRow(byteId, []uint8(FilmList[0].Title), []uint8(FilmList[0].TitleOriginal), []uint8(FilmList[0].Description), byteRelease, []uint8(FilmList[0].PosterUrl))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFilms)).WithArgs(uint64(1)).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilms)).WithArgs(uint64(1)).WillReturnRows(rowsRecom)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmRating)).WithArgs(uint64(1)).WillReturnRows(rowsRate)
	pool.ExpectCommit()

	actual, err := repository.GetCollectionFilms(1)
	assert.NoError(t, err)
	assert.Equal(t, FilmList, actual)
}

func TestGetCollInfoSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewCollectionsRepository(mdb)
	defer pool.Close()

	coll := domain.Collection{
		Id:           1,
		AuthId:       1,
		CollName:     "a",
		Description:  "bb",
		PicUrl:       "d",
		CreationTime: "01.01.2000",
	}
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, coll.Id)
	byteIdAuth := make([]byte, 8)
	binary.BigEndian.PutUint64(byteIdAuth, coll.AuthId)
	byteStars := make([]byte, 8)
	binary.BigEndian.PutUint64(byteStars, math.Float64bits(2))
	byteRelease := make([]byte, 8)
	binary.BigEndian.PutUint32(byteRelease, 0)

	rowsRecom := pgxmock.NewRows([]string{"id", "authid", "name", "desc", "date", "url"})
	rowsRecom.AddRow(byteId, byteIdAuth, []uint8(coll.CollName), []uint8(coll.Description), byteRelease, []uint8(coll.PicUrl))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetCollection)).WithArgs(uint64(1)).WillReturnRows(rowsRecom)
	pool.ExpectCommit()


	actual, err := repository.GetCollectionInfo(1)
	assert.NoError(t, err)
	assert.Equal(t,coll, actual)
}


