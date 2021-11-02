package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/film"
	"2021_2_MAMBa/internal/pkg/user"
	"encoding/binary"
	"github.com/jackc/pgx/pgtype"
	"golang.org/x/crypto/bcrypt"
	"math"
)

type dbUserRepository struct {
	dbm *database.DBManager
}

func NewUserRepository(manager *database.DBManager) domain.UserRepository {
	return &dbUserRepository{dbm: manager}
}

const (
	queryGetById              = "SELECT * FROM Profile WHERE User_ID = $1"
	queryGetByEmail           = "SELECT * FROM Profile WHERE email = $1"
	queryAddUser              = "INSERT INTO Profile(first_name, surname, email, password, picture_url, register_date) VALUES ($1, $2, $3, $4, $5, current_timestamp) RETURNING User_ID"
	queryCountBookmarksById   = "SELECT COUNT(*) FROM bookmark WHERE user_id = $1"
	queryCountSubscribersById = "SELECT COUNT(*) FROM subscription WHERE author_id = $1"
	queryCheckSubscription    = "SELECT COUNT(1) FROM subscription WHERE subscriber_id = $1 AND author_id = $2;"
	queryUpdProfile           = "UPDATE Profile SET first_name = $2, surname = $3, picture_url = $4, email = $5, gender = $6 WHERE user_id = $1"
	querySubscribe            = "INSERT INTO subscription VALUES ($1, $2) ON CONFLICT DO NOTHING"
	queryGetAuthorName        = "SELECT first_name, surname, picture_url FROM profile WHERE user_id = $1"
	queryGetFilmShort         = "SELECT title, title_original, poster_url FROM FILM WHERE Film_ID = $1"
	queryCountFilmReviews     = "SELECT COUNT(*) FROM Review WHERE author_id = $1 AND (NOT type = 0)"
	queryGetReviewByUserID    = "SELECT * FROM review WHERE author_id = $1 AND (NOT type = 0) LIMIT $2 OFFSET $3"
)

func (ur *dbUserRepository) GetUserByEmail(email string) (domain.User, error) {
	result, err := ur.dbm.Query(queryGetByEmail, email)
	if err != nil {
		return domain.User{}, user.ErrorInternalServer
	}
	if len(result) == 0 {
		return domain.User{}, user.ErrorNoUser
	}
	raw := result[0]
	found := domain.User{
		ID:             binary.BigEndian.Uint64(raw[0]),
		FirstName:      string(raw[1]),
		Surname:        string(raw[2]),
		Email:          string(raw[3]),
		Password:       string(raw[4]),
		PasswordRepeat: "",
		ProfilePic:     string(raw[5]),
	}
	return found, nil
}

func (ur *dbUserRepository) GetUserById(id uint64) (domain.User, error) {
	result, err := ur.dbm.Query(queryGetById, id)
	if err != nil {
		return domain.User{}, user.ErrorInternalServer
	}
	if len(result) == 0 {
		return domain.User{}, user.ErrorNoUser
	}
	raw := result[0]
	found := domain.User{
		ID:             binary.BigEndian.Uint64(raw[0]),
		FirstName:      string(raw[1]),
		Surname:        string(raw[2]),
		Email:          string(raw[3]),
		Password:       string(raw[4]),
		PasswordRepeat: "",
		ProfilePic:     string(raw[5]),
	}
	return found, nil
}

func (ur *dbUserRepository) AddUser(us *domain.User) (uint64, error) {
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, user.ErrorInternalServer
	}
	us.Password = string(passwordByte)
	us.ProfilePic = domain.BasePicture
	result, err := ur.dbm.Query(queryAddUser, us.FirstName, us.Surname, us.Email, us.Password, us.ProfilePic)

	us.ID = binary.BigEndian.Uint64(result[0][0])

	return us.ID, nil
}

