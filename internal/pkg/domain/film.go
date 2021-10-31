package domain

import "time"

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

type Review struct {
	Id                uint64    `json:"id"`
	FilmId            uint64    `json:"film_id"`
	FilmTitleRu       string    `json:"film_title_ru,omitempty"`
	FilmTitleOriginal string    `json:"film_title_original,omitempty"`
	FilmPictureUrl    string    `json:"film_picture_url,omitempty"`
	AuthorName        string    `json:"author_name,omitempty"`
	AuthorPictureUrl  string    `json:"author_picture_url,omitempty"`
	ReviewText        string    `json:"review_text,omitempty"`
	ReviewType        int       `json:"review_type,omitempty"`
	Stars             float64   `json:"stars"`
	Date              time.Time `json:"date"`
}

type FilmRecommendations struct {
	RecommendationList  []Film `json:"recommendation_list"`
	MoreAvaliable       bool   `json:"more_avaliable"`
	RecommendationTotal int    `json:"recommendation_total"`
	CurrentLimit        int    `json:"current_limit"`
	CurrentSkip         int    `json:"current_skip"`
}

type FilmReviews struct {
	ReviewList    []Review `json:"review_list"`
	MoreAvaliable bool     `json:"more_avaliable"`
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
	GetFilmReviews (id uint64, skip int, limit int) (FilmReviews, error)
	GetFilmRecommendations (id uint64, skip int, limit int) (FilmRecommendations, error)
	PostRating (id uint64, authorId uint64, rating float64) (float64, error)
	GetMyReview (id uint64, authorId uint64) (Review, error)
}

type FilmUsecase interface {
	GetFilm(id uint64, skipReviews int, limitReviews int, skipRecommend int, limitRecommend int) (FilmPageInfo, error)
	PostRating(id uint64, authorId uint64, rating float64) (float64, error)
	LoadFilmReviews(id uint64, skip int, limit int) (FilmReviews, error)
	LoadFilmRecommendations (id uint64, skip int, limit int) (FilmRecommendations, error)
	LoadMyReview (id uint64, authorId uint64) (Review, error)
}
