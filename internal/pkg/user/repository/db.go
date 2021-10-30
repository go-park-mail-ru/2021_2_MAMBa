package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/user"
	"encoding/binary"
	"golang.org/x/crypto/bcrypt"
)

type dbUserRepository struct {
	dbm *database.DBManager
}

func NewUserRepository(manager *database.DBManager) domain.UserRepository {
	return &dbUserRepository{dbm: manager}
}

const (
	queryGetById    = "SELECT * FROM Profile WHERE User_ID = $1"
	queryGetByEmail = "SELECT * FROM Profile WHERE email = $1"
	queryAddUser    = "INSERT INTO Profile(first_name, surname, email, password, picture_url, register_date) VALUES ($1, $2, $3, $4, $5, current_timestamp) RETURNING User_ID"
)

func (ur *dbUserRepository) GetByEmail(email string) (domain.User, error) {
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

func (ur *dbUserRepository) GetById(id uint64) (domain.User, error) {
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
