package domain

import "time"

const BasePicture = "/pic/1.jpg"

type Profile struct {
	ID        uint64 `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	Surname       string    `json:"surname,omitempty"`
	PictureUrl    string    `json:"picture_url,omitempty"`
	Email         string    `json:"email,omitempty"`
	Gender        string    `json:"gender,omitempty"`
	RegisterDate  time.Time `json:"register_date,omitempty"`
	SubCount      int       `json:"sub_count,omitempty"`
	BookmarkCount int       `json:"bookmark_count,omitempty"`
	AmSubscribed  bool      `json:"am_subscribed,omitempty"`
}

type User struct {
	ID             uint64 `json:"id,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	Surname        string `json:"surname,omitempty"`
	Email          string `json:"email,omitempty"`
	Password       string `json:"password,omitempty"`
	PasswordRepeat string `json:"password_repeat,omitempty"`
	ProfilePic     string `json:"profile_pic,omitempty"`
}

type UserToLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRepository interface {
	GetProfileById(whoAskID, id uint64) (Profile, error)
	UpdateProfile(profile Profile) (Profile, error)
	CheckSubscription(src, dst uint64) (bool, error)
	CreateSubscription(src, dst uint64) (Profile, error)
	LoadUserReviews(id uint64, skip int, limit int) (FilmReviews, error)
	GetUserById(id uint64) (User, error)
	GetUserByEmail(email string) (User, error)
	AddUser(user *User) (uint64, error)
}

type UserUsecase interface {
	GetBasicInfo(id uint64) (User, error)
	Register(u *User) (User, error)
	Login(u *UserToLogin) (User, error)
	CheckAuth(id uint64) (User, error)
}

func (us *User) OmitPassword() {
	us.Password = ""
	us.PasswordRepeat = ""
}

func (us *User) OmitId() {
	us.ID = 0
}

func (us *User) OmitPic() {
	us.ProfilePic = ""
}
