package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
)

type dbUserRepository struct {
	dbm *database.DBManager
}

func NewUserRepository(manager *database.DBManager) domain.UserRepository {
	return &dbUserRepository{dbm: manager}
}

const (
	queryGetById              = "SELECT * FROM profile WHERE user_ID = $1"
	queryGetByEmail           = "SELECT * FROM profile WHERE email = $1"
	queryAddUser              = "INSERT INTO profile(first_name, surname, email, password, picture_url, register_date) VALUES ($1, $2, $3, $4, $5, current_timestamp) RETURNING user_ID"
	queryCountBookmarksById   = "SELECT COUNT(*) FROM bookmark WHERE user_id = $1"
	queryCountSubscribersById = "SELECT COUNT(*) FROM subscription WHERE author_id = $1"
	queryCheckSubscription    = "SELECT COUNT(1) FROM subscription WHERE subscriber_id = $1 AND author_id = $2;"
	queryUpdProfile           = "UPDATE profile SET first_name = $2, surname = $3, picture_url = $4, email = $5, gender = $6 WHERE user_id = $1"
	querySubscribe            = "INSERT INTO subscription VALUES ($1, $2) ON CONFLICT DO NOTHING"
	queryGetAuthorName        = "SELECT first_name, surname, picture_url FROM profile WHERE user_id = $1"
	queryGetFilmShort         = "SELECT title, title_original, poster_url FROM FILM WHERE film_ID = $1"
	queryCountFilmReviews     = "SELECT COUNT(*) FROM review WHERE author_id = $1"
	queryGetReviewByUserID    = "SELECT * FROM review WHERE author_id = $1 LIMIT $2 OFFSET $3"
	queryUpdAvatarByUsID      = "UPDATE profile SET picture_url = $2 WHERE user_id = $1"
)

func (ur *dbUserRepository) GetUserByEmail(email string) (domain.User, error) {
	result, err := ur.dbm.Query(queryGetByEmail, email)
	if err != nil {
		return domain.User{}, customErrors.ErrorInternalServer
	}
	if len(result) == 0 {
		return domain.User{}, customErrors.ErrorNoUser
	}
	raw := result[0]
	found := domain.User{
		ID:             cast.ToUint64(raw[0]),
		FirstName:      cast.ToString(raw[1]),
		Surname:        cast.ToString(raw[2]),
		Email:          cast.ToString(raw[3]),
		Password:       cast.ToString(raw[4]),
		PasswordRepeat: "",
		ProfilePic:     cast.ToString(raw[5]),
	}
	return found, nil
}

func (ur *dbUserRepository) GetUserById(id uint64) (domain.User, error) {
	result, err := ur.dbm.Query(queryGetById, id)
	if err != nil {
		return domain.User{}, customErrors.ErrorInternalServer
	}
	if len(result) == 0 {
		return domain.User{}, customErrors.ErrorNoUser
	}
	raw := result[0]
	found := domain.User{
		ID:             cast.ToUint64(raw[0]),
		FirstName:      cast.ToString(raw[1]),
		Surname:        cast.ToString(raw[2]),
		Email:          cast.ToString(raw[3]),
		Password:       cast.ToString(raw[4]),
		PasswordRepeat: "",
		ProfilePic:     cast.ToString(raw[5]),
	}
	return found, nil
}

func (ur *dbUserRepository) AddUser(us *domain.User) (uint64, error) {
	result, err := ur.dbm.Query(queryAddUser, us.FirstName, us.Surname, us.Email, us.Password, us.ProfilePic)
	if err != nil {
		return 0, err
	}
	us.ID = cast.ToUint64(result[0][0])
	return us.ID, nil
}

