package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	mylog "2021_2_MAMBa/internal/pkg/utils/log"
	"encoding/binary"
	"errors"
	"github.com/jackc/pgtype"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"math"
	"regexp"
	"testing"
)

func MockDatabase() (*database.DBManager, pgxmock.PgxPoolIface, error) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		mylog.Error(errors.New("failed to create mock"))
	}
	return &database.DBManager{Pool: mock}, mock, err
}

func TestGetSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewReviewRepository(mdb)
	defer pool.Close()

	r := domain.Review{
		Id:                1,
		FilmId:            1,
		FilmTitleRu:       "тест",
		FilmTitleOriginal: "TestFilm",
		FilmPictureUrl:    "12.jpg",
		AuthorName:        "Ivan Ivanov",
		AuthorPictureUrl:  "auth.jpg",
		ReviewText:        "sdknldsnvlksnlbdlsnbljD;ABSIFBLASN",
		ReviewType:        3,
		Stars:             0,
		Date:              "",
	}

	byteId := make([]uint8, 8)
	binary.BigEndian.PutUint64(byteId, r.Id)
	timeBuffer := pgtype.Timestamp{}
	byteTime, err := timeBuffer.EncodeBinary(nil, nil)
	bytetype := make([]uint8, 4)
	binary.BigEndian.PutUint32(bytetype, uint32(r.ReviewType))
	byteStars := make([]byte, 8)
	binary.BigEndian.PutUint64(byteStars, math.Float64bits(r.Stars))

	rowsRev := pgxmock.NewRows([]string{"review_id", "film_id", "authid", "text", "type", "stars", "date"})
	rowsRev.AddRow(byteId, byteId, byteId, []uint8(r.ReviewText), bytetype, byteStars, byteTime)
	rowsAuth := pgxmock.NewRows([]string{"name", "surname", "url"}).AddRow([]uint8("Ivan"), []uint8("Ivanov"), []uint8(r.AuthorPictureUrl))
	rowsFilm := pgxmock.NewRows([]string{"title", "title_o", "poster"}).AddRow([]uint8(r.FilmTitleRu), []uint8(r.FilmTitleOriginal), []uint8(r.FilmPictureUrl))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetReviewByID)).WithArgs(r.Id).WillReturnRows(rowsRev)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetAuthorName)).WithArgs(r.Id).WillReturnRows(rowsAuth)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmShort)).WithArgs(r.Id).WillReturnRows(rowsFilm)
	pool.ExpectCommit()

	actual, err := repository.GetReview(r.Id)
	assert.NoError(t, err)
	assert.Equal(t, r, actual)
}

func TestPostSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewReviewRepository(mdb)
	defer pool.Close()

	r := domain.Review{
		Id:                1,
		FilmId:            1,
		FilmTitleRu:       "тест",
		FilmTitleOriginal: "TestFilm",
		FilmPictureUrl:    "12.jpg",
		AuthorName:        "Ivan Ivanov",
		AuthorId:          1,
		AuthorPictureUrl:  "auth.jpg",
		ReviewText:        "sdknldsnvlksnlbdlsnbljD;ABSIFBLASN",
		ReviewType:        3,
		Stars:             0,
		Date:              "",
	}

	byteId := make([]uint8, 8)
	binary.BigEndian.PutUint64(byteId, r.Id)
	timeBuffer := pgtype.Timestamp{}
	byteTime, err := timeBuffer.EncodeBinary(nil, nil)
	bytetype := make([]uint8, 4)
	binary.BigEndian.PutUint32(bytetype, uint32(r.ReviewType))
	byteStars := make([]byte, 8)
	binary.BigEndian.PutUint64(byteStars, math.Float64bits(r.Stars))

	rowsRev := pgxmock.NewRows([]string{"review_id", "film_id", "authid", "text", "type", "stars", "date"})
	rowsRev.AddRow(byteId, byteId, byteId, []uint8(r.ReviewText), bytetype, byteStars, byteTime)
	rowsId := pgxmock.NewRows([]string{"newID"}).AddRow(byteId)

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(querySearchReview)).WithArgs(r.FilmId, r.AuthorId).WillReturnRows(rowsRev)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryUpdateReview)).WithArgs(r.ReviewText, r.ReviewType, r.FilmId, r.AuthorId).WillReturnRows(rowsId)
	pool.ExpectCommit()

	actual, err := repository.PostReview(r)
	assert.NoError(t, err)
	assert.Equal(t, r.Id, actual)
}

func TestGetExceptSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewReviewRepository(mdb)
	defer pool.Close()

	r := domain.Review{
		Id:                1,
		FilmId:            1,
		FilmTitleRu:       "тест",
		FilmTitleOriginal: "TestFilm",
		FilmPictureUrl:    "12.jpg",
		AuthorName:        "Ivan Ivanov",
		AuthorPictureUrl:  "auth.jpg",
		ReviewText:        "sdknldsnvlksnlbdlsnbljD;ABSIFBLASN",
		ReviewType:        3,
		Stars:             0,
		Date:              "",
	}

	rlist := domain.FilmReviews{
		ReviewList:    []domain.Review{r},
		MoreAvailable: false,
		ReviewTotal:   2,
		CurrentLimit:  10,
		CurrentSkip:   10,
	}
	byteCount := make([]uint8, 8)
	binary.BigEndian.PutUint64(byteCount, r.Id+1)
	byteId := make([]uint8, 8)
	binary.BigEndian.PutUint64(byteId, r.Id)
	timeBuffer := pgtype.Timestamp{}
	byteTime, err := timeBuffer.EncodeBinary(nil, nil)
	bytetype := make([]uint8, 4)
	binary.BigEndian.PutUint32(bytetype, uint32(r.ReviewType))
	byteStars := make([]byte, 8)
	binary.BigEndian.PutUint64(byteStars, math.Float64bits(r.Stars))

	rowsCount := pgxmock.NewRows([]string{"newID"}).AddRow(byteCount)
	rowsRev := pgxmock.NewRows([]string{"review_id", "film_id", "authid", "text", "type", "stars", "date"})
	rowsRev.AddRow(byteId, byteId, byteId, []uint8(r.ReviewText), bytetype, byteStars, byteTime)
	rowsAuth := pgxmock.NewRows([]string{"name", "surname", "url"}).AddRow([]uint8("Ivan"), []uint8("Ivanov"), []uint8(r.AuthorPictureUrl))
	rowsFilm := pgxmock.NewRows([]string{"title", "title_o", "poster"}).AddRow([]uint8(r.FilmTitleRu), []uint8(r.FilmTitleOriginal), []uint8(r.FilmPictureUrl))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFilmReviews)).WithArgs(r.FilmId).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetReviewByFilmIDEXCEPT)).WithArgs(r.FilmId, r.Id+1, 10, 0).WillReturnRows(rowsRev)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetAuthorName)).WithArgs(r.Id).WillReturnRows(rowsAuth)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmShort)).WithArgs(r.Id).WillReturnRows(rowsFilm)
	pool.ExpectCommit()

	actual, err := repository.LoadReviewsExcept(r.Id+1, r.FilmId, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, rlist, actual)
}
