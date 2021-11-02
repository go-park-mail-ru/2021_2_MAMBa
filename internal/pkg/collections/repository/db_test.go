package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
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
		CollArray:       []domain.CollectionPreview{{
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
	rowsColl := pgxmock.NewRows([]string{"id","name","url"}).AddRow(idByte, []uint8(mockCollections.CollArray[0].Title), []uint8(mockCollections.CollArray[0].PictureUrl))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountCollections)).WithArgs().WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetCollections)).WithArgs(10,0).WillReturnRows(rowsColl)
	pool.ExpectCommit()

	actual, err := repository.GetCollections(0, 10)
	assert.NoError(t, err)
	assert.Equal(t, mockCollections, actual)
}


