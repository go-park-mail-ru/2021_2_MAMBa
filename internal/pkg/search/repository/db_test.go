package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	mylog "2021_2_MAMBa/internal/pkg/utils/log"
	"encoding/binary"
	"github.com/pashagolub/pgxmock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
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

func TestCountFoundSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewSearchRepository(mdb)
	defer pool.Close()

	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))

	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFoundFilms)).WithArgs("%%").WillReturnRows(rowsCount)
	pool.ExpectCommit()
	actual, err := repository.CountFoundFilms("")
	assert.NoError(t, err)
	assert.Equal(t, 1, actual)

	rowsCount = pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFoundPersons)).WithArgs("%%").WillReturnRows(rowsCount)
	pool.ExpectCommit()
	actual, err = repository.CountFoundPersons("")
	assert.NoError(t, err)
	assert.Equal(t, 1, actual)
}

func TestSearchIDlistSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewSearchRepository(mdb)
	defer pool.Close()

	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))

	rowsCount := pgxmock.NewRows([]string{"id"}).AddRow(byteCount)
	rowsCount.AddRow(byteCount)
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(querySearchFilmsByString)).WithArgs("%%", 10, 0).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	actual, err := repository.SearchFilmsIDList("", 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1, 1}, actual)

	rowsCount = pgxmock.NewRows([]string{"id"}).AddRow(byteCount)
	rowsCount.AddRow(byteCount)
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(querySearchPersonsByString)).WithArgs("%%", 10, 0).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	actual, err = repository.SearchPersonsIDList("", 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, []uint64{1, 1}, actual)

}
