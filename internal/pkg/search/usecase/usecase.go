package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
)

type SearchUsecase struct {
	searchRepo domain.SearchRepository
}

func NewSearchUsecase(r domain.SearchRepository) domain.SearchUsecase {
	return &SearchUsecase{
		searchRepo: r,
	}
}

func (uc *SearchUsecase) GetSearch(query string, skipFilms int, limitFilms int, skipPersons int, limitPersons int) (domain.SearchResult, error) {
	result, err := uc.searchRepo.GetSearch(query, skipFilms, limitFilms, skipPersons, limitPersons)
	if err != nil {
		return domain.SearchResult{}, err
	}

	return result, nil
}
