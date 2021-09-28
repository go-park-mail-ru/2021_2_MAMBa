package database

import (
	"errors"
	"strings"
	"sync"
)

type Database struct {
	sync.RWMutex
	users []User
}

type User struct {
	ID         uint64
	FirstName  string `json:"first_name"`
	Surname    string `json:"surname"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ProfilePic string
}

var errorNoUser = errors.New("error: no user")

func (db *Database) AddUser(us *User) {
	us.ID = uint64(len(db.users) + 1)
	us.Email = strings.ToLower(us.Email)
	us.Surname = strings.Title(us.Surname)
	us.FirstName = strings.Title(us.FirstName)

	db.Lock()
	db.users = append(db.users, *us)
	db.Unlock()
}

func (db *Database) FindEmail(email string) (us User, err error) {
	db.RLock()
	defer db.RUnlock()
	loweredEmail := strings.ToLower(email)
	for _, us := range db.users {
		if us.Email == loweredEmail {
			return us, nil
		}
	}
	return User{}, errorNoUser
}

func (db *Database) FindId(id uint64) (us User, err error) {
	db.RLock()
	defer db.RUnlock()
	if int(id) < len(db.users) {
		return db.users[id], nil
	}
	return User{}, errorNoUser
}
