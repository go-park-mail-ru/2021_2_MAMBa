package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	mock2 "2021_2_MAMBa/internal/pkg/user/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testRow struct {
	inId       uint64
	inUser     domain.User
	inProfile  domain.Profile
	inId2      int64
	outUser    domain.User
	outProfile domain.Profile
	err        error
	name       string
}

var testTableGetSuccess = [...]testRow{
	{
		inId: 1,
		outUser: domain.User{
			ID:         1,
			FirstName:  "Test",
			Surname:    "Testovich",
			Email:      "abcd12",
			ProfilePic: "/1/pic.jpg",
		},
		err:  nil,
		name: `usecase get id`,
	},
}

var testTableGetFailure = [...]testRow{
	{
		inId:    10,
		outUser: domain.User{},
		err:     customErrors.ErrorInternalServer,
		name:    `usecase get id fl`,
	},
}

func TestGetBasicInfoSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range testTableGetSuccess {
		mock := mock2.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mock)
		mock.EXPECT().GetUserById(test.inId).Return(test.outUser, nil)
		actual, err := usecase.GetBasicInfo(test.inId)
		assert.Equal(t, test.outUser, actual, "Test: "+test.name)
		assert.Equal(t, test.err, err, "Test: "+test.name)
	}
}

func TestGetBasicInfoFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range testTableGetFailure {
		mock := mock2.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mock)
		mock.EXPECT().GetUserById(test.inId).Return(test.outUser, customErrors.ErrorNoUser)
		actual, err := usecase.GetBasicInfo(test.inId)
		assert.Equal(t, test.outUser, actual, "Test: "+test.name)
		assert.Equal(t, test.err, err, "Test: "+test.name)
	}
}

var testTableRegSuccess = [...]testRow{
	{
		inUser: domain.User{
			FirstName:      "Test",
			Surname:        "Testovich",
			Email:          "abcd12",
			Password:       "123",
			PasswordRepeat: "123",
		},
		outUser: domain.User{
			FirstName:  "Test",
			Surname:    "Testovich",
			Email:      "abcd12",
			ProfilePic: domain.BasePicture,
		},
		err:  nil,
		name: `usecase get id`,
	},
}

var testTableRegFailure = [...]testRow{
	{
		inUser: domain.User{
			FirstName:      "Test",
			Surname:        "Testovich",
			Email:          "abcd12",
			Password:       "124",
			PasswordRepeat: "123",
		},
		outUser: domain.User{},
		err:  customErrors.ErrorBadInput,
		name: `usecase get id`,
	},
	{
		inUser: domain.User{
			FirstName:      "Test",
			Surname:        "Testovich",
			Password:       "123",
			PasswordRepeat: "123",
		},
		outUser: domain.User{},
		err:  customErrors.ErrorBadInput,
		name: `usecase get id`,
	},
	{
		inUser: domain.User{
			FirstName:      "Test",
			Surname:        "Testovich",
			Email:          "abcd12",
			Password:       "123",
			PasswordRepeat: "123",
		},
		outUser: domain.User{},
		err:  customErrors.ErrorAlreadyExists,
		name: `usecase get id`,
	},
}

func TestRegisterSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for _, test := range testTableRegSuccess {
		mock := mock2.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mock)
		mock.EXPECT().GetUserByEmail(test.inUser.Email).Return(domain.User{}, customErrors.ErrorInternalServer)
		mock.EXPECT().AddUser(gomock.Any()).Return(test.outUser.ID, nil)
		actual, err := usecase.Register(&test.inUser)
		assert.Equal(t, test.outUser, actual, "Test: "+test.name)
		assert.Equal(t, test.err, err, "Test: "+test.name)
	}
}

func TestRegisterFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	for i, test := range testTableRegFailure {
		mock := mock2.NewMockUserRepository(ctrl)
		usecase := NewUserUsecase(mock)
		if i == 2 {
			mock.EXPECT().GetUserByEmail(test.inUser.Email).Return(domain.User{}, nil)
		}
		actual, err := usecase.Register(&test.inUser)
		assert.Equal(t, test.outUser, actual, "Test: "+test.name)
		assert.Equal(t, test.err, err, "Test: "+test.name)
	}
}
