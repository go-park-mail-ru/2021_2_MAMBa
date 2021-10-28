package domain

const BasePicture = "/pic/1.jpg"

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
	GetById(id uint64) (User, error)
	GetByEmail(email string) (User, error)
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
