package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"encoding/binary"
	"github.com/jackc/pgx/pgtype"
	"math"
	"time"
)

type dbFilmRepository struct {
	dbm *database.DBManager
}

func NewFilmRepository(manager *database.DBManager) domain.FilmRepository {
	return &dbFilmRepository{dbm: manager}
}

const (
	queryCountFilmReviews = "SELECT COUNT(*) FROM Review WHERE Film_ID = $1 AND (NOT type = 0)"
	queryCountFilmRecommendations = "SELECT COUNT(*) FROM recommended WHERE Film_ID = $1"
	queryGetFilmId = "SELECT * FROM film WHERE film_id = $1"
	queryGetFilmDirScr = "SELECT film.*, p.person_id,p.name_en,p.name_rus,p.picture_url,p.career, p1.person_id,p1.name_en,p1.name_rus,p1.picture_url,p1.career FROM film JOIN person p on film.director = p.person_id JOIN person p1 on film.screenwriter = p1.person_id WHERE film_id = $1"
	queryGetFilmCountries = "SELECT country.country_name FROM country JOIN countryproduction c ON country.country_id = c.country_id WHERE c.Film_ID = $1"
	queryGetFilmGenres = "SELECT genre.* FROM genre JOIN filmgenres f on genre.genre_id = f.genre_id WHERE f.film_id = $1"
	queryGetFilmCast = "SELECT p.person_id,p.name_en,p.name_rus,p.picture_url,p.career  FROM person p JOIN filmcast f on p.person_id = f.person_id WHERE f.film_id = $1"
	queryGetFilmReviews = "SELECT review.*, p.first_name, p.surname, p.picture_url FROM review join profile p on p.user_id = review.author_id WHERE film_id = $1 AND (NOT type = 0) LIMIT $2 OFFSET $3"
	queryGetFilmRecommendations = "SELECT f.film_id, f.title, f.poster_url FROM recommended r join film f on f.film_id = r.recommended_id WHERE r.film_id = $1 LIMIT $2 OFFSET $3"
	queryGetFilmRating = "SELECT AVG(stars) FROM review WHERE film_id = $1AND (NOT stars = 0)"
	queryPostRating = "INSERT INTO  review VALUES (DEFAULT, $1, $2, '', $3, $4, $5)"
	queryGetReviewByAuthor = "SELECT * FROM review WHERE film_id = $1 AND author_id = $2"
	queryUpdateRating = "UPDATE review SET stars = $1 WHERE film_id = $2 AND author_id = $3"
	queryGetAuthorName = "SELECT first_name, surname, picture_url FROM profile WHERE user_id = $1"
	queryGetFilmShort = "SELECT title, title_original, poster_url FROM FILM WHERE Film_ID = $1"
)

