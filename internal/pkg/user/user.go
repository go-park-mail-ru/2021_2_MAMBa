package user

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/sessions"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

type userToLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	db database.UserMockDatabase
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
	userInfo := user.Multipurpose()
	userInfo.OmitPassword()
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
	userForm := new(database.UserMultiPurpose)
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

	newUser := userForm.ToUser()

	idReg := db.AddUser(&newUser)
	err = sessions.StartSession(w, r, newUser.ID)
	if err != nil && idReg != 0 {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
	}
	userInfo := newUser.Multipurpose()
	userInfo.OmitPassword()
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
	errPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userForm.Password))
	if err != nil || errPassword != nil {
		http.Error(w, errorBadCredentials, http.StatusUnauthorized)
		return
	}
	_, err = sessions.CheckSession(r)
	if err != sessions.ErrUserNotLoggedIn {
		http.Error(w, errorAlreadyIn, http.StatusBadRequest)
		return
	}
	userInfo := user.Multipurpose()
	userInfo.OmitPassword()
	b, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
	}

	err = sessions.StartSession(w, r, user.ID)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
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

	userInfo := database.UserMultiPurpose{ID: userID}
	b, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, errorInternalServer, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
}
