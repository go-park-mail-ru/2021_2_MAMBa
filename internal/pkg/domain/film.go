package domain

import (
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"encoding/json"
	"net/http"
	"strconv"
)

type Response struct {
	Body   json.RawMessage `json:"body,omitempty"`
	Error  json.RawMessage `json:"error,omitempty"`
	Status int             `json:"status"`
}

func (r *Response) Write(w http.ResponseWriter) {
	w.WriteHeader(r.Status)
	x, err := json.Marshal(r)
	_, err = w.Write(x)
	if err != nil {
		http.Error(w, customErrors.ErrEncMsg, http.StatusInternalServerError)
		return
	}
}

type Banner struct {
	Id          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	PictureURL  string `json:"picture_url,omitempty"`
	Link        string `json:"link,omitempty"`
}

type BannersList struct {
	BannersList []Banner `json:"banners_list"`
}

type Country struct {
	Id          uint64
	CountryName string
}

type Genre struct {
	Id         uint64 `json:"id"`
	Name       string `json:"name"`
	PictureURL string `json:"picture_url,omitempty"`
}

type GenresList struct {
	GenresList []Genre `json:"genres_list"`
}

type GenreFilmList struct {
	Id        uint64   `json:"id"`
	Name      string   `json:"name"`
	FilmsList FilmList `json:"films"`
}

//easyjson:skip
type Film struct {
	Id              uint64   `json:"id,omitempty"`
	Title           string   `json:"title,omitempty"`
	TitleOriginal   string   `json:"title_original,omitempty"`
	Rating          float64  `json:"rating"`
	Description     string   `json:"description,omitempty"`
	TotalRevenue    string   `json:"total_revenue,omitempty"`
	PosterUrl       string   `json:"poster_url,omitempty"`
	TrailerUrl      string   `json:"trailer_url,omitempty"`
	ContentType     string   `json:"content_type,omitempty"`
	ReleaseYear     int      `json:"release_year,omitempty"`
	Duration        int      `json:"duration,omitempty"`
	PremiereRu      string   `json:"premiere_ru,omitempty"`
	OriginCountries []string `json:"origin_countries,omitempty"`
	Cast            []Person `json:"cast,omitempty"`
	Director        Person   `json:"director,omitempty"`
	Screenwriter    Person   `json:"screenwriter,omitempty"`
	Genres          []Genre  `json:"genres,omitempty"`
}

type FilmJson struct {
	Id              uint64      `json:"id,omitempty"`
	Title           string      `json:"title,omitempty"`
	TitleOriginal   string      `json:"title_original,omitempty"`
	Rating          json.Number `json:"rating"`
	Description     string      `json:"description,omitempty"`
	TotalRevenue    string      `json:"total_revenue,omitempty"`
	PosterUrl       string      `json:"poster_url,omitempty"`
	TrailerUrl      string      `json:"trailer_url,omitempty"`
	ContentType     string      `json:"content_type,omitempty"`
	ReleaseYear     int         `json:"release_year,omitempty"`
	Duration        int         `json:"duration,omitempty"`
	PremiereRu      string      `json:"premiere_ru,omitempty"`
	OriginCountries []string    `json:"origin_countries,omitempty"`
	Cast            []Person    `json:"cast,omitempty"`
	Director        Person      `json:"director,omitempty"`
	Screenwriter    Person      `json:"screenwriter,omitempty"`
	Genres          []Genre     `json:"genres,omitempty"`
}

func (film *Film) toJsonNum() FilmJson {
	s := json.Number(strconv.FormatFloat(film.Rating, 'f', 1, 64))
	return FilmJson{
		Id:              film.Id,
		Title:           film.Title,
		TitleOriginal:   film.TitleOriginal,
		Rating:          s,
		Description:     film.Description,
		TotalRevenue:    film.TotalRevenue,
		PosterUrl:       film.PosterUrl,
		TrailerUrl:      film.TrailerUrl,
		ContentType:     film.ContentType,
		ReleaseYear:     film.ReleaseYear,
		Duration:        film.Duration,
		PremiereRu:      film.PremiereRu,
		OriginCountries: film.OriginCountries,
		Cast:            film.Cast,
		Director:        film.Director,
		Screenwriter:    film.Screenwriter,
		Genres:          film.Genres,
	}
}

/*
func (f *Film) CustomEasyJSON() ([]byte, error) {
	return f.toJsonNum().MarshalJSON()
}*/

