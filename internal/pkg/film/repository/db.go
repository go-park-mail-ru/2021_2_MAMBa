package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"time"
)

type dbFilmRepository struct {
	dbm *database.DBManager
}

func NewFilmRepository(manager *database.DBManager) domain.FilmRepository {
	return &dbFilmRepository{dbm: manager}
}

const (
	queryCountFilmReviews         = "SELECT COUNT(*) FROM Review WHERE Film_ID = $1 AND (NOT type = 0)"
	queryCountFilmRecommendations = "SELECT COUNT(*) FROM recommended WHERE Film_ID = $1"
	queryCountUserBookmarks       = "SELECT COUNT(*) FROM bookmark b WHERE user_id = $1"
	queryCheckFilmBookmarked      = "SELECT COUNT(*) FROM bookmark b WHERE user_id = $1 AND film_id = $2"
	queryGetFilmId                = "SELECT * FROM film WHERE film_id = $1"
	queryGetFilmDirScr            = "SELECT film.*, p.person_id,p.name_en,p.name_rus,p.picture_url,p.career, p1.person_id,p1.name_en,p1.name_rus,p1.picture_url,p1.career FROM film JOIN person p on film.director = p.person_id JOIN person p1 on film.screenwriter = p1.person_id WHERE film_id = $1"
	queryGetFilmCountries         = "SELECT country.country_name FROM country JOIN countryproduction c ON country.country_id = c.country_id WHERE c.Film_ID = $1"
	queryGetFilmGenres            = "SELECT genre.* FROM genre JOIN filmgenres f on genre.genre_id = f.genre_id WHERE f.film_id = $1"
	queryGetFilmCast              = "SELECT p.person_id,p.name_en,p.name_rus,p.picture_url,p.career  FROM person p JOIN filmcast f on p.person_id = f.person_id WHERE f.film_id = $1"
	queryGetFilmReviews           = "SELECT review.*, p.first_name, p.surname, p.picture_url FROM review join profile p on p.user_id = review.author_id WHERE film_id = $1 AND (NOT type = 0) ORDER BY review.review_date DESC LIMIT $2 OFFSET $3"
	queryGetFilmRecommendations   = "SELECT f.film_id, f.title, f.poster_url FROM recommended r join film f on f.film_id = r.recommended_id WHERE r.film_id = $1 LIMIT $2 OFFSET $3"
	queryGetFilmRating            = "SELECT AVG(stars) FROM review WHERE film_id = $1AND (NOT stars = 0)"
	queryPostRating               = "INSERT INTO  review (review_id, film_id, author_id, review_text, type, stars, review_date) VALUES (DEFAULT, $1, $2, '', $3, $4, $5)"
	queryGetReviewByAuthor        = "SELECT * FROM review WHERE film_id = $1 AND author_id = $2"
	queryUpdateRating             = "UPDATE review SET stars = $1 WHERE film_id = $2 AND author_id = $3"
	queryGetAuthorName            = "SELECT first_name, surname, picture_url FROM profile WHERE user_id = $1"
	queryGetFilmShort             = "SELECT title, title_original, poster_url FROM FILM WHERE Film_ID = $1"
	queryGetBookmarksByUserID     = "SELECT b.film_id FROM bookmark b WHERE user_id = $1 LIMIT $2 OFFSET $3"
	queryAddBookmark              = "INSERT INTO bookmark (film_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING"
	queryDeleteBookmark           = "DELETE FROM bookmark WHERE film_id = $1 AND user_id = $2"
	queryGetFilmsByMonthYear      = "SELECT film_id, title, title_original, release_year, description, poster_url, premiere_ru FROM film WHERE EXTRACT(MONTH FROM premiere_ru) = $1 AND EXTRACT(YEAR FROM premiere_ru) = $2 LIMIT $3 OFFSET $4"
	queryCountFilmsByMonthYear    = "SELECT COUNT(*) FROM film WHERE EXTRACT(MONTH FROM premiere_ru) = $1 AND EXTRACT(YEAR FROM premiere_ru) = $2"
	queryGetGenres                = "SELECT genre_id, genre_name, picture_url FROM genre LIMIT 15"
	queryGetFilmsByGenreID        = "SELECT f.film_id, f.title, f.title_original, f.release_year, f.description, f.poster_url, f.premiere_ru FROM film f JOIN filmgenres g on f.film_id = g.film_id WHERE genre_id = $1 LIMIT $2 OFFSET $3"
	queryCountFilmsByGenreID      = "SELECT COUNT(*) FROM film f JOIN filmgenres g on f.film_id = g.film_id WHERE genre_id = $1"
	queryGetGenreName             = "SELECT genre_name FROM genre WHERE genre_id = $1"
)

