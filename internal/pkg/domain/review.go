package domain

import "time"

type Review struct {
	Id                uint64    `json:"id"`
	FilmId            uint64    `json:"film_id"`
	FilmTitleRu       string    `json:"film_title_ru,omitempty"`
	FilmTitleOriginal string    `json:"film_title_original,omitempty"`
	FilmPictureUrl    string    `json:"film_picture_url,omitempty"`
	AuthorId          uint64    `json:"author_id,omitempty"`
	AuthorName        string    `json:"author_name,omitempty"`
	AuthorPictureUrl  string    `json:"author_picture_url,omitempty"`
	ReviewText        string    `json:"review_text,omitempty"`
	ReviewType        int       `json:"review_type"`
	Stars             float64   `json:"stars"`
	Date              time.Time `json:"date"`
}

type ReviewRepository interface {
	GetReview(id uint64) (Review, error)
	PostReview(review Review) (uint64, error)
	LoadReviewsExcept(id uint64, film_id uint64, skip int, limit int) (FilmReviews, error)
}

type ReviewUsecase interface {
	GetReview(id uint64) (Review, error)
	PostReview(review Review) (Review, error)
	LoadReviewsExcept(id uint64, film_id uint64, skip int, limit int) (FilmReviews, error)
}
