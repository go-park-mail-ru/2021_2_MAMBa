package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var fullFilled = User{
	ID:             1,
	FirstName:      "Ivan",
	Surname:        "Ivanov",
	Email:          "ivan1@mail.ru",
	Password:       "123456",
	PasswordRepeat: "123456",
	ProfilePic:     "/pic/1",
}

func TestOmitPassword(t *testing.T) {
	actual := fullFilled
	actual.OmitPassword()

	expected := fullFilled
	expected.Password = ""
	expected.PasswordRepeat = ""

	assert.ObjectsAreEqual(expected, actual)
}

func TestOmitId(t *testing.T) {
	actual := fullFilled
	actual.OmitId()

	expected := fullFilled
	expected.Password = ""

	assert.ObjectsAreEqual(expected, actual)
}

func TestOmitPic(t *testing.T) {
	actual := fullFilled
	actual.OmitPic()

	expected := fullFilled
	expected.ProfilePic = ""

	assert.ObjectsAreEqual(expected, actual)
}
