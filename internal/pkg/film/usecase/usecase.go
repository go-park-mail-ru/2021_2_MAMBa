package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
)

type FilmUsecase struct {
	FilmRepo domain.FilmRepository
}

func NewFilmUsecase(u domain.FilmRepository) domain.FilmUsecase {
	return &FilmUsecase{
		FilmRepo: u,
	}
}

func (uc *FilmUsecase) GetFilm(userID, filmID uint64, skipReviews int, limitReviews int, skipRecommend int, limitRecommend int) (domain.FilmPageInfo, error) {
	film, err := uc.FilmRepo.GetFilm(filmID)
	if err != nil {
		return domain.FilmPageInfo{}, err
	}
	Reviews, err := uc.FilmRepo.GetFilmReviews(filmID, skipReviews, limitReviews)
	if err != nil {
		return domain.FilmPageInfo{}, err
	}
	Recommended, err := uc.FilmRepo.GetFilmRecommendations(filmID, skipReviews, limitRecommend)
	if err != nil {
		return domain.FilmPageInfo{}, err
	}

	myReview := domain.Review{}
	if userID != 0 {
		myReview, err = uc.FilmRepo.GetMyReview(filmID, userID)
		if err != nil && err != customErrors.ErrorNoReviewForFilm {
			return domain.FilmPageInfo{}, err
		}
	}
	result := domain.FilmPageInfo{
		FilmMain:        &film,
		Reviews:         Reviews,
		Recommendations: Recommended,
		MyReview:        myReview,
	}
	return result, nil
}

func (uc *FilmUsecase) PostRating(id uint64, authorId uint64, rating float64) (float64, error) {
	rating, err := uc.FilmRepo.PostRating(id, authorId, rating)
	if err != nil {
		return 0, err
	}
	return rating, nil
}

func (uc *FilmUsecase) LoadFilmReviews(id uint64, skip int, limit int) (domain.FilmReviews, error) {
	Reviews, err := uc.FilmRepo.GetFilmReviews(id, skip, limit)
	if err != nil {
		return domain.FilmReviews{}, err
	}
	return Reviews, nil
}
func (uc *FilmUsecase) LoadFilmRecommendations(id uint64, skip int, limit int) (domain.FilmRecommendations, error) {
	Recommendations, err := uc.FilmRepo.GetFilmRecommendations(id, skip, limit)
	if err != nil {
		return domain.FilmRecommendations{}, err
	}
	return Recommendations, nil
}
func (uc *FilmUsecase) LoadMyReview(id uint64, authorId uint64) (domain.Review, error) {
	myRev, err := uc.FilmRepo.GetMyReview(id, authorId)
	if err != nil {
		return domain.Review{}, err
	}
	return myRev, nil
}
