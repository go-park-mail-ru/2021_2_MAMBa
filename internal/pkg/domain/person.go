package domain

type Person struct {
	Id           uint64   `json:"id,omitempty"`
	NameEn       string   `json:"name_en,omitempty"`
	NameRus      string   `json:"name_rus,omitempty"`
	PictureUrl   string   `json:"picture_url,omitempty"`
	Career       []string `json:"career,omitempty"`
	Height       float64  `json:"height,omitempty"`
	Age          float64  `json:"age,omitempty"`
	Birthday     string   `json:"birthday,omitempty"`
	Death        string   `json:"death,omitempty"`
	BirthPlace   string   `json:"birth_place,omitempty"`
	DeathPlace   string   `json:"death_place,omitempty"`
	Gender       string   `json:"gender,omitempty"`
	FamilyStatus string   `json:"family_status,omitempty"`
	FilmNumber   string   ` json:"film_number,omitempty"`
}
