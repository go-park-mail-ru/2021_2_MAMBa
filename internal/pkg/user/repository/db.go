package repository

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/user"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"sync"
)

type dbUserRepository struct {
	sync.RWMutex
	users []domain.User
}

func NewUserRepository() domain.UserRepository {
	return &dbUserRepository{}
}

func (ur *dbUserRepository) GetByEmail(email string) (domain.User, error) {
	ur.RLock()
	defer ur.RUnlock()
	loweredEmail := strings.ToLower(email)
	for _, us := range ur.users {
		if us.Email == loweredEmail {
			return us, nil
		}
	}
	return domain.User{}, user.ErrorNoUser
}

func (ur *dbUserRepository) GetById(id uint64) (domain.User, error) {
	ur.RLock()
	defer ur.RUnlock()
	if int(id) <= len(ur.users) && id != 0 {
		return ur.users[id-1], nil
	}
	return domain.User{}, user.ErrorNoUser
}

// AddUser TODO: Разобраться с локами, чтобы не было одинаковых ID у пользователей
func (ur *dbUserRepository) AddUser(us *domain.User) (uint64, error) {
	ur.RLock()
	us.ID = uint64(len(ur.users) + 1)
	ur.RUnlock()
	us.Email = strings.ToLower(us.Email)
	us.Surname = strings.Title(us.Surname)
	us.FirstName = strings.Title(us.FirstName)
	passwordByte, err := bcrypt.GenerateFromPassword([]byte(us.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, user.ErrorInternalServer
	}
	us.Password = string(passwordByte)
	us.ProfilePic = domain.BasePicture

	ur.Lock()
	ur.users = append(ur.users, *us)
	ur.Unlock()
	return us.ID, nil
}
