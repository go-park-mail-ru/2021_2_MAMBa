package repository

import (
	"2021_2_MAMBa/internal/pkg/collections"
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	"2021_2_MAMBa/internal/pkg/user"
	"encoding/binary"
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
)

func (cr *dbCollectionsRepository) GetCollections(skip int, limit int) (domain.Collections, error) {

	result, err := cr.dbm.Query(queryCountCollections)
	if err != nil {
		return domain.Collections{}, user.ErrorInternalServer
	}
	dbSizeRaw := binary.BigEndian.Uint64(result[0][0])
	dbSize := int(dbSizeRaw)
	if skip >= dbSize {
		return domain.Collections{}, collections.ErrorSkip
	}

	if err != nil {
		return domain.Collections{}, user.ErrorInternalServer
	}
	if len(result) == 0 {
		return domain.Collections{}, user.ErrorNoUser
	}

	moreAvailable := skip+limit < dbSize
	result, err = cr.dbm.Query(queryGetCollections, limit, skip)

	previews := make([]domain.CollectionPreview, 0)
	for i := range result {
		previewBuffer := domain.CollectionPreview{
			Id:         uint(binary.BigEndian.Uint64(result[i][0])),
			Title:      string(result[i][1]),
			PictureUrl: string(result[i][2]),
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
