package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
	"strings"
)

type dbPersonRepository struct {
	dbm *database.DBManager
}

func NewPersonRepository(manager *database.DBManager) domain.PersonRepository {
	return &dbPersonRepository{dbm: manager}
}

const (
	queryGetPerson             = "SELECT * FROM person WHERE person_id = $1"
	queryGetPersonFilms        = "SELECT f.film_id, f.title, f.description, f.release_year, f.poster_url FROM filmcast JOIN film f on filmcast.film_id = f.film_id where filmcast.person_id = $1 LIMIT $2 OFFSET $3"
	queryCountFilm             = "SELECT COUNT(*) FROM filmcast JOIN film f on filmcast.film_id = f.film_id where filmcast.person_id = $1"
	queryGetPersonFilmsPopular = "SELECT f.film_id, f.title, f.description, f.release_year, f.poster_url, ord.avg FROM (filmcast JOIN film f on filmcast.film_id = f.film_id) JOIN (SELECT AVG(stars) as avg, film_id FROM review WHERE (NOT type = 0) GROUP BY film_id) ord ON ord.film_id = f.film_id where filmcast.person_id = $1 ORDER BY ord.avg DESC LIMIT $2 OFFSET $3"
	queryCountPersonFilms      = "SELECT COUNT(*) FROM filmcast WHERE person_id = $1"
)

func (pr *dbPersonRepository) GetPerson(id uint64) (domain.Person, error) {
	result, err := pr.dbm.Query(queryGetPerson, id)
	if err != nil {
		return domain.Person{}, err
	}
	if len(result) == 0 {
		return domain.Person{}, customErrors.ErrNotFound
	}
	height := float64(int(cast.ToUint32(result[0][5])))
	person := domain.Person{
		Id:           cast.ToUint64(result[0][0]),
		NameEn:       cast.ToString(result[0][1]),
		NameRus:      cast.ToString(result[0][2]),
		PictureUrl:   cast.ToString(result[0][3]),
		Career:       strings.Split(cast.ToString(result[0][4]), ","),
		Height:       height,
		Age:          int(cast.ToUint32(result[0][6])),
		BirthPlace:   cast.ToString(result[0][9]),
		DeathPlace:   cast.ToString(result[0][10]),
		Gender:       cast.ToString(result[0][11]),
		FamilyStatus: cast.ToString(result[0][12]),
	}
	timestamp, err := cast.DateToString(result[0][7])
	if err != nil {
		return domain.Person{}, err
	}
	timestamp2, err := cast.DateToString(result[0][8])
	if err != nil {
		return domain.Person{}, err
	}
	person.Birthday = timestamp
	person.Death = timestamp2

	result, err = pr.dbm.Query(queryCountPersonFilms, id)
	if err != nil {
		return domain.Person{}, err
	}
	filmNumber := int(cast.ToUint64(result[0][0]))
	person.FilmNumber = filmNumber

	return person, nil
}
func (pr *dbPersonRepository) GetFilms(id uint64, skip int, limit int) (domain.FilmList, error) {
	result, err := pr.dbm.Query(queryCountFilm, id)
	if err != nil {
		return domain.FilmList{}, err
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	if skip >= dbSize && skip != 0 {
		return domain.FilmList{}, customErrors.ErrorSkip
	}

	result, err = pr.dbm.Query(queryGetPersonFilms, id, limit, skip)
	if err != nil {
		return domain.FilmList{}, err
	}
	filmList := make([]domain.Film, 0)
	for i := range result {
		film := domain.Film{
			Id:          cast.ToUint64(result[i][0]),
			Title:       cast.ToString(result[i][1]),
			Description: cast.ToString(result[i][2]),
			ReleaseYear: int(cast.ToUint32(result[i][3])),
			PosterUrl:   cast.ToString(result[i][4]),
		}
		filmList = append(filmList, film)
	}

	personFilms := domain.FilmList{
		FilmList:      filmList,
		MoreAvailable: skip+limit < dbSize,
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
	dbSize := int(cast.ToUint64(result[0][0]))
	if skip >= dbSize && skip != 0 {
		return domain.FilmList{}, customErrors.ErrorBadCredentials
	}

	result, err = pr.dbm.Query(queryGetPersonFilmsPopular, id, limit, skip)
	if err != nil {
		return domain.FilmList{}, err
	}
	filmList := make([]domain.Film, 0)
	for i := range result {
		film := domain.Film{
			Id:          cast.ToUint64(result[i][0]),
			Title:       cast.ToString(result[i][1]),
			Description: cast.ToString(result[i][2]),
			ReleaseYear: int(cast.ToUint32(result[i][3])),
			PosterUrl:   cast.ToString(result[i][4]),
		}
		filmList = append(filmList, film)
	}

	personFilms := domain.FilmList{
		FilmList:      filmList,
		MoreAvailable: skip+limit < dbSize,
		FilmTotal:     dbSize,
		CurrentLimit:  limit,
		CurrentSkip:   skip + limit,
	}
	return personFilms, nil
}