func (ur *dbUserRepository) GetProfileById(whoAskID, id uint64) (domain.Profile, error) {
	result, err := ur.dbm.Query(queryGetById, id)
	if err != nil {
		return domain.Profile{}, err
	}
	if len(result) == 0 {
		return domain.Profile{}, user.ErrorNoUser
	}

	rawRow := result[0]

	timeBuffer := pgtype.Timestamp{}
	err = timeBuffer.DecodeBinary(nil, result[0][7])
	if err != nil {
		return domain.Profile{}, err
	}

	resultBookmarks, err := ur.dbm.Query(queryCountBookmarksById, id)
	if err != nil {
		return domain.Profile{}, err
	}

	resultSubscribers, err := ur.dbm.Query(queryCountSubscribersById, id)
	if err != nil {
		return domain.Profile{}, err
	}

	amSubscribed, err := ur.CheckSubscription(whoAskID, id)
	if err != nil {
		return domain.Profile{}, err
	}

	found := domain.Profile{
		ID:            binary.BigEndian.Uint64(rawRow[0]),
		FirstName:     string(rawRow[1]),
		Surname:       string(rawRow[2]),
		PictureUrl:    string(rawRow[5]),
		Email:         string(rawRow[3]),
		Gender:        string(rawRow[6]),
		RegisterDate:  timeBuffer.Time,
		SubCount:      int(binary.BigEndian.Uint32(resultSubscribers[0][0])),
		BookmarkCount: int(binary.BigEndian.Uint32(resultBookmarks[0][0])),
		AmSubscribed:  amSubscribed,
	}
	return found, nil
}

func (ur *dbUserRepository) CheckSubscription(src, dst uint64) (bool, error) {
	result, err := ur.dbm.Query(queryCheckSubscription, src, dst)
	if err != nil {
		return false, err
	}

	var count = int(binary.BigEndian.Uint32(result[0][0]))
	if count == 0 {
		return false, nil
	} else if count == 1 {
		return true, nil
	} else {
		return false, user.ErrorInternalServer
	}
}

func (ur *dbUserRepository) UpdateProfile(profile domain.Profile) (domain.Profile, error) {
	_, err := ur.dbm.Query(queryUpdProfile, profile.ID, profile.FirstName,
		profile.Surname, profile.PictureUrl, profile.Email, profile.Gender)
	if err != nil {
		return domain.Profile{}, err
	}
	updated, err := ur.GetProfileById(profile.ID, profile.ID)
	if err != nil {
		return domain.Profile{}, err
	}
	return updated, err
}

func (ur *dbUserRepository) CreateSubscription(src, dst uint64) (domain.Profile, error) {
	_, err := ur.dbm.Query(querySubscribe, src, dst)
	if err != nil {
		return domain.Profile{}, err
	}
	updated, err := ur.GetProfileById(src, src)
	if err != nil {
		return domain.Profile{}, err
	}
	return updated, err
}

func (ur *dbUserRepository) LoadUserReviews(id uint64, skip int, limit int) (domain.FilmReviews, error) {
	result, err := ur.dbm.Query(queryCountFilmReviews, id)
	if err != nil {
		return domain.FilmReviews{}, user.ErrorInternalServer
	}
	dbSizeRaw := binary.BigEndian.Uint64(result[0][0])
	dbSize := int(dbSizeRaw)
	if skip >= dbSize {
		return domain.FilmReviews{}, film.ErrorSkip
	}

	moreAvailable := skip+limit < dbSize

	result, err = ur.dbm.Query(queryGetReviewByUserID, id, limit, skip)
	if err != nil {
		return domain.FilmReviews{}, err
	}
	reviewList := make([]domain.Review, 0)
	for i := range result {
		timeBuffer := pgtype.Timestamp{}
		err = timeBuffer.DecodeBinary(nil, result[i][6])
		if err != nil {
			return domain.FilmReviews{}, err
		}
		review := domain.Review{
			Id:         binary.BigEndian.Uint64(result[i][0]),
			FilmId:     binary.BigEndian.Uint64(result[i][1]),
			ReviewText: string(result[i][3]),
			ReviewType: int(binary.BigEndian.Uint32(result[i][4])),
			Stars:      math.Float64frombits(binary.BigEndian.Uint64(result[i][5])),
			Date:       timeBuffer.Time,
		}
		filmId := binary.BigEndian.Uint64(result[i][1])
		authId := binary.BigEndian.Uint64(result[i][2])
		result1, err := ur.dbm.Query(queryGetAuthorName, authId)
		if err != nil {
			return domain.FilmReviews{}, err
		}
		review.AuthorName = string(result1[0][0]) + " " + string(result1[0][1])
		review.AuthorPictureUrl = string(result1[0][2])
		result1, err = ur.dbm.Query(queryGetFilmShort, filmId)
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
		CurrentSkip:   skip + limit,
	}
	return reviews, nil
}
