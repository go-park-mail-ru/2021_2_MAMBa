package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/log"
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
	user, err := uc.userRepo.GetById(id)
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	log.Info(string(passwordByte))
	user.OmitPassword()
	if err != nil {
		return domain.User{}, customErrors.ErrorInternalServer
	}
	return user, nil
}

func (uc userUsecase) Register(u *domain.User) (domain.User, error) {
	if u.FirstName == "" || u.Surname == "" || u.Email == "" ||
		u.Password == "" || u.Password != u.PasswordRepeat {
		return domain.User{}, customErrors.ErrorBadInput
	}
	_, err := uc.userRepo.GetByEmail(u.Email)
	if err == nil {
		return domain.User{}, customErrors.ErrorAlreadyExists
	}

	// соль пароль
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, customErrors.ErrorInternalServer
	}
	u.Password = string(passwordByte)
	u.ProfilePic = domain.BasePicture

	_, err = uc.userRepo.AddUser(u)
	if err != nil {
		return domain.User{}, customErrors.ErrorInternalServer
	}
	u.OmitPassword()
	return *u, nil
}

func (uc userUsecase) Login(u *domain.UserToLogin) (domain.User, error) {
	if u.Email == "" || u.Password == "" {
		return domain.User{}, customErrors.ErrorBadInput
	}
	us, err := uc.userRepo.GetByEmail(u.Email)
	errPassword := bcrypt.CompareHashAndPassword([]byte(us.Password), []byte(u.Password))
	if err != nil || errPassword != nil {
		return domain.User{}, customErrors.ErrorBadCredentials
	}
	us.OmitPassword()
	return us, nil
}

func (uc userUsecase) CheckAuth(id uint64) (domain.User, error) {
	us := domain.User{ID: id}
	return us, nil
}
