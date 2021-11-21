package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"context"
)
import "2021_2_MAMBa/internal/pkg/collections/delivery/grpc"

type collectionsUsecase struct {
	collectionsRepo domain.CollectionsRepository
}

func NewCollectionsUsecase(u domain.CollectionsRepository) grpc.CollectionsRPCServer {
	return &collectionsUsecase{
		collectionsRepo: u,
	}
}

func (uc *collectionsUsecase) GetCollections(ctx context.Context, skl *grpc.SkipLimit) (*grpc.Collections, error) {
	skip := skl.Skip
	limit := skl.Limit
	colls, err := uc.collectionsRepo.GetCollections(int(skip), int(limit))
	if err != nil {
		return &grpc.Collections{}, err
	}
	buffer := make([]*grpc.CollectionPreview, 0)
	for _, elem := range colls.CollArray {
		buffer = append(buffer, &grpc.CollectionPreview{
			Id:         elem.Id,
			Title:      elem.Title,
			PictureUrl: elem.PictureUrl,
		})
	}

	return &grpc.Collections{
		CollArray:       buffer,
		MoreAvailable:   colls.MoreAvailable,
		CollectionTotal: int64(colls.CollectionTotal),
		CurrentSort:     colls.CurrentSort,
		CurrentLimit:    int64(colls.CurrentLimit),
		CurrentSkip:     int64(colls.CurrentSkip),
	}, err
}

func (uc *collectionsUsecase) GetCollectionPage(ctx context.Context, id *grpc.ID) (*grpc.CollectionPage, error) {
	collId := id.Id
	films, err := uc.collectionsRepo.GetCollectionFilms(collId)
	bufferFilms := make([]*grpc.Film, 0)
	for _, elem := range films {
		bufferFilms = append(bufferFilms, &grpc.Film{
			Id:          elem.Id,
			Title:       elem.Title,
			Description: elem.Description,
			ReleaseYear: int64(elem.ReleaseYear),
			PosterUrl:   elem.PosterUrl,
		})
	}
	if err != nil {
		return &grpc.CollectionPage{}, err
	}
	coll, err := uc.collectionsRepo.GetCollectionInfo(collId)
	if err != nil {
		return &grpc.CollectionPage{}, err
	}
	collPage := grpc.CollectionPage{
		Films: bufferFilms,
		Coll:  &grpc.Collection{
			Id:           coll.Id,
			AuthId:       coll.AuthId,
			CollName:     coll.CollName,
			Description:  coll.Description,
			CreationTime: coll.CreationTime,
			PicUrl:       coll.PicUrl,
		},
	}
	return &collPage, err
}
