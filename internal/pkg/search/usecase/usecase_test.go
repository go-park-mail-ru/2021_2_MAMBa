package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	mockF "2021_2_MAMBa/internal/pkg/film/repository/mock"
	mockP "2021_2_MAMBa/internal/pkg/person/repository/mock"
	mockS "2021_2_MAMBa/internal/pkg/search/repository/mock"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSearch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cl := domain.SearchResult{
		Films:   domain.FilmList{FilmList: []domain.Film{{Id: 1}}, MoreAvailable: true, FilmTotal: 1, CurrentLimit: 10, CurrentSkip: 10},
		Persons: domain.PersonList{PersonList: []domain.Person{{Id: 1}}, MoreAvailable: true, PersonTotal: 1, CurrentLimit: 10, CurrentSkip: 10},
	}

	mockSearch := mockS.NewMockSearchRepository(ctrl)
	mockPerson := mockP.NewMockPersonRepository(ctrl)
	mockFilm := mockF.NewMockFilmRepository(ctrl)
	usecase := NewSearchUsecase(mockSearch, mockPerson, mockFilm)

	mockSearch.EXPECT().SearchFilmsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().SearchPersonsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().CountFoundFilms("aaa").Return(1, nil)
	mockSearch.EXPECT().CountFoundPersons("aaa").Return(1, nil)
	mockFilm.EXPECT().GetFilm(uint64(1)).Return(domain.Film{Id: 1}, nil)
	mockPerson.EXPECT().GetPerson(uint64(1)).Return(domain.Person{Id: 1}, nil)
	actual, err := usecase.GetSearch("aaa", 0, 10, 0, 10)
	assert.Equal(t, cl, actual)
	assert.Equal(t, err, nil)

	testErr := errors.New("test")

	mockSearch.EXPECT().SearchFilmsIDList("aaa", 0, 10).Return([]uint64{1}, testErr)
	actual, err = usecase.GetSearch("aaa", 0, 10, 0, 10)
	assert.Equal(t, domain.SearchResult{}, actual)
	assert.Equal(t, err, testErr)

	mockSearch.EXPECT().SearchFilmsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().SearchPersonsIDList("aaa", 0, 10).Return([]uint64{1}, err)
	actual, err = usecase.GetSearch("aaa", 0, 10, 0, 10)
	assert.Equal(t, domain.SearchResult{}, actual)
	assert.Equal(t, err, testErr)

	mockSearch.EXPECT().SearchFilmsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().SearchPersonsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().CountFoundFilms("aaa").Return(1, err)
	actual, err = usecase.GetSearch("aaa", 0, 10, 0, 10)
	assert.Equal(t, domain.SearchResult{}, actual)
	assert.Equal(t, err, testErr)

	mockSearch.EXPECT().SearchFilmsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().SearchPersonsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().CountFoundFilms("aaa").Return(1, nil)
	mockSearch.EXPECT().CountFoundPersons("aaa").Return(1, err)
	actual, err = usecase.GetSearch("aaa", 0, 10, 0, 10)
	assert.Equal(t, domain.SearchResult{}, actual)
	assert.Equal(t, err, testErr)

	mockSearch.EXPECT().SearchFilmsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().SearchPersonsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().CountFoundFilms("aaa").Return(1, nil)
	mockSearch.EXPECT().CountFoundPersons("aaa").Return(1, nil)
	mockFilm.EXPECT().GetFilm(uint64(1)).Return(domain.Film{Id: 1}, err)
	actual, err = usecase.GetSearch("aaa", 0, 10, 0, 10)
	assert.Equal(t, domain.SearchResult{}, actual)
	assert.Equal(t, err, testErr)

	mockSearch.EXPECT().SearchFilmsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().SearchPersonsIDList("aaa", 0, 10).Return([]uint64{1}, nil)
	mockSearch.EXPECT().CountFoundFilms("aaa").Return(1, nil)
	mockSearch.EXPECT().CountFoundPersons("aaa").Return(1, nil)
	mockFilm.EXPECT().GetFilm(uint64(1)).Return(domain.Film{Id: 1}, nil)
	mockPerson.EXPECT().GetPerson(uint64(1)).Return(domain.Person{Id: 1}, err)
	actual, err = usecase.GetSearch("aaa", 0, 10, 0, 10)
	assert.Equal(t, domain.SearchResult{}, actual)
	assert.Equal(t, err, testErr)

}
