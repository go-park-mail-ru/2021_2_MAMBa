package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	userErrors "2021_2_MAMBa/internal/pkg/user"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(u domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: u,
	}
}

func (uc userUsecase) GetBasicInfo(id uint64) (domain.User, error) {
	user, err := uc.userRepo.GetUserById(id)
	user.OmitPassword()
	if err != nil {
		return domain.User{}, userErrors.ErrorInternalServer
	}
	return user, nil
}

func (uc userUsecase) Register(u *domain.User) (domain.User, error) {
	if u.FirstName == "" || u.Surname == "" || u.Email == "" ||
		u.Password == "" || u.Password != u.PasswordRepeat {
		return domain.User{}, userErrors.ErrorBadInput
	}
	_, err := uc.userRepo.GetUserByEmail(u.Email)
	if err == nil {
		return domain.User{}, userErrors.ErrorAlreadyExists
	}

	_, err = uc.userRepo.AddUser(u)
	if err != nil {
		return domain.User{}, userErrors.ErrorInternalServer
	}
	u.OmitPassword()
	return *u, nil
}

func (uc userUsecase) Login(u *domain.UserToLogin) (domain.User, error) {
	if u.Email == "" || u.Password == "" {
		return domain.User{}, userErrors.ErrorBadInput
	}
	us, err := uc.userRepo.GetUserByEmail(u.Email)
	errPassword := bcrypt.CompareHashAndPassword([]byte(us.Password), []byte(u.Password))
	if err != nil || errPassword != nil {
		return domain.User{}, userErrors.ErrorBadCredentials
	}
	us.OmitPassword()
	return us, nil
}

func (uc userUsecase) CheckAuth(id uint64) (domain.User, error) {
	us := domain.User{ID: id}
	return us, nil
}

func (uc userUsecase) GetProfileById(whoAskID, id uint64) (domain.Profile, error) {
	us, err := uc.userRepo.GetProfileById(whoAskID, id)
	if err != nil {
		return domain.Profile{}, err
	}
	return us, nil
}

func (uc userUsecase) UpdateProfile(profile domain.Profile) (domain.Profile, error) {
	us, err := uc.userRepo.UpdateProfile(profile)
	if err != nil {
		return domain.Profile{}, err
	}
	return us, nil
}

func (uc userUsecase) CreateSubscription(src, dst uint64) (domain.Profile, error) {
	us, err := uc.userRepo.CreateSubscription(src, dst)
	if err != nil {
		return domain.Profile{}, err
	}
	return us, nil
}

func (uc userUsecase) LoadUserReviews(id uint64, skip int, limit int) (domain.FilmReviews, error) {
	reviews, err := uc.userRepo.LoadUserReviews(id, skip, limit)
	if err != nil {
		return domain.FilmReviews{}, err
	}
	return reviews, nil
}
