package user

import (
	"net/http"
)

type userToLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userSignupForm struct {
	FirstName      string `json:"first_name"`
	Surname        string `json:"surname"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PasswordRepeat string `json:"password_repeat"`
}

func GetBasicInfo(w http.ResponseWriter, r *http.Request) {

}

func Register(w http.ResponseWriter, r *http.Request) {

}

func Login(w http.ResponseWriter, r *http.Request) {

}

func Logout(w http.ResponseWriter, r *http.Request) {

}

func CheckAuth(w http.ResponseWriter, r *http.Request) {

}