func (fr *dbFilmRepository) GetGenres() (domain.GenresList, error) {
	result, err := fr.dbm.Query(queryGetGenres)
	if err != nil {
		return domain.GenresList{}, customErrors.ErrorInternalServer
	}

	bufferGenres := make([]domain.Genre, 0)
	for i := range result {
		genre := domain.Genre{
			Id:         uint64(cast.ToUint32(result[i][0])),
			Name:       cast.ToString(result[i][1]),
			PictureURL: cast.ToString(result[i][2]),
		}
		bufferGenres = append(bufferGenres, genre)
	}
	return domain.GenresList{GenresList: bufferGenres}, nil
}

func (fr *dbFilmRepository) GetFilmsByGenre(genreID uint64, limit int, skip int) (domain.GenreFilmList, error) {
	result, err := fr.dbm.Query(queryCountFilmsByGenreID, genreID)
	if err != nil {
		return domain.GenreFilmList{}, customErrors.ErrorInternalServer
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	if skip >= dbSize && skip != 0 {
		return domain.GenreFilmList{}, customErrors.ErrorSkip
	}
	moreAvailable := skip+limit < dbSize

	result, err = fr.dbm.Query(queryGetFilmsByGenreID, genreID, limit, skip)
	if err != nil {
		return domain.GenreFilmList{}, err
	}

	bufferFilms := make([]domain.Film, 0)
	for i := range result {
		film := domain.Film{
			Id:            cast.ToUint64(result[i][0]),
			Title:         cast.ToString(result[i][1]),
			TitleOriginal: cast.ToString(result[i][2]),
			ReleaseYear:   int(cast.ToUint32(result[i][3])),
			Description:   cast.ToString(result[i][4]),
			PosterUrl:     cast.ToString(result[i][5]),
		}
		dateString, err := cast.DateToStringUnderscore(result[i][6])
		if err != nil {
			return domain.GenreFilmList{}, err
		}
		film.PremiereRu = dateString
		bufferFilms = append(bufferFilms, film)
	}

	filmList := domain.FilmList{
		FilmList:      bufferFilms,
		MoreAvailable: moreAvailable,
		FilmTotal:     dbSize,
		CurrentLimit:  limit,
		CurrentSkip:   skip + limit,
	}

	result, err = fr.dbm.Query(queryGetGenreName, genreID)
	if err != nil {
		return domain.GenreFilmList{}, err
	}
	genreName := cast.ToString(result[0][0])

	genreFilmList := domain.GenreFilmList{
		Id:        genreID,
		Name:      genreName,
		FilmsList: filmList,
	}
	return genreFilmList, nil
}

func (fr *dbFilmRepository) BookmarkFilm(userID uint64, filmID uint64, bookmarked bool) error {
	var err error
	if bookmarked {
		_, err = fr.dbm.Query(queryAddBookmark, filmID, userID)
	} else {
		_, err = fr.dbm.Query(queryDeleteBookmark, filmID, userID)
	}
	if err != nil {
		return err
	}
	return nil
}

func (fr *dbFilmRepository) GetFilmsByMonthYear(month int, year int, limit int, skip int) (domain.FilmList, error) {
	result, err := fr.dbm.Query(queryCountFilmsByMonthYear, month, year)
	if err != nil {
		return domain.FilmList{}, customErrors.ErrorInternalServer
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	if skip >= dbSize && skip != 0 {
		return domain.FilmList{}, customErrors.ErrorSkip
	}
	moreAvailable := skip+limit < dbSize

	result, err = fr.dbm.Query(queryGetFilmsByMonthYear, month, year, limit, skip)
	if err != nil {
		return domain.FilmList{}, err
	}

	bufferFilms := make([]domain.Film, 0)
	for i := range result {
		film := domain.Film{
			Id:            cast.ToUint64(result[i][0]),
			Title:         cast.ToString(result[i][1]),
			TitleOriginal: cast.ToString(result[i][2]),
			ReleaseYear:   int(cast.ToUint32(result[i][3])),
			Description:   cast.ToString(result[i][4]),
			PosterUrl:     cast.ToString(result[i][5]),
		}
		dateString, err := cast.DateToStringUnderscore(result[i][6])
		if err != nil {
			return domain.FilmList{}, err
		}
		film.PremiereRu = dateString
		bufferFilms = append(bufferFilms, film)
	}
	filmList := domain.FilmList{
		FilmList:      bufferFilms,
		MoreAvailable: moreAvailable,
		FilmTotal:     dbSize,
		CurrentLimit:  limit,
		CurrentSkip:   skip + limit,
	}
	return filmList, nil
}

func (fr *dbFilmRepository) LoadUserBookmarkedFilmsID(userID uint64, skip int, limit int) ([]uint64, error) {
	result, err := fr.dbm.Query(queryGetBookmarksByUserID, userID, limit, skip)
	filmIdxList := make([]uint64, 0)
	if err != nil {
		return filmIdxList, err
	}
	for i := range result {
		filmIdxList = append(filmIdxList, cast.ToUint64(result[i][0]))
	}
	return filmIdxList, nil
}

func (fr *dbFilmRepository) CheckFilmBookmarked(userID uint64, filmID uint64) (bool, error) {
	result, err := fr.dbm.Query(queryCheckFilmBookmarked, userID, filmID)
	if err != nil {
		return false, err
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	if dbSize > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (fr *dbFilmRepository) GetFilm(id uint64) (domain.Film, error) {
	result, err := fr.dbm.Query(queryGetFilmDirScr, id)
	if err != nil || len(result) == 0 {
		return domain.Film{}, err
	}
	raw := result[0]
	film := domain.Film{
		Id:              cast.ToUint64(raw[0]),
		Title:           cast.ToString(raw[1]),
		TitleOriginal:   cast.ToString(raw[2]),
		Rating:          cast.ToFloat64(raw[3]),
		Description:     cast.ToString(raw[4]),
		TotalRevenue:    cast.ToString(raw[7]),
		PosterUrl:       cast.ToString(raw[5]),
		TrailerUrl:      cast.ToString(raw[6]),
		ContentType:     cast.ToString(raw[12]),
		ReleaseYear:     int(cast.ToUint32(raw[8])),
		Duration:        int(cast.ToUint32(raw[9])),
		OriginCountries: nil,
		Cast:            nil,
		Director: domain.Person{
			Id:         cast.ToUint64(raw[14]),
			NameEn:     cast.ToString(raw[15]),
			NameRus:    cast.ToString(raw[16]),
			PictureUrl: cast.ToString(raw[17]),
			Career:     []string{cast.ToString(raw[18])},
		},
		Screenwriter: domain.Person{
			Id:         cast.ToUint64(raw[19]),
			NameEn:     cast.ToString(raw[20]),
			NameRus:    cast.ToString(raw[21]),
			PictureUrl: cast.ToString(raw[22]),
			Career:     []string{cast.ToString(raw[23])},
		},
		Genres: nil,
	}

	dateString, err := cast.DateToStringUnderscore(raw[13])
	if err != nil {
		return domain.Film{}, err
	}
	film.PremiereRu = dateString

	result, err = fr.dbm.Query(queryGetFilmRating, id)
	if err != nil {
		return domain.Film{}, err
	}
	film.Rating = cast.ToFloat64(result[0][0])
	result, err = fr.dbm.Query(queryGetFilmCountries, id)
	if err != nil || len(result) == 0 {
		return domain.Film{}, err
	}
	bufferCountries := make([]string, 0)
	for i := range result {
		bufferCountries = append(bufferCountries, cast.ToString(result[i][0]))
	}
	film.OriginCountries = bufferCountries

	result, err = fr.dbm.Query(queryGetFilmGenres, id)
	if err != nil {
		return domain.Film{}, err
	}
	if len(result) == 0 {
		return domain.Film{}, err
	}
	bufferGenres := make([]domain.Genre, 0)
	for i := range result {
		bufferGenres = append(bufferGenres, domain.Genre{
			Id:   uint64(cast.ToUint32(result[i][0])),
			Name: cast.ToString(result[i][1]),
		})
	}
	film.Genres = bufferGenres

	result, err = fr.dbm.Query(queryGetFilmCast, id)
	if err != nil {
		return domain.Film{}, err
	}
	if len(result) == 0 {
		return domain.Film{}, err
	}
	bufferCast := make([]domain.Person, 0)
	for i := range result {
		bufferCast = append(bufferCast, domain.Person{
			Id:         cast.ToUint64(result[i][0]),
			NameEn:     cast.ToString(result[i][1]),
			NameRus:    cast.ToString(result[i][2]),
			PictureUrl: cast.ToString(result[i][3]),
			Career:     []string{cast.ToString(result[i][4])},
		})
	}
	film.Cast = bufferCast
	return film, nil
}

func (fr *dbFilmRepository) GetFilmReviews(id uint64, skip int, limit int) (domain.FilmReviews, error) {
	result, err := fr.dbm.Query(queryCountFilmReviews, id)
	if err != nil {
		return domain.FilmReviews{}, customErrors.ErrorInternalServer
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	if skip >= dbSize && skip != 0 {
		return domain.FilmReviews{}, customErrors.ErrorSkip
	}
	moreAvailable := skip+limit < dbSize
	result, err = fr.dbm.Query(queryGetFilmReviews, id, limit, skip)
	if err != nil {
		return domain.FilmReviews{}, customErrors.ErrorInternalServer
	}
	reviews := make([]domain.Review, 0)
	for i := range result {
		timestamp, err := cast.TimestampToString(result[i][6])
		if err != nil {
			return domain.FilmReviews{}, err
		}
		reviews = append(reviews, domain.Review{
			Id:               cast.ToUint64(result[i][0]),
			FilmId:           cast.ToUint64(result[i][1]),
			AuthorName:       cast.ToString(result[i][7]) + " " + cast.ToString(result[i][8]),
			ReviewText:       cast.ToString(result[i][3]),
			AuthorPictureUrl: cast.ToString(result[i][9]),
			ReviewType:       int(cast.ToUint32(result[i][4])),
			Stars:            cast.ToFloat64(result[i][5]),
			Date:             timestamp,
		})
	}
	reviewsList := domain.FilmReviews{
		ReviewList:    reviews,
		MoreAvailable: moreAvailable,
		ReviewTotal:   dbSize,
		CurrentSort:   "",
		CurrentLimit:  limit,
		CurrentSkip:   skip + limit,
	}
	return reviewsList, nil
}

func (fr *dbFilmRepository) GetFilmRecommendations(id uint64, skip int, limit int) (domain.FilmRecommendations, error) {
	result, err := fr.dbm.Query(queryCountFilmRecommendations, id)
	if err != nil {
		return domain.FilmRecommendations{}, customErrors.ErrorInternalServer
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	if skip >= dbSize && skip != 0 {
		return domain.FilmRecommendations{}, customErrors.ErrorSkip
	}

	moreAvailable := skip+limit < dbSize
	result, err = fr.dbm.Query(queryGetFilmRecommendations, id, limit, skip)
	reviews := make([]domain.Film, 0)
	for i := range result {
		reviews = append(reviews, domain.Film{
			Id:        cast.ToUint64(result[i][0]),
			Title:     cast.ToString(result[i][1]),
			PosterUrl: cast.ToString(result[i][2]),
		})
	}
	reviewsList := domain.FilmRecommendations{
		RecommendationList:  reviews,
		MoreAvailable:       moreAvailable,
		RecommendationTotal: dbSize,
		CurrentLimit:        limit,
		CurrentSkip:         skip + limit,
	}
	return reviewsList, nil
}

func (fr *dbFilmRepository) CountBookmarks(userID uint64) (int, error) {
	result, err := fr.dbm.Query(queryCountUserBookmarks, userID)
	if err != nil {
		return 0, customErrors.ErrorInternalServer
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	return dbSize, nil
}

func (fr *dbFilmRepository) PostRating(id uint64, author_id uint64, rating float64) (float64, error) {
	result, err := fr.dbm.Query(queryGetReviewByAuthor, id, author_id)
	if err != nil {
		return 0, err
	}
	if len(result) == 0 {
		err = fr.dbm.Execute(queryPostRating, id, author_id, 0, rating, time.Now())
		if err != nil {
			return 0, err
		}
	} else {
		err = fr.dbm.Execute(queryUpdateRating, rating, id, author_id)
		if err != nil {
			return 0, err
		}
	}
	result, err = fr.dbm.Query(queryGetFilmRating, id)
	if err != nil {
		return 0, err
	}
	newRating := cast.ToFloat64(result[0][0])
	return newRating, nil
}

func (fr *dbFilmRepository) GetMyReview(id uint64, author_id uint64) (domain.Review, error) {
	result, err := fr.dbm.Query(queryGetReviewByAuthor, id, author_id)
	if err != nil {
		return domain.Review{}, err
	}
	if len(result) == 0 {
		return domain.Review{}, customErrors.ErrorNoReviewForFilm
	}
	timestamp, err := cast.TimestampToString(result[0][6])
	urtype := cast.ToUint32(result[0][4])
	rtype := int(urtype)
	if err != nil {
		return domain.Review{}, err
	}
	review := domain.Review{
		Id:         cast.ToUint64(result[0][0]),
		FilmId:     cast.ToUint64(result[0][1]),
		ReviewText: cast.ToString(result[0][3]),
		ReviewType: rtype,
		Stars:      cast.ToFloat64(result[0][5]),
		Date:       timestamp,
	}
	result, err = fr.dbm.Query(queryGetAuthorName, author_id)
	if err != nil {
		return domain.Review{}, err
	}
	review.AuthorName = cast.ToString(result[0][0]) + " " + cast.ToString(result[0][1])
	review.AuthorPictureUrl = cast.ToString(result[0][2])
	result, err = fr.dbm.Query(queryGetFilmShort, id)
	if err != nil {
		return domain.Review{}, err
	}
	review.FilmTitleRu = cast.ToString(result[0][0])
	review.FilmTitleOriginal = cast.ToString(result[0][1])
	review.FilmPictureUrl = cast.ToString(result[0][2])

	return review, nil
}