func (fr *dbFilmRepository) GetFilm (id uint64) (domain.Film, error) {
	result, err := fr.dbm.Query(queryGetFilmDirScr, id)
	if err != nil {
		return domain.Film{}, err
	}
	if len(result) == 0 {
		return domain.Film{}, err
	}
	raw := result[0]
	film := domain.Film{
		Id:              binary.BigEndian.Uint64(raw[0]),
		Title:           string(raw[1]),
		TitleOriginal:   string(raw[2]),
		Rating:          math.Float64frombits(binary.BigEndian.Uint64(raw[3])),
		Description:     string(raw[4]),
		TotalRevenue:    string(raw[7]),
		PosterUrl:       string(raw[5]),
		TrailerUrl:      string(raw[6]),
		ContentType:     string(raw[12]),
		ReleaseYear:     int(binary.BigEndian.Uint32(raw[8])),
		Duration:        int(binary.BigEndian.Uint32(raw[9])),
		OriginCountries: nil,
		Cast:            nil,
		Director: domain.Person{
			Id: binary.BigEndian.Uint64(raw[13]),
			NameEn: string(raw[14]),
			NameRus:string(raw[15]),
			PictureUrl: string(raw[16]),
			Career: []string{string(raw[17])},
		},
		Screenwriter: domain.Person{
			Id: binary.BigEndian.Uint64(raw[18]),
			NameEn: string(raw[19]),
			NameRus:string(raw[20]),
			PictureUrl: string(raw[21]),
			Career: []string{string(raw[22])},
		},
		Genres:       nil,
	}
	result, err = fr.dbm.Query(queryGetFilmRating, id)
	if err != nil {
		return domain.Film{}, err
	}
	film.Rating = math.Float64frombits(binary.BigEndian.Uint64(result[0][0]))
	result, err = fr.dbm.Query(queryGetFilmCountries, id)
	if err != nil {
		return domain.Film{}, err
	}
	if len(result) == 0 {
		return domain.Film{}, err
	}
	bufferCountries := make([]string, 0)
	for i := range result {
		bufferCountries = append(bufferCountries, string(result[i][0]))
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
			Id:   uint64(binary.BigEndian.Uint32(result[i][0])),
			Name: string(result[i][1]),
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
			Id: binary.BigEndian.Uint64(result[i][0]),
			NameEn: string(result[i][1]),
			NameRus:string(result[i][2]),
			PictureUrl: string(result[i][3]),
			Career: []string{string(result[i][4])},
		})
	}
	film.Cast = bufferCast
	return film, nil
}
func (fr *dbFilmRepository) GetFilmReviews (id uint64, skip int, limit int) (domain.FilmReviews, error) {
	result, err := fr.dbm.Query(queryCountFilmReviews, id)
	if err != nil {
		return domain.FilmReviews{}, customErrors.ErrorInternalServer
	}
	dbSizeRaw := binary.BigEndian.Uint64(result[0][0])
	dbSize := int(dbSizeRaw)
	if skip >= dbSize {
		return domain.FilmReviews{}, customErrors.ErrorSkip
	}
	

	moreAvailable := skip+limit < dbSize
	result, err = fr.dbm.Query(queryGetFilmReviews, id, limit, skip)
	reviews := make([]domain.Review, 0)
	for i := range result {
		timeBuffer := pgtype.Timestamp{}
		err =timeBuffer.DecodeBinary(nil, result[i][6])
		if err != nil {
			return domain.FilmReviews{}, err
		}
		reviews = append(reviews, domain.Review{
			Id:                binary.BigEndian.Uint64(result[i][0]),
			FilmId:            binary.BigEndian.Uint64(result[i][1]),
			AuthorName:        string(result[i][7])+ " " +string(result[i][8]),
			ReviewText:        string(result[i][3]),
			AuthorPictureUrl:  string(result[i][9]),
			ReviewType:        int(binary.BigEndian.Uint32(result[i][4])),
			Stars:             math.Float64frombits(binary.BigEndian.Uint64(result[i][5])),
			Date:              timeBuffer.Time,
		})
	}
	reviewsList := domain.FilmReviews{
		ReviewList:    reviews,
		MoreAvaliable: moreAvailable,
		ReviewTotal:   dbSize,
		CurrentSort:   "",
		CurrentLimit:  limit,
		CurrentSkip:   skip+limit,
	}
	return reviewsList, nil
}
func (fr *dbFilmRepository) GetFilmRecommendations (id uint64, skip int, limit int) (domain.FilmRecommendations, error) {
	result, err := fr.dbm.Query(queryCountFilmRecommendations, id)
	if err != nil {
		return domain.FilmRecommendations{}, customErrors.ErrorInternalServer
	}
	dbSizeRaw := binary.BigEndian.Uint64(result[0][0])
	dbSize := int(dbSizeRaw)
	if skip >= dbSize {
		return domain.FilmRecommendations{}, customErrors.ErrorSkip
	}


	moreAvailable := skip+limit < dbSize
	result, err = fr.dbm.Query(queryGetFilmRecommendations, id, limit, skip)
	reviews := make([]domain.Film, 0)
	for i := range result {
		reviews = append(reviews, domain.Film{
			Id:                binary.BigEndian.Uint64(result[i][0]),
			Title: string(result[i][1]),
			PosterUrl: string(result[i][2]),
		})
	}
	reviewsList := domain.FilmRecommendations{
		RecommendationList:    reviews,
		MoreAvaliable: moreAvailable,
		RecommendationTotal:   dbSize,
		CurrentLimit:  limit,
		CurrentSkip:   skip+limit,
	}
	return reviewsList, nil
}

func  (fr *dbFilmRepository) PostRating (id uint64, author_id uint64, rating float64) (float64, error) {
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
	newRating := math.Float64frombits(binary.BigEndian.Uint64(result[0][0]))
	return newRating, nil
}


func (fr *dbFilmRepository) GetMyReview (id uint64, author_id uint64) (domain.Review, error) {
	result, err := fr.dbm.Query(queryGetReviewByAuthor, id, author_id)
	if err != nil {
		return domain.Review{}, err
	}
	if len(result) == 0 {
		return domain.Review{}, customErrors.ErrorNoReviewForFilm
	}
	timeBuffer := pgtype.Timestamp{}
	err =timeBuffer.DecodeBinary(nil, result[0][6])
	urtype := binary.BigEndian.Uint32(result[0][4])
	rtype := int(urtype)
	if err != nil {
		return domain.Review{}, err
	}
	review := domain.Review{
		Id:					binary.BigEndian.Uint64(result[0][0]),
		FilmId:            binary.BigEndian.Uint64(result[0][1]),
		ReviewText:        string(result[0][3]),
		ReviewType:        rtype,
		Stars:             math.Float64frombits(binary.BigEndian.Uint64(result[0][5])),
		Date:              timeBuffer.Time,
	}
	result, err = fr.dbm.Query(queryGetAuthorName, author_id)
	if err != nil {
		return domain.Review{}, err
	}
	review.AuthorName = string(result[0][0]) +" "+ string(result[0][1])
	review.AuthorPictureUrl = string(result[0][2])
	result, err = fr.dbm.Query(queryGetFilmShort, id)
	if err != nil {
		return domain.Review{}, err
	}
	review.FilmTitleRu = string(result[0][0])
	review.FilmTitleOriginal = string(result[0][1])
	review.FilmPictureUrl = string(result[0][2])

	return review, nil
}