func (ur *dbUserRepository) GetProfileById(whoAskID, id uint64) (domain.Profile, error) {
	result, err := ur.dbm.Query(queryGetById, id)
	if err != nil {
		return domain.Profile{}, err
	}
	if len(result) == 0 {
		return domain.Profile{}, customErrors.ErrorNoUser
	}

	rawRow := result[0]
	timeBuffer, err := cast.TimestampToString(result[0][7])
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
		ID:            cast.ToUint64(rawRow[0]),
		FirstName:     cast.ToString(rawRow[1]),
		Surname:       cast.ToString(rawRow[2]),
		PictureUrl:    cast.ToString(rawRow[5]),
		Email:         cast.ToString(rawRow[3]),
		Gender:        cast.ToString(rawRow[6]),
		RegisterDate:  timeBuffer,
		SubCount:      int(cast.ToUint64(resultSubscribers[0][0])),
		BookmarkCount: int(cast.ToUint64(resultBookmarks[0][0])),
		AmSubscribed:  amSubscribed,
	}
	return found, nil
}

func (ur *dbUserRepository) CheckSubscription(src, dst uint64) (bool, error) {
	result, err := ur.dbm.Query(queryCheckSubscription, src, dst)
	if err != nil {
		return false, err
	}
	if len(result) == 0 {
		return false, nil
	}

	count := cast.ToUint64(result[0][0])
	if count == 0 {
		return false, nil
	} else if count == 1 {
		return true, nil
	} else {
		return false, customErrors.ErrorInternalServer
	}
}

func (ur *dbUserRepository) UpdateProfile(profile domain.Profile) (domain.Profile, error) {
	err := ur.dbm.Execute(queryUpdProfile, profile.ID, profile.FirstName,
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

func (ur *dbUserRepository) UpdateAvatar(clientID uint64, url string) (domain.Profile, error) {
	_, err := ur.dbm.Query(queryUpdAvatarByUsID, clientID, url)
	if err != nil {
		return domain.Profile{}, err
	}
	updated, err := ur.GetProfileById(clientID, clientID)
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
		return domain.FilmReviews{}, customErrors.ErrorInternalServer
	}
	dbSizeRaw := cast.ToUint64(result[0][0])
	dbSize := int(dbSizeRaw)
	if skip >= dbSize && skip != 0 {
		return domain.FilmReviews{}, customErrors.ErrorSkip
	}

	moreAvailable := skip+limit < dbSize

	result, err = ur.dbm.Query(queryGetReviewByUserID, id, limit, skip)
	if err != nil {
		return domain.FilmReviews{}, err
	}
	reviewList := make([]domain.Review, 0)
	for i := range result {
		timeBuffer, err := cast.TimestampToString(result[i][6])
		if err != nil {
			return domain.FilmReviews{}, err
		}
		review := domain.Review{
			Id:         cast.ToUint64(result[i][0]),
			FilmId:     cast.ToUint64(result[i][1]),
			AuthorId:   cast.ToUint64(result[i][2]),
			ReviewText: cast.ToString(result[i][3]),
			ReviewType: int(cast.ToUint32(result[i][4])),
			Stars:      cast.ToFloat64(result[i][5]),
			Date:       timeBuffer,
		}
		filmId := cast.ToUint64(result[i][1])
		authId := cast.ToUint64(result[i][2])
		result1, err := ur.dbm.Query(queryGetAuthorName, authId)
		if err != nil {
			return domain.FilmReviews{}, err
		}
		review.AuthorName = cast.ToString(result1[0][0]) + " " + cast.ToString(result1[0][1])
		review.AuthorPictureUrl = cast.ToString(result1[0][2])
		result1, err = ur.dbm.Query(queryGetFilmShort, filmId)
		if err != nil {
			return domain.FilmReviews{}, err
		}
		review.FilmTitleRu = cast.ToString(result1[0][0])
		review.FilmTitleOriginal = cast.ToString(result1[0][1])
		review.FilmPictureUrl = cast.ToString(result1[0][2])
		reviewList = append(reviewList, review)
	}
	reviews := domain.FilmReviews{
		ReviewList:    reviewList,
		MoreAvailable: moreAvailable,
		ReviewTotal:   dbSize,
		CurrentSort:   "",
		CurrentLimit:  limit,
		CurrentSkip:   skip + limit,
	}
	return reviews, nil
}
