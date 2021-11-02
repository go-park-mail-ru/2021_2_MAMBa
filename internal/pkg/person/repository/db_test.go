package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	mylog "2021_2_MAMBa/internal/pkg/utils/log"
	"encoding/binary"
	"errors"
	"github.com/jackc/pgx/pgtype"
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
	repository := NewPersonRepository(mdb)
	defer pool.Close()
	timeBuffer1 := pgtype.Timestamp{}
	err = timeBuffer1.DecodeBinary(nil, nil)
	timeBuffer2 := pgtype.Timestamp{}
	err = timeBuffer2.DecodeBinary(nil, nil)
	pers := domain.Person{
		Id:           1,
		NameEn:       "Test",
		NameRus:      "Тест",
		PictureUrl:   "1.jpeg",
		Career:       []string{"Тестер"},
		Height:       130.0,
		Age:          20,
		Birthday:     timeBuffer1.Time.String(),
		Death:        timeBuffer2.Time.String(),
		BirthPlace:   "",
		DeathPlace:   "",
		Gender:       "",
		FamilyStatus: "",
		FilmNumber:   "",
	}
	byteId := make([]uint8, 8)
	binary.BigEndian.PutUint64(byteId, pers.Id)
	byteHeight := make([]uint8, 4)
	binary.BigEndian.PutUint32(byteHeight, uint32(int(pers.Height)))
	byteAge := make([]uint8, 4)
	binary.BigEndian.PutUint32(byteAge, uint32(pers.Age))

	rowPerson := pgxmock.NewRows([]string{"id", "name-en", "name_ru", "pic-url", "career", "height", "age", "bday", "dday",
		"bp", "dp", "gender", "status", "number"})
	rowPerson.AddRow(byteId, []uint8(pers.NameEn), []uint8(pers.NameRus), []uint8(pers.PictureUrl), []uint8(pers.Career[0]),
		byteHeight, byteAge, []uint8(nil), []uint8(nil), []uint8(""), []uint8(""), []uint8(""), []uint8(""), []uint8(""))
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetPerson)).WithArgs(pers.Id).WillReturnRows(rowPerson)
	pool.ExpectCommit()

	actual, err := repository.GetPerson(pers.Id)
	assert.NoError(t, err)
	assert.Equal(t, pers, actual)
}

func TestGetFilms(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewPersonRepository(mdb)
	defer pool.Close()

	pId := uint64(1)

	fl := domain.FilmList{
		FilmList: []domain.Film{
			domain.Film{
				Id:          1,
				Title:       "TestFilm",
				Description: "ahahah ehehehe",
				PosterUrl:   "joebiden.pres",
				ReleaseYear: 2020,
			},
		},
		MoreAvaliable: false,
		FilmTotal:     1,
		CurrentLimit:  10,
		CurrentSkip:   10,
	}
	byteId := make([]uint8, 8)
	binary.BigEndian.PutUint64(byteId, pId)
	byteRelease := make([]uint8, 4)
	binary.BigEndian.PutUint32(byteRelease, uint32(fl.FilmList[0].ReleaseYear))

	rowCount := pgxmock.NewRows([]string{"count"}).AddRow(byteId)
	rowFilm := pgxmock.NewRows([]string{"id", "title", "desc", "release", "poster"})
	rowFilm.AddRow(byteId, []uint8(fl.FilmList[0].Title), []uint8(fl.FilmList[0].Description), byteRelease, []uint8(fl.FilmList[0].PosterUrl))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFilm)).WithArgs(pId).WillReturnRows(rowCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetPersonFilms)).WithArgs(pId, 10, 0).WillReturnRows(rowFilm)
	pool.ExpectCommit()

	actual, err := repository.GetFilms(uint64(1), 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, fl, actual)
}

func TestGetPopularFilms(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewPersonRepository(mdb)
	defer pool.Close()

	pId := uint64(1)

	fl := domain.FilmList{
		FilmList: []domain.Film{
			domain.Film{
				Id:          1,
				Title:       "TestFilm",
				Description: "ahahah ehehehe",
				PosterUrl:   "joebiden.pres",
				ReleaseYear: 2020,
			},
		},
		MoreAvaliable: false,
		FilmTotal:     1,
		CurrentLimit:  10,
		CurrentSkip:   10,
	}
	byteId := make([]uint8, 8)
	binary.BigEndian.PutUint64(byteId, pId)
	byteRelease := make([]uint8, 4)
	binary.BigEndian.PutUint32(byteRelease, uint32(fl.FilmList[0].ReleaseYear))
	bytefloat := make([]uint8, 8)
	binary.BigEndian.PutUint64(bytefloat, math.Float64bits(10.0))

	rowCount := pgxmock.NewRows([]string{"count"}).AddRow(byteId)
	rowFilm := pgxmock.NewRows([]string{"id", "title", "desc", "release", "poster", "rate"})
	rowFilm.AddRow(byteId, []uint8(fl.FilmList[0].Title), []uint8(fl.FilmList[0].Description), byteRelease, []uint8(fl.FilmList[0].PosterUrl), bytefloat)

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFilm)).WithArgs(pId).WillReturnRows(rowCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetPersonFilmsPopular)).WithArgs(pId, 10, 0).WillReturnRows(rowFilm)
	pool.ExpectCommit()

	actual, err := repository.GetFilmsPopular(uint64(1), 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, fl, actual)
}