func (film *Film) MarshalJSON() ([]byte, error) {
	return json.Marshal(film.toJsonNum())
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

type FilmBookmarks struct {
	FilmsList     []Film `json:"bookmarks_list"`
	MoreAvailable bool   `json:"more_available"`
	FilmsTotal    int    `json:"films_total"`
	CurrentSort   string `json:"current_sort"`
	CurrentLimit  int    `json:"current_limit"`
	CurrentSkip   int    `json:"current_skip"`
}

type FilmPageInfo struct {
	FilmMain        *Film               `json:"film"`
	Reviews         FilmReviews         `json:"reviews"`
	Recommendations FilmRecommendations `json:"recommendations"`
	MyReview        Review              `json:"my_review"`
	Bookmarked      bool                `json:"bookmarked"`
}

type FilmPageInfoJson struct {
	FilmMain        FilmJson            `json:"film"`
	Reviews         FilmReviews         `json:"reviews"`
	Recommendations FilmRecommendations `json:"recommendations"`
	MyReview        Review              `json:"my_review"`
	Bookmarked      bool                `json:"bookmarked"`
}

func (filmPage *FilmPageInfo) CustomEasyJSON() ([]byte, error) {
	fpiJSON := FilmPageInfoJson{
		FilmMain:        filmPage.FilmMain.toJsonNum(),
		Reviews:         filmPage.Reviews,
		Recommendations: filmPage.Recommendations,
		MyReview:        filmPage.MyReview,
		Bookmarked:      filmPage.Bookmarked,
	}
	ret, _ := json.Marshal(fpiJSON)
	return ret, nil
}

type NewRate struct {
	Rating json.Number `json:"rating"`
}

type PostBookmarkResult struct {
	FilmID     uint64 `json:"film_id"`
	Bookmarked bool   `json:"bookmarked"`
}
//go:generate mockgen -destination=../film/repository/mock/db_mock.go  -package=mock 2021_2_MAMBa/internal/pkg/domain FilmRepository
type FilmRepository interface {
	GetFilm(id uint64) (Film, error)
	GetFilmReviews(id uint64, skip int, limit int) (FilmReviews, error)
	GetFilmRecommendations(id uint64, skip int, limit int) (FilmRecommendations, error)
	PostRating(id uint64, authorId uint64, rating float64) (float64, error)
	GetMyReview(id uint64, authorId uint64) (Review, error)
	CheckFilmBookmarked(userID uint64, filmID uint64) (bool, error)
	LoadUserBookmarkedFilmsID(userID uint64, skip int, limit int) ([]uint64, error)
	CountBookmarks(userID uint64) (int, error)
	BookmarkFilm(userID uint64, filmID uint64, bookmarked bool) error
	GetFilmsByMonthYear(month int, year int, limit int, skip int) (FilmList, error)
	GetGenres() (GenresList, error)
	GetFilmsByGenre(genreID uint64, limit int, skip int) (GenreFilmList, error)
	GetBanners() (BannersList, error)
	GetPopularFilms() (FilmList, error)
	GetRandomFilms (genre1 uint64, genre2 uint64, genre3 uint64, dateStart int, dateEnd int) (FilmList, error)
}

//go:generate mockgen -destination=../film/usecase/mock/usecase_mock.go  -package=mock 2021_2_MAMBa/internal/pkg/domain FilmUsecase
type FilmUsecase interface {
	GetFilm(userID, filmID uint64, skipReviews int, limitReviews int, skipRecommend int, limitRecommend int) (FilmPageInfo, error)
	PostRating(id uint64, authorId uint64, rating float64) (float64, error)
	LoadFilmReviews(id uint64, skip int, limit int) (FilmReviews, error)
	LoadFilmRecommendations(id uint64, skip int, limit int) (FilmRecommendations, error)
	LoadMyReview(id uint64, authorId uint64) (Review, error)
	LoadUserBookmarks(userID uint64, skip int, limit int) (FilmBookmarks, error)
	BookmarkFilm(userID uint64, filmID uint64, bookmarked bool) error
	GetFilmsByMonthYear(month int, year int, limit int, skip int) (FilmList, error)
	GetGenres() (GenresList, error)
	GetFilmsByGenre(genreID uint64, limit int, skip int) (GenreFilmList, error)
	GetBanners() (BannersList, error)
	GetPopularFilms() (FilmList, error)
	GetRandomFilms (genre1 uint64, genre2 uint64, genre3 uint64, dateStart int, dateEnd int) (FilmList, error)
}
