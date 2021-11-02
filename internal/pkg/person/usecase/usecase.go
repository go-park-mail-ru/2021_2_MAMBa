package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
)

type PersonUsecase struct {
	PersonRepo domain.PersonRepository
}

func NewPersonUsecase(u domain.PersonRepository) domain.PersonUsecase {
	return &PersonUsecase{
		PersonRepo: u,
	}
}

func (uc *PersonUsecase) GetPerson(id uint64) (domain.PersonPage, error) {
	if id == 0 {
		return domain.PersonPage{}, customErrors.ErrorBadInput
	}
	person, err := uc.PersonRepo.GetPerson(id)
	if err != nil {
		return domain.PersonPage{}, err
	}
	films, err := uc.PersonRepo.GetFilms(id, 0, 10)
	if err != nil {
		return domain.PersonPage{}, err
	}
	pop, err := uc.PersonRepo.GetFilmsPopular(id, 0, 10)
	if err != nil {
		return domain.PersonPage{}, err
	}
	return domain.PersonPage{
		Actor:        person,
		Films:        films,
		PopularFilms: pop,
	}, nil
}
func (uc *PersonUsecase) GetFilms(id uint64, skip int, limit int) (domain.FilmList, error) {
	films, err := uc.PersonRepo.GetFilms(id, skip, limit)
	if err != nil {
		return domain.FilmList{}, err
	}
	return films, nil
}
