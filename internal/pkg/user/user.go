package user

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/sessions"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

type userBasicInfo struct {
	ID         uint64 `json:"id"`
	FirstName  string `json:"first_name"`
	Surname    string `json:"surname"`
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
}

var (
	db database.Database

	errorBadInput       = "error - bad input"
	errorAlreadyIn      = "error - already in"
	errorBadCredentials = "error - bad credentials"
	errorInternalServer = "Internal server error"
)

func GetBasicInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	u64, err := strconv.ParseUint(params["id"], 10, 64)
	if err != nil {
		http.Error(w, errorBadInput, http.StatusBadRequest)
		return
	}
	user, err := db.FindId(u64)
	if err != nil {
		http.Error(w, errorBadInput, http.StatusNotFound)
		return
	}
	userInfo := &userBasicInfo{
		ID:         u64,
		FirstName:  user.FirstName,
		Surname:    user.Surname,
		Email:      user.Email,
		ProfilePic: user.ProfilePic,
	}
	b, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
		return
	}
}

func Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userForm := new(userSignupForm)
	err := json.NewDecoder(r.Body).Decode(&userForm)
	if err != nil {
		http.Error(w, errorBadInput, http.StatusBadRequest)
		return
	}

	if userForm.FirstName == "" || userForm.Surname == "" || userForm.Email == "" ||
		userForm.Password == "" || userForm.Password != userForm.PasswordRepeat {
		http.Error(w, errorBadInput, http.StatusBadRequest)
		return
	}
	_, err = db.FindEmail(userForm.Email)
	if err == nil {
		http.Error(w, errorAlreadyIn, http.StatusConflict)
		return
	}

	newUser := &database.User{
		// ID устанавливается в addUser под заблокированным RWMutex
		Email:      userForm.Email,
		FirstName:  userForm.FirstName,
		Surname:    userForm.Surname,
		Password:   userForm.Password,
		ProfilePic: "/pic/1.jpg"}

	idReg := db.AddUser(newUser)
	err = sessions.StartSession(w, r, newUser.ID)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
	}
	userInfo := &userBasicInfo{
		ID:         idReg,
		FirstName:  userForm.FirstName,
		Surname:    userForm.Surname,
		Email:      userForm.Email,
		ProfilePic: "/pic/1.jpg",
	}
	b, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(b)
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	userForm := new(userToLogin)
	err := json.NewDecoder(r.Body).Decode(&userForm)
	if err != nil {
		http.Error(w, errorBadInput, http.StatusBadRequest)
		return
	}
	if userForm.Email == "" || userForm.Password == "" {
		http.Error(w, errorBadInput, http.StatusBadRequest)
		return
	}
	user, err := db.FindEmail(userForm.Email)
	if err != nil || user.Password != userForm.Password {
		http.Error(w, errorBadCredentials, http.StatusUnauthorized)
		return
	}
	_, err = sessions.CheckSession(r)
	if err != sessions.ErrUserNotLoggedIn {
		http.Error(w, errorAlreadyIn, http.StatusBadRequest)
		return
	}
	userInfo := &userBasicInfo{
		ID:         user.ID,
		FirstName:  user.FirstName,
		Surname:    user.Surname,
		Email:      user.Email,
		ProfilePic: user.ProfilePic,
	}
	b, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	err = sessions.StartSession(w, r, user.ID)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	id, err := sessions.CheckSession(r)
	if err == sessions.ErrUserNotLoggedIn {
		http.Error(w, errorBadInput, http.StatusForbidden)
		return
	}
	err = sessions.EndSession(w, r, id)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	userID, err := sessions.CheckSession(r)
	if err == sessions.ErrUserNotLoggedIn {
		http.Error(w, errorBadInput, http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(strconv.FormatUint(userID, 10)))
}
