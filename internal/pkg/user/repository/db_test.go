package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"

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

func TestAddSuccess (t *testing.T) {
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

