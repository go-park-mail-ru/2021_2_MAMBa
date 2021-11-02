package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/pgtype"
	"math"
	"time"

	mylog "2021_2_MAMBa/internal/pkg/utils/log"
	"encoding/binary"
	"errors"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
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

func TestGetSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewUserRepository(mdb)
	defer pool.Close()

	mu := domain.User{
		ID:             1,
		FirstName:      "Test",
		Surname:        "Testovich",
		Email:          "Testosteron@mail.ru",
		Password:       "abcd1234",
		PasswordRepeat: "",
		ProfilePic:     "/pic/1.jpg",
	}
	rows := pgxmock.NewRows([]string{"id", "firstname", "surname", "email", "password", "profilepic"})
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, mu.ID)
	rows.AddRow(byteId, []uint8(mu.FirstName), []uint8(mu.Surname), []uint8(mu.Email), []uint8(mu.Password), []uint8(mu.ProfilePic))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetById)).WithArgs(mu.ID).WillReturnRows(rows)
	pool.ExpectCommit()

	actual, err := repository.GetUserById(mu.ID)
	assert.NoError(t, err)
	assert.Equal(t, mu, actual)
}

func TestGetFailure(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewUserRepository(mdb)
	defer pool.Close()

	mu := domain.User{
		ID: 1,
	}
	rows := pgxmock.NewRows([]string{"id", "firstname", "surname", "email", "password", "profilepic"})
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, mu.ID)

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetById)).WithArgs(mu.ID).WillReturnRows(rows)
	pool.ExpectCommit()

	actual, err := repository.GetUserById(mu.ID)
	assert.Equal(t, customErrors.ErrorNoUser, err)
	assert.Equal(t, domain.User{}, actual)
}

func TestGetEmailSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewUserRepository(mdb)
	defer pool.Close()

	mu := domain.User{
		ID:             1,
		FirstName:      "Test",
		Surname:        "Testovich",
		Email:          "Testosteron@mail.ru",
		Password:       "abcd1234",
		PasswordRepeat: "",
		ProfilePic:     "/pic/1.jpg",
	}
	rows := pgxmock.NewRows([]string{"id", "firstname", "surname", "email", "password", "profilepic"})
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, mu.ID)
	rows.AddRow(byteId, []uint8(mu.FirstName), []uint8(mu.Surname), []uint8(mu.Email), []uint8(mu.Password), []uint8(mu.ProfilePic))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetByEmail)).WithArgs(mu.Email).WillReturnRows(rows)
	pool.ExpectCommit()

	actual, err := repository.GetUserByEmail(mu.Email)
	assert.NoError(t, err)
	assert.Equal(t, mu, actual)
}

func TestGetEmailFailure(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewUserRepository(mdb)
	defer pool.Close()

	mu := domain.User{
		ID:    1,
		Email: "Testosteron@mail.ru",
	}
	rows := pgxmock.NewRows([]string{"id", "firstname", "surname", "email", "password", "profilepic"})
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, mu.ID)

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetByEmail)).WithArgs(mu.Email).WillReturnRows(rows)
	pool.ExpectCommit()

	actual, err := repository.GetUserByEmail(mu.Email)
	assert.Equal(t, customErrors.ErrorNoUser, err)
	assert.Equal(t, domain.User{}, actual)
}

func TestAddSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewUserRepository(mdb)
	defer pool.Close()

	mu := domain.User{
		ID:             1,
		FirstName:      "Test",
		Surname:        "Testovich",
		Email:          "Testosteron@mail.ru",
		Password:       "HASHED_PASS",
		PasswordRepeat: "",
		ProfilePic:     "/pic/1.jpg",
	}
	rows := pgxmock.NewRows([]string{"id"})
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, mu.ID)
	rows.AddRow(byteId)

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryAddUser)).WithArgs(mu.FirstName, mu.Surname, mu.Email, mu.Password, mu.ProfilePic).WillReturnRows(rows)
	pool.ExpectCommit()

	actual, err := repository.AddUser(&mu)
	assert.NoError(t, err)
	assert.Equal(t, mu.ID, actual)
}

func TestGetReviewsSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewUserRepository(mdb)
	defer pool.Close()
	timeBuffer := pgtype.Timestamp{}
	timeBuffer.Time = time.Now()
	reviews := domain.FilmReviews{
		ReviewList: []domain.Review{domain.Review{
			Id:                1,
			FilmId:            1,
			FilmTitleRu:       "Фильм",
			FilmTitleOriginal: "Film",
			FilmPictureUrl:    "Film.jpg",
			AuthorName:        "Ivan Ivanovich",
			ReviewText:        "Test review on film",
			AuthorPictureUrl:  "pic1.jopeg",
			ReviewType:        3,
			Stars:             4.0,
			Date:              time.Time{},
		}},
		MoreAvaliable: false,
		ReviewTotal:   1,
		CurrentLimit:  10,
		CurrentSkip:   10,
	}
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, uint64(1))
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))
	byteTime := make([]byte, 16)
	byteTime, err = timeBuffer.EncodeBinary(nil, byteTime)
	byteType := make([]byte, 4)
	binary.BigEndian.PutUint32(byteType, uint32(reviews.ReviewList[0].ReviewType))
	byteStars := make([]byte, 8)
	binary.BigEndian.PutUint64(byteStars, math.Float64bits(reviews.ReviewList[0].Stars))

	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	rowsRev := pgxmock.NewRows([]string{"review_id", "film_id", "authid", "text", "type", "stars", "date"})
	rowsRev.AddRow(byteId, byteId, byteId, []uint8(reviews.ReviewList[0].ReviewText), byteType, byteStars, byteTime)
	rowsAuth := pgxmock.NewRows([]string{"name", "sname", "url"}).AddRow([]uint8("Ivan"), []uint8("Ivanovich"), []uint8(reviews.ReviewList[0].AuthorPictureUrl))
	rowsFilm := pgxmock.NewRows([]string{"title", "otitle", "url"}).AddRow([]uint8("Фильм"), []uint8("Film"), []uint8("Film.jpg"))
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFilmReviews)).WithArgs(uint64(1)).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetReviewByUserID)).WithArgs(uint64(1), 10, 0).WillReturnRows(rowsRev)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetAuthorName)).WithArgs(uint64(1)).WillReturnRows(rowsAuth)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmShort)).WithArgs(uint64(1)).WillReturnRows(rowsFilm)
	pool.ExpectCommit()

	actual, err := repository.LoadUserReviews(1, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, reviews, actual)
}

func TestGetProfileSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewUserRepository(mdb)
	defer pool.Close()

	timeBuffer := pgtype.Timestamp{}
	timeBuffer.EncodeBinary(nil, nil)
	byteTime := make([]byte, 16)
	byteTime, err = timeBuffer.EncodeBinary(nil, byteTime)

	mu := domain.Profile{
		ID:            1,
		FirstName:     "Test",
		Surname:       "Testovich",
		Email:         "Testosteron@mail.ru",
		PictureUrl:    "/pic/1.jpg",
		Gender:        "male",
		RegisterDate:  timeBuffer.Time,
		SubCount:      2,
		BookmarkCount: 2,
		AmSubscribed:  false,
	}
	rows := pgxmock.NewRows([]string{"id", "firstname", "surname", "email", "password", "profilepic", "gen", "date"})
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, mu.ID)
	rows.AddRow(byteId, []uint8(mu.FirstName), []uint8(mu.Surname), []uint8(mu.Email), []uint8("abcd1234"), []uint8(mu.PictureUrl), []uint8(mu.Gender), byteTime)
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(2))
	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	rowsCount2 := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetById)).WithArgs(mu.ID).WillReturnRows(rows)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountBookmarksById)).WithArgs(mu.ID).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountSubscribersById)).WithArgs(mu.ID).WillReturnRows(rowsCount2)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCheckSubscription)).WithArgs(uint64(2), mu.ID).WillReturnRows(rowsCount)
	pool.ExpectCommit()

	actual, err := repository.GetProfileById(uint64(2), mu.ID)
	assert.NoError(t, err)
	assert.Equal(t, mu, actual)
}

func TestUpdateProfileSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewUserRepository(mdb)
	defer pool.Close()

	timeBuffer := pgtype.Timestamp{}
	timeBuffer.EncodeBinary(nil, nil)
	byteTime := make([]byte, 16)
	byteTime, err = timeBuffer.EncodeBinary(nil, byteTime)

	mu := domain.Profile{
		ID:            1,
		FirstName:     "Test",
		Surname:       "Testovich",
		Email:         "Testosteron@mail.ru",
		PictureUrl:    "/pic/1.jpg",
		Gender:        "male",
		RegisterDate:  timeBuffer.Time,
		SubCount:      2,
		BookmarkCount: 2,
		AmSubscribed:  false,
	}
	rows := pgxmock.NewRows([]string{"id", "firstname", "surname", "email", "password", "profilepic", "gen", "date"})
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, mu.ID)
	rows.AddRow(byteId, []uint8(mu.FirstName), []uint8(mu.Surname), []uint8(mu.Email), []uint8("abcd1234"), []uint8(mu.PictureUrl), []uint8(mu.Gender), byteTime)
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(2))
	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	rowsCount2 := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)

	pool.ExpectBegin()
	pool.ExpectExec(regexp.QuoteMeta(queryUpdProfile)).WithArgs(mu.ID, mu.FirstName, mu.Surname, mu.PictureUrl, mu.Email, mu.Gender).WillReturnResult(pgconn.CommandTag{})
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetById)).WithArgs(mu.ID).WillReturnRows(rows)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountBookmarksById)).WithArgs(mu.ID).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountSubscribersById)).WithArgs(mu.ID).WillReturnRows(rowsCount2)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCheckSubscription)).WithArgs(uint64(1), mu.ID).WillReturnRows(rowsCount)
	pool.ExpectCommit()

	actual, err := repository.UpdateProfile(mu)
	assert.NoError(t, err)
	assert.Equal(t, mu, actual)
}
