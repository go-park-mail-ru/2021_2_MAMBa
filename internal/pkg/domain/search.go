package domain

type SearchResult struct {
	films   FilmList
	persons PersonList
}

type SearchRepository interface {
	GetSearch(query string, skipFilms int, limitFilms int, skipPersons int, limitPersons int) (SearchResult, error)
}

type SearchUsecase interface {
	GetSearch(query string, skipFilms int, limitFilms int, skipPersons int, limitPersons int) (SearchResult, error)
}
