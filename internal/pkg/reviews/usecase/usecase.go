package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
)

type ReviewUsecase struct {
	reviewRepo domain.ReviewRepository
}

func NewReviewUsecase(r domain.ReviewRepository) domain.ReviewUsecase {
	return &ReviewUsecase{
		reviewRepo: r,
	}
}

func (uc *ReviewUsecase) GetReview(id uint64) (domain.Review, error) {
	review, err := uc.reviewRepo.GetReview(id)
	if err != nil {
		return domain.Review{}, err
	}
	return review, nil
}
func (uc *ReviewUsecase) PostReview(review domain.Review) (domain.Review, error) {
	id, err := uc.reviewRepo.PostReview(review)
	if err != nil {
		return domain.Review{}, err
	}
	newReview, err := uc.reviewRepo.GetReview(id)
	return newReview, nil
}
func (uc *ReviewUsecase) LoadReviewsExcept(id uint64, film_id uint64, skip int, limit int) (domain.FilmReviews, error) {
	reviews, err := uc.reviewRepo.LoadReviewsExcept(id, film_id, skip, limit)
	if err != nil {
		return domain.FilmReviews{}, err
	}
	return reviews, nil
}
