package domain

import (
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"encoding/json"
	"net/http"
)

type Response struct {
	Body   json.RawMessage `json:"body"`
	Status int    `json:"status"`
}

func (r *Response) Write (w http.ResponseWriter)  {
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		http.Error(w, customErrors.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

type Country struct {
	Id          uint64
	CountryName string
}

type Genre struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type Film struct {
	Id              uint64   `json:"id,omitempty"`
	Title           string   `json:"title,omitempty"`
	TitleOriginal   string   `json:"title_original,omitempty"`
	Rating          float64  `json:"rating,omitempty"`
	Description     string   `json:"description,omitempty"`
	TotalRevenue    string   `json:"total_revenue,omitempty"`
	PosterUrl       string   `json:"poster_url,omitempty"`
	TrailerUrl      string   `json:"trailer_url,omitempty"`
	ContentType     string   `json:"content_type,omitempty"`
	ReleaseYear     int      `json:"release_year,omitempty"`
	Duration        int      `json:"duration,omitempty"`
	OriginCountries []string `json:"origin_countries,omitempty"`
	Cast            []Person `json:"cast,omitempty"`
	Director        Person   `json:"director,omitempty"`
	Screenwriter    Person   `json:"screenwriter,omitempty"`
	Genres          []Genre  `json:"genres,omitempty"`
}

type FilmRecommendations struct {
	RecommendationList  []Film `json:"recommendation_list"`
	MoreAvailable       bool   `json:"more_available"`
	RecommendationTotal int    `json:"recommendation_total"`
	CurrentLimit        int    `json:"current_limit"`
	CurrentSkip         int    `json:"current_skip"`
}

type FilmReviews struct {
	ReviewList    []Review `json:"review_list"`
	MoreAvailable bool     `json:"more_available"`
	ReviewTotal   int      `json:"review_total"`
	CurrentSort   string   `json:"current_sort"`
	CurrentLimit  int      `json:"current_limit"`
	CurrentSkip   int      `json:"current_skip"`
}

type FilmPageInfo struct {
	Film            Film                `json:"film"`
	Reviews         FilmReviews         `json:"reviews"`
	Recommendations FilmRecommendations `json:"recommendations"`
	MyRating        float64             `json:"my_rating"`
}

type NewRate struct {
	Rating float64 `json:"rating,omitempty"`
}

type FilmRepository interface {
	GetFilm(id uint64) (Film, error)
	GetFilmReviews(id uint64, skip int, limit int) (FilmReviews, error)
	GetFilmRecommendations(id uint64, skip int, limit int) (FilmRecommendations, error)
	PostRating(id uint64, authorId uint64, rating float64) (float64, error)
	GetMyReview(id uint64, authorId uint64) (Review, error)
}

//go:generate mockgen -destination=../film/usecase/mock/usecase_mock.go  -package=mock 2021_2_MAMBa/internal/pkg/domain FilmUsecase
type FilmUsecase interface {
	GetFilm(id uint64, skipReviews int, limitReviews int, skipRecommend int, limitRecommend int) (FilmPageInfo, error)
	PostRating(id uint64, authorId uint64, rating float64) (float64, error)
	LoadFilmReviews(id uint64, skip int, limit int) (FilmReviews, error)
	LoadFilmRecommendations(id uint64, skip int, limit int) (FilmRecommendations, error)
	LoadMyReview(id uint64, authorId uint64) (Review, error)
}
