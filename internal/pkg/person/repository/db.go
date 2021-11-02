package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/person"
	"encoding/binary"
	"github.com/jackc/pgx/pgtype"
	"strings"
)

type dbPersonRepository struct {
	dbm *database.DBManager
}

func NewPersonRepository(manager *database.DBManager) domain.PersonRepository {
	return &dbPersonRepository{dbm: manager}
}

const (
	queryGetPerson      = "SELECT * FROM person WHERE person_id = $1"
	queryGetPersonFilms = "SELECT f.film_id, f.title, f.description, f.release_year, f.poster_url FROM filmcast JOIN film f on filmcast.film_id = f.film_id where filmcast.person_id = $1 LIMIT $2 OFFSET $3"
	queryCountFilm      = "SELECT COUNT(*) FROM filmcast JOIN film f on filmcast.film_id = f.film_id where filmcast.person_id = $1"
	queryGetPersonFilmsPopular = "SELECT f.film_id, f.title, f.description, f.release_year, f.poster_url, ord.avg FROM (filmcast JOIN film f on filmcast.film_id = f.film_id) JOIN (SELECT AVG(stars) as avg, film_id FROM review WHERE (NOT type = 0) GROUP BY film_id) ord ON ord.film_id = f.film_id where filmcast.person_id = $1 ORDER BY ord.avg DESC LIMIT $2 OFFSET $3"
)

func (pr *dbPersonRepository) GetPerson(id uint64) (domain.Person, error) {
	result, err := pr.dbm.Query(queryGetPerson, id)
	if err != nil {
		return domain.Person{}, err
	}
	intheight := int(binary.BigEndian.Uint32(result[0][5]))
	height := float64(intheight)
	person := domain.Person{
		Id:           binary.BigEndian.Uint64(result[0][0]),
		NameEn:       string(result[0][1]),
		NameRus:      string(result[0][2]),
		PictureUrl:   string(result[0][3]),
		Career:       strings.Split(string(result[0][4]), ","),
		Height:       height,
		Age:          int(binary.BigEndian.Uint32(result[0][6])),
		BirthPlace:   string(result[0][9]),
		DeathPlace:   string(result[0][10]),
		Gender:       string(result[0][11]),
		FamilyStatus: string(result[0][12]),
		FilmNumber:   string(result[0][13]),
	}
	timeBuffer1 := pgtype.Timestamp{}
	err = timeBuffer1.DecodeBinary(nil, result[0][7])
	if err != nil {
		return domain.Person{}, err
	}
	timeBuffer2 := pgtype.Timestamp{}
	err = timeBuffer2.DecodeBinary(nil, result[0][8])
	if err != nil {
		return domain.Person{}, err
	}
	person.Birthday = timeBuffer1.Time.String()
	person.Death = timeBuffer2.Time.String()
	return person, nil
}
func (pr *dbPersonRepository) GetFilms(id uint64, skip int, limit int) (domain.FilmList, error) {
	result, err := pr.dbm.Query(queryCountFilm, id)
	if err != nil {
		return domain.FilmList{}, err
	}
	dbSizeRaw := binary.BigEndian.Uint64(result[0][0])
	dbSize := int(dbSizeRaw)
	if skip >= dbSize {
		return domain.FilmList{}, person.ErrorSkip
	}

	result, err = pr.dbm.Query(queryGetPersonFilms, id, limit, skip)
	if err != nil {
		return domain.FilmList{}, err
	}
	filmList := make([]domain.Film, 0)
	for i := range result {
		film := domain.Film{
			Id:          binary.BigEndian.Uint64(result[i][0]),
			Title:       string(result[i][1]),
			Description: string(result[i][2]),
			ReleaseYear: int(binary.BigEndian.Uint32(result[i][3])),
			PosterUrl:   string(result[i][4]),
		}
		filmList = append(filmList, film)
	}

	personFilms := domain.FilmList{
		FilmList:      filmList,
		MoreAvaliable: skip + limit < dbSize,
		FilmTotal:     dbSize,
		CurrentLimit:  limit,
		CurrentSkip:   skip + limit,
	}
	return personFilms, nil
}
func (pr *dbPersonRepository) GetFilmsPopular(id uint64, skip int, limit int) (domain.FilmList, error) {
	result, err := pr.dbm.Query(queryCountFilm, id)
	if err != nil {
		return domain.FilmList{}, err
	}
	dbSizeRaw := binary.BigEndian.Uint64(result[0][0])
	dbSize := int(dbSizeRaw)
	if skip >= dbSize {
		return domain.FilmList{}, person.ErrorBadCredentials
	}

	result, err = pr.dbm.Query(queryGetPersonFilmsPopular, id, limit, skip)
	if err != nil {
		return domain.FilmList{}, err
	}
	filmList := make([]domain.Film, 0)
	for i := range result {
		film := domain.Film{
			Id:          binary.BigEndian.Uint64(result[i][0]),
			Title:       string(result[i][1]),
			Description: string(result[i][2]),
			ReleaseYear: int(binary.BigEndian.Uint32(result[i][3])),
			PosterUrl:   string(result[i][4]),
		}
		filmList = append(filmList, film)
	}

	personFilms := domain.FilmList{
		FilmList:      filmList,
		MoreAvaliable: skip + limit < dbSize,
		FilmTotal:     dbSize,
		CurrentLimit:  limit,
		CurrentSkip:   skip + limit,
	}
	return personFilms, nil
}
