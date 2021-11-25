package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
	"2021_2_MAMBa/internal/pkg/utils/cast"
)

type dbCollectionsRepository struct {
	dbm *database.DBManager
}

func NewCollectionsRepository(manager *database.DBManager) domain.CollectionsRepository {
	return &dbCollectionsRepository{dbm: manager}
}

const (
	queryCountCollections = "SELECT COUNT(*) FROM Collection"
	queryGetCollections   = "SELECT collection_id, collection_name, picture_url FROM Collection LIMIT $1 OFFSET $2 "
	queryCountFilms       = "SELECT COUNT(*) from film f join collectionconnection c on f.film_id = c.film_id WHERE collection_id = $1"
	queryGetFilms         = "SELECT f.film_id, f.title, f.description, f.release_year, f.poster_url from film f join collectionconnection c on f.film_id = c.film_id WHERE collection_id = $1 ORDER BY f.film_id"
	queryGetCollection    = "SELECT * FROM collection WHERE collection_id = $1"
)

func (cr *dbCollectionsRepository) GetCollections(skip int, limit int) (domain.Collections, error) {
	result, err := cr.dbm.Query(queryCountCollections)
	if err != nil {
		return domain.Collections{}, customErrors.ErrorInternalServer
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	if skip >= dbSize && skip != 0 {
		return domain.Collections{}, customErrors.ErrorSkip
	}

	moreAvailable := skip+limit < dbSize
	result, err = cr.dbm.Query(queryGetCollections, limit, skip)

	previews := make([]domain.CollectionPreview, 0)
	for i := range result {
		previewBuffer := domain.CollectionPreview{
			Id:         cast.ToUint64(result[i][0]),
			Title:      cast.ToString(result[i][1]),
			PictureUrl: cast.ToString(result[i][2]),
		}
		previews = append(previews, previewBuffer)
	}

	collect := domain.Collections{
		CollArray:       previews,
		MoreAvailable:   moreAvailable,
		CollectionTotal: dbSize,
		CurrentLimit:    limit,
		CurrentSkip:     skip + limit,
	}
	return collect, nil
}

func (cr *dbCollectionsRepository) GetCollectionFilms(id uint64) ([]domain.Film, error) {
	result, err := cr.dbm.Query(queryCountFilms, id)
	if err != nil {
		return []domain.Film{}, err
	}
	dbSize := int(cast.ToUint64(result[0][0]))
	if dbSize == 0 {
		return []domain.Film{}, customErrors.ErrorBadInput
	}
	result, err = cr.dbm.Query(queryGetFilms, id)
	if err != nil {
		return []domain.Film{}, err
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
	return filmList, nil
}

func (cr *dbCollectionsRepository) GetCollectionInfo(id uint64) (domain.Collection, error) {
	result, err := cr.dbm.Query(queryGetCollection, id)
	if err != nil {
		return domain.Collection{}, err
	}
	timeString, err := cast.TimestampToString(result[0][4])
	if err != nil {
		return domain.Collection{}, err
	}
	collResult := domain.Collection{
		Id:           cast.ToUint64(result[0][0]),
		AuthId:       cast.ToUint64(result[0][1]),
		CollName:     cast.ToString(result[0][2]),
		Description:  cast.ToString(result[0][3]),
		CreationTime: timeString,
		PicUrl:       cast.ToString(result[0][5]),
	}
	return collResult, nil
}
