package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
)

type SearchUsecase struct {
	searchRepo domain.SearchRepository
	personRepo domain.PersonRepository
	filmRepo   domain.FilmRepository
}

func NewSearchUsecase(sr domain.SearchRepository, pr domain.PersonRepository, fr domain.FilmRepository) domain.SearchUsecase {
	return &SearchUsecase{
		searchRepo: sr,
		personRepo: pr,
		filmRepo:   fr,
	}
}

func (uc *SearchUsecase) GetSearch(query string, skipFilms int, limitFilms int, skipPersons int, limitPersons int) (domain.SearchResult, error) {
	filmsIDList, err := uc.searchRepo.SearchFilmsIDList(query, skipFilms, limitFilms)
	if err != nil {
		return domain.SearchResult{}, err
	}

	personsIDList, err := uc.searchRepo.SearchPersonsIDList(query, skipPersons, limitPersons)
	if err != nil {
		return domain.SearchResult{}, err
	}

	totalCountFilms, err := uc.searchRepo.CountFoundFilms(query)
	if err != nil {
		return domain.SearchResult{}, err
	}
	if skipFilms >= totalCountFilms && skipFilms != 0 {
		return domain.SearchResult{}, customErrors.ErrorSkip
	}
	moreAvailableFilms := skipFilms+skipFilms < totalCountFilms

	totalCountPersons, err := uc.searchRepo.CountFoundPersons(query)
	if err != nil {
		return domain.SearchResult{}, err
	}
	if skipPersons >= totalCountPersons && skipPersons != 0 {
		return domain.SearchResult{}, customErrors.ErrorSkip
	}
	moreAvailablePersons := skipPersons+skipPersons < totalCountPersons

	filmsList := make([]domain.Film, 0)
	for _, filmID := range filmsIDList {
		film, err := uc.filmRepo.GetFilm(filmID)
		if err != nil {
			return domain.SearchResult{}, err
		}
		filmsList = append(filmsList, film)
	}

	personsList := make([]domain.Person, 0)
	for _, personID := range personsIDList {
		person, err := uc.personRepo.GetPerson(personID)
		if err != nil {
			return domain.SearchResult{}, err
		}
		personsList = append(personsList, person)
	}

	result := domain.SearchResult{
		Films: domain.FilmList{
			FilmList:      filmsList,
			MoreAvailable: moreAvailableFilms,
			FilmTotal:     totalCountFilms,
			CurrentLimit:  limitFilms,
			CurrentSkip:   skipFilms + limitFilms,
		},
		Persons: domain.PersonList{
			PersonList:    personsList,
			MoreAvailable: moreAvailablePersons,
			PersonTotal:   totalCountPersons,
			CurrentLimit:  limitPersons,
			CurrentSkip:   skipPersons + limitPersons,
		},
	}
	return result, nil
}
