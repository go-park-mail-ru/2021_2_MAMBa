package usecase

import "2021_2_MAMBa/internal/pkg/domain"

type FilmUsecase struct {
	FilmRepo domain.FilmRepository
}

func NewFilmUsecase(u domain.FilmRepository) domain.FilmUsecase {
	return &FilmUsecase{
		FilmRepo: u,
	}
}

func (uc *FilmUsecase) GetFilm(id uint64, skipReviews int, limitReviews int, skipRecommend int, limitRecommend int) (domain.FilmPageInfo, error) {
	film, err := uc.FilmRepo.GetFilm(id)
	if err != nil {
		return domain.FilmPageInfo{}, err
	}
	Reviews, err := uc.FilmRepo.GetFilmReviews(id, skipReviews, limitReviews)
	if err != nil {
		return domain.FilmPageInfo{}, err
	}
	Recommended , err := uc.FilmRepo.GetFilmRecommendations(id, skipReviews, limitRecommend)
	if err != nil {
		return domain.FilmPageInfo{}, err
	}
	result := domain.FilmPageInfo{
		Film:            film,
		Reviews:         Reviews,
		Recommendations: Recommended,
		MyRating:        -1,
	}
	return result, nil
}