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
	Id:           1,
	NameEn:       "Miley",
	NameRus:      "Сайрус",
	PictureUrl:   "/miley.webp",
	Career:       []string{"Актриса"},
}

func TestGetSuccess(t *testing.T) {
	mylog.Info("test get success")
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()

	mf := domain.Film{
		Id:              1,
		Title:           "Тест",
		TitleOriginal:   "Test",
		Rating:          4.1,
		Description:     "this is a test film",
		TotalRevenue:    "$100",
		PosterUrl:       "/img/poster1.webp",
		TrailerUrl:      "/trailer/mov.mov",
		ContentType:     "film",
		ReleaseYear:     2021,
		Duration:        80,
		OriginCountries: []string {"США", "Канада"} ,
		Cast:            []domain.Person{mockPersonPreview},
		Director:        mockPersonPreview,
		Screenwriter:    mockPersonPreview,
		Genres:          []domain.Genre{{Id: 1, Name: "Комедия"}},
	}
	rowsFilm := pgxmock.NewRows([]string{"id", "title", "title_original", "rating", "description",
		"poster_url", "trailer", "total_revenue", "release", "duration", "scr", "dir",
		"cont_type", "pid", "pnameen", "pnameru", "ppicture", "pcareer",
		"p1id", "p1nameen", "p1nameru", "p1picture", "p1career"})
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, mf.Id)
	byteRating := make([]byte, 8)
	binary.BigEndian.PutUint64(byteRating, math.Float64bits(mf.Rating))
	byteRelease := make([]byte, 4)
	binary.BigEndian.PutUint32(byteRelease, uint32(mf.ReleaseYear))
	byteDuration := make([]byte, 4)
	binary.BigEndian.PutUint32(byteDuration, uint32(mf.Duration))
	byteGenreId := make([]byte, 4)
	binary.BigEndian.PutUint32(byteGenreId, 1)
	rowsFilm.AddRow(byteId, []uint8(mf.Title), []uint8(mf.TitleOriginal), byteRating, []uint8(mf.Description),
		[]uint8(mf.PosterUrl), []uint8(mf.TrailerUrl), []uint8(mf.TotalRevenue), byteRelease, byteDuration, byteId, byteId,
		[]uint8(mf.ContentType), byteId, []uint8(mf.Director.NameEn), []uint8(mf.Director.NameRus), []uint8(mf.Director.PictureUrl), []uint8(mf.Director.Career[0]),
		byteId, []uint8(mf.Director.NameEn), []uint8(mf.Director.NameRus), []uint8(mf.Director.PictureUrl), []uint8(mf.Director.Career[0]))
	rowsRate := pgxmock.NewRows([]string{"rating"}).AddRow(byteRating)
	rowsCountry := pgxmock.NewRows([]string{"country_name"}).AddRow([]uint8("США"))
	rowsCountry.AddRow([]uint8("Канада"))
	rowsGenres:= pgxmock.NewRows([]string{"id", "genre_name"}).AddRow(byteGenreId, []uint8(mf.Genres[0].Name))
	rowsCast := pgxmock.NewRows([]string{"pid", "pnameen", "pnameru", "ppicture", "pcareer"})
	rowsCast.AddRow(byteId, []uint8(mf.Director.NameEn), []uint8(mf.Director.NameRus), []uint8(mf.Director.PictureUrl), []uint8(mf.Director.Career[0]))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmDirScr)).WithArgs(mf.Id).WillReturnRows(rowsFilm)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmRating)).WithArgs(mf.Id).WillReturnRows(rowsRate)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmCountries)).WithArgs(mf.Id).WillReturnRows(rowsCountry)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmGenres)).WithArgs(mf.Id).WillReturnRows(rowsGenres)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmCast)).WithArgs(mf.Id).WillReturnRows(rowsCast)
	pool.ExpectCommit()

	actual, err := repository.GetFilm(mf.Id)
	assert.NoError(t, err)
	assert.Equal(t, mf, actual)
	mylog.Info("test get success done")
}

