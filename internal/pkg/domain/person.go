package domain

type Person struct {
	Id           uint64   `json:"id,omitempty"`
	NameEn       string   `json:"name_en,omitempty"`
	NameRus      string   `json:"name_rus,omitempty"`
	PictureUrl   string   `json:"picture_url,omitempty"`
	Career       []string `json:"career,omitempty"`
	Height       float64  `json:"height,omitempty"`
	Age          int      `json:"age,omitempty"`
	Birthday     string   `json:"birthday,omitempty"`
	Death        string   `json:"death,omitempty"`
	BirthPlace   string   `json:"birth_place,omitempty"`
	DeathPlace   string   `json:"death_place,omitempty"`
	Gender       string   `json:"gender,omitempty"`
	FamilyStatus string   `json:"family_status,omitempty"`
	FilmNumber   string   ` json:"film_number,omitempty"`
}

type FilmList struct {
	FilmList      []Film `json:"film_list"`
	MoreAvaliable bool   `json:"more_avaliable"`
	FilmTotal     int    `json:"film_total"`
	CurrentLimit  int    `json:"current_limit"`
	CurrentSkip   int    `json:"current_skip"`
}

type PersonPage struct {
	Actor        Person   `json:"actor"`
	Films        FilmList `json:"films"`
	PopularFilms FilmList `json:"popular_films"`
}

type PersonRepository interface {
	GetPerson(id uint64) (Person, error)
	GetFilms(id uint64, skip int, limit int) (FilmList, error)
	GetFilmsPopular(id uint64, skip int, limit int) (FilmList, error)
}

type PersonUsecase interface {
	GetPerson(id uint64) (PersonPage, error)
	GetFilms(id uint64, skip int, limit int) (FilmList, error)
}
