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

type dbReviewRepository struct {
	dbm *database.DBManager
}

func NewReviewRepository(manager *database.DBManager) domain.ReviewRepository {
	return &dbReviewRepository{dbm: manager}
}
const (
	queryCountFilmReviews = "SELECT COUNT(*) FROM Review WHERE Film_ID = $1 AND (NOT type = 0)"
	queryGetReviewByID = "SELECT * FROM review WHERE review_id = $1"
	queryGetAuthorName = "SELECT first_name, surname, picture_url FROM profile WHERE user_id = $1"
	queryGetFilmShort = "SELECT title, title_original, poster_url FROM FILM WHERE Film_ID = $1"
	queryGetReviewByFilmIDEXCEPT = "SELECT * FROM review WHERE film_id = $1 AND (NOT review_id = $2) AND (NOT type = 0) LIMIT $3 OFFSET $4"
	querySearchReview = "SELECT * FROM review WHERE film_id = $1 AND author_id = $2"
	queryInsertReview = "INSERT INTO review VALUES (DEFAULT, $1, $2, $3, $4, $5, $6) RETURNING review_id"
	queryUpdateReview = "UPDATE review SET review_text = $1, type = $2 WHERE film_id = $3 AND author_id = $4 RETURNING review_id"
)

func (rr *dbReviewRepository) GetReview (id uint64) (domain.Review, error) {
	result, err := rr.dbm.Query(queryGetReviewByID, id)
	if err != nil {
		return domain.Review{}, err
	}
	timeBuffer := pgtype.Timestamp{}
	err =timeBuffer.DecodeBinary(nil, result[0][6])
	if err != nil {
		return domain.Review{}, err
	}
	review := domain.Review{
		Id:                binary.BigEndian.Uint64(result[0][0]),
		FilmId:            binary.BigEndian.Uint64(result[0][1]),
		ReviewText:        string(result[0][3]),
		ReviewType:        int(binary.BigEndian.Uint32(result[0][4])),
		Stars:             math.Float64frombits(binary.BigEndian.Uint64(result[0][5])),
		Date:              timeBuffer.Time,
	}
	filmId := binary.BigEndian.Uint64(result[0][1])
	authId := binary.BigEndian.Uint64(result[0][2])
	result, err = rr.dbm.Query(queryGetAuthorName, authId)
	if err != nil {
		return domain.Review{}, err
	}
	review.AuthorName = string(result[0][0]) +" "+ string(result[0][1])
	review.AuthorPictureUrl = string(result[0][2])
	result, err = rr.dbm.Query(queryGetFilmShort, filmId)
	if err != nil {
		return domain.Review{}, err
	}
	review.FilmTitleRu = string(result[0][0])
	review.FilmTitleOriginal = string(result[0][1])
	review.FilmPictureUrl = string(result[0][2])
	return review, nil
}

func (rr *dbReviewRepository)  PostReview (review domain.Review) (uint64, error) {
	result, err := rr.dbm.Query(querySearchReview, review.FilmId, review.AuthorId)
	if err != nil {
		return 0, customErrors.ErrorInternalServer
	}
	if len(result) == 0 {
		result, err = rr.dbm.Query(queryInsertReview, review.FilmId, review.AuthorId, review.ReviewText, review.ReviewType, 0, time.Now())
		if err != nil {
			return 0, customErrors.ErrorInternalServer
		}
	} else {
		result, err = rr.dbm.Query(queryUpdateReview, review.ReviewText, review.ReviewType, review.FilmId, review.AuthorId)
		if err != nil {
			return 0, customErrors.ErrorInternalServer
		}
	}
	return binary.BigEndian.Uint64(result[0][0]), nil
}

func (rr *dbReviewRepository)  LoadReviewsExcept(id uint64, film_id uint64, skip int, limit int) (domain.FilmReviews, error) {
	result, err := rr.dbm.Query(queryCountFilmReviews, film_id)
	if err != nil {
		return domain.FilmReviews{}, customErrors.ErrorInternalServer
	}
	dbSizeRaw := binary.BigEndian.Uint64(result[0][0])
	dbSize := int(dbSizeRaw)
	if skip >= dbSize {
		return domain.FilmReviews{}, customErrors.ErrorSkip
	}

	moreAvailable := skip+limit < dbSize

	result, err = rr.dbm.Query(queryGetReviewByFilmIDEXCEPT, film_id, id, limit, skip)
	if err != nil {
		return domain.FilmReviews{}, err
	}
	reviewList := make([]domain.Review, 0)
	for i := range result{
		timeBuffer := pgtype.Timestamp{}
		err =timeBuffer.DecodeBinary(nil, result[i][6])
		if err != nil {
			return domain.FilmReviews{}, err
		}
		review := domain.Review{
			Id:                binary.BigEndian.Uint64(result[i][0]),
			FilmId:            binary.BigEndian.Uint64(result[i][1]),
			ReviewText:        string(result[i][3]),
			ReviewType:        int(binary.BigEndian.Uint32(result[i][4])),
			Stars:             math.Float64frombits(binary.BigEndian.Uint64(result[i][5])),
			Date:              timeBuffer.Time,
		}
		filmId := binary.BigEndian.Uint64(result[i][1])
		authId := binary.BigEndian.Uint64(result[i][2])
		result1, err := rr.dbm.Query(queryGetAuthorName, authId)
		if err != nil {
			return domain.FilmReviews{}, err
		}
		review.AuthorName = string(result1[0][0]) + " " + string(result1[0][1])
		review.AuthorPictureUrl = string(result1[0][2])
		result1, err = rr.dbm.Query(queryGetFilmShort, filmId)
		if err != nil {
			return domain.FilmReviews{}, err
		}
		review.FilmTitleRu = string(result1[0][0])
		review.FilmTitleOriginal = string(result1[0][1])
		review.FilmPictureUrl = string(result1[0][2])
		reviewList = append(reviewList, review)
	}
	reviews := domain.FilmReviews{
		ReviewList:    reviewList,
		MoreAvaliable: moreAvailable,
		ReviewTotal:   dbSize,
		CurrentSort:   "",
		CurrentLimit:  limit,
		CurrentSkip:   skip+limit,
	}
	return reviews, nil
}