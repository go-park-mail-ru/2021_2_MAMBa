package user

import (
	"2021_2_MAMBa/internal/pkg/sessions"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
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
	FirstName      string `json:"first_name"`
	Surname        string `json:"surname"`
	Email          string `json:"email"`
	ProfilePic	   string `json:"profile_description"`
}

type User struct {
	ID             uint
	FirstName      string `json:"first_name"`
	Surname        string `json:"surname"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ProfilePic     string
}

type mockDB struct {
	sync.RWMutex
	users []User
}

func (db *mockDB) addUser (user User){
	db.Lock()
	db.users = append(db.users, user)
	db.Unlock()
}

var (
	db mockDB
	errorNoUser = errors.New(`"error": "no user"`)
)



func (db *mockDB) findEmail (email string) (user User, err error) {
	db.RLock()
	defer db.RUnlock()
	for _, us := range db.users {
		if us.Email == email {
			return us, nil
		}
	}
	return User{}, errorNoUser
}

func (db *mockDB) findId (id uint) (user User, err error) {
	db.RLock()
	defer db.RUnlock()
	if int(id) < len(db.users) {
		return db.users[id], nil
	}
	return User{}, errorNoUser
}


func GetBasicInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	u64, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		http.Error(w, `error - bad input`, http.StatusBadRequest)
		return
	}
	user, err := db.findId(uint(u64))
	if err != nil {
		http.Error(w, `error - bad input`, http.StatusBadRequest)
		return
	}
	userInfo := userBasicInfo{
		FirstName: user.FirstName,
		Surname: user.Surname,
		Email: user.Email,
		ProfilePic: user.ProfilePic,
		}
	err = json.NewEncoder(w).Encode(userInfo)
	if err != nil {
		http.Error(w, `error - bad input`, http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	userForm := new(userSignupForm)
	err := decoder.Decode(userForm)
	if err != nil {
		http.Error(w, `error - bad input`, http.StatusBadRequest)
		return
	}
	if userForm.Email == "" || userForm.Password == "" || userForm.Password != userForm.PasswordRepeat{
		http.Error(w, `error - bad input`, http.StatusBadRequest)
		return
	}
	_, err = db.findEmail(userForm.Email)
	if err == nil {
		http.Error(w, `error - already in`, http.StatusConflict)
		return
	}
	db.addUser(User{ID: uint(len(db.users) + 3),
					Email: userForm.Email,
					FirstName: userForm.FirstName,
					Password: userForm.Password,
					ProfilePic: "/pic/1.jpg"})
	w.WriteHeader(http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	userForm := new(userToLogin)
	err := decoder.Decode(userForm)
	if err != nil {
		http.Error(w, `error - bad input`, http.StatusBadRequest)
		return
	}
	if userForm.Email == "" || userForm.Password == "" {
		http.Error(w, `error - bad input`, http.StatusBadRequest)
		return
	}
	user, err := db.findEmail(userForm.Email)
	if err == nil || user.Password != userForm.Password {
		http.Error(w, `error - bad credentials`, http.StatusUnauthorized)
		return
	}
	err = sessions.StartSession(w, r, user.ID)
	if err != nil {
		http.Error(w, `Internal server error`, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	id, err := sessions.CheckSession(r)
	if err == sessions.ErrUserNotLoggedIn {
		http.Error(w, `error - bad input`, http.StatusForbidden)
		return
	}
	err = sessions.EndSession(w, r, id)
	if err != nil {
		http.Error(w, `Internal server error`, http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	_, err := sessions.CheckSession(r)
	if err == sessions.ErrUserNotLoggedIn {
		http.Error(w, `error - bad input`, http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}
