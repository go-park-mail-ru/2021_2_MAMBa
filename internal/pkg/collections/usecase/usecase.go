package usecase

import "2021_2_MAMBa/internal/pkg/domain"

type collectionsUsecase struct {
	collectionsRepo domain.CollectionsRepository
}

func NewCollectionsUsecase(u domain.CollectionsRepository) domain.CollectionsUsecase {
	return &collectionsUsecase{
		collectionsRepo: u,
	}
}

func (uc *collectionsUsecase) GetCollections(skip int, limit int) (domain.Collections, error) {
	colls, err := uc.collectionsRepo.GetCollections(skip, limit)
	if err != nil {
		return domain.Collections{}, err
	}
	return colls, err
}

func (uc *collectionsUsecase) GetCollectionPage (collId uint64) (domain.CollectionPage, error) {
	films, err := uc.collectionsRepo.GetCollectionFilms(collId)
	if err != nil {
		return domain.CollectionPage{}, err
	}
	coll, err := uc.collectionsRepo.GetCollectionInfo(collId)
	if err != nil {
		return domain.CollectionPage{}, err
	}

	collPage := domain.CollectionPage{
		Films: films,
		Coll:  coll,
	}
	return collPage, err
}
