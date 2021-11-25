package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
)

type dbSearchRepository struct {
	dbm *database.DBManager
}

func NewSearchRepository(manager *database.DBManager) domain.SearchRepository {
	return &dbSearchRepository{dbm: manager}
}

const (
	queryCountFoundFilms       = "SELECT COUNT(*) FROM film WHERE title ILIKE $1 OR title_original ILIKE $1"
	queryCountFoundPersons     = "SELECT COUNT(*) FROM person WHERE name_en ILIKE $1 OR name_rus ILIKE $1"
	querySearchFilmsByString   = "SELECT film_id FROM film WHERE title ILIKE $1 OR title_original ILIKE $1 LIMIT $2 OFFSET $3"
	querySearchPersonsByString = "SELECT person_id FROM person WHERE name_en ILIKE $1 OR name_rus ILIKE $1 LIMIT $2 OFFSET $3"
)

func (fr *dbSearchRepository) CountFoundFilms(query string) (int, error) {
	result, err := fr.dbm.Query(queryCountFoundFilms, "%"+query+"%")
	if err != nil {
		return 0, customErrors.ErrorInternalServer
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	return dbSize, nil
}

func (fr *dbSearchRepository) CountFoundPersons(query string) (int, error) {
	result, err := fr.dbm.Query(queryCountFoundPersons, "%"+query+"%")
	if err != nil {
		return 0, customErrors.ErrorInternalServer
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	return dbSize, nil
}

func (fr *dbSearchRepository) SearchFilmsIDList(query string, skip int, limit int) ([]uint64, error) {
	result, err := fr.dbm.Query(querySearchFilmsByString, "%"+query+"%", limit, skip)
	filmIdxList := make([]uint64, 0)
	if err != nil {
		return filmIdxList, err
	}
	for i := range result {
		filmIdxList = append(filmIdxList, cast.ToUint64(result[i][0]))
	}
	return filmIdxList, nil
}

func (fr *dbSearchRepository) SearchPersonsIDList(query string, skip int, limit int) ([]uint64, error) {
	result, err := fr.dbm.Query(querySearchPersonsByString, "%"+query+"%", limit, skip)
	personIdxList := make([]uint64, 0)
	if err != nil {
		return personIdxList, err
	}
	for i := range result {
		personIdxList = append(personIdxList, cast.ToUint64(result[i][0]))
	}
	return personIdxList, nil
}
