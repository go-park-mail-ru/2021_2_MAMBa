package domain

type SearchResult struct {
	Films   FilmList   `json:"films,omitempty"`
	Persons PersonList `json:"persons,omitempty"`
}

//go:generate mockgen -destination=../search/repository/mock/db_mock.go  -package=mock 2021_2_MAMBa/internal/pkg/domain SearchRepository
type SearchRepository interface {
	CountFoundPersons(query string) (int, error)
	CountFoundFilms(query string) (int, error)
	SearchFilmsIDList(query string, skip int, limit int) ([]uint64, error)
	SearchPersonsIDList(query string, skip int, limit int) ([]uint64, error)
}

//go:generate mockgen -destination=../search/usecase/mock/usecase_mock.go  -package=mock 2021_2_MAMBa/internal/pkg/domain SearchUsecase
type SearchUsecase interface {
	GetSearch(query string, skipFilms int, limitFilms int, skipPersons int, limitPersons int) (SearchResult, error)
}
