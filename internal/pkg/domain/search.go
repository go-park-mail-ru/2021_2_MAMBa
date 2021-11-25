package domain

type SearchResult struct {
	Films   FilmList   `json:"films,omitempty"`
	Persons PersonList `json:"persons,omitempty"`
}

type SearchRepository interface {
	CountFoundPersons(query string) (int, error)
	CountFoundFilms(query string) (int, error)
	SearchFilmsIDList(query string, skip int, limit int) ([]uint64, error)
	SearchPersonsIDList(query string, skip int, limit int) ([]uint64, error)
}

type SearchUsecase interface {
	GetSearch(query string, skipFilms int, limitFilms int, skipPersons int, limitPersons int) (SearchResult, error)
}
