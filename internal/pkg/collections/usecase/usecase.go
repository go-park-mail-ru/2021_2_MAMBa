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

func (uc collectionsUsecase) GetCollections(skip int, limit int) (domain.Collections, error) {
	colls, err := uc.collectionsRepo.GetCollections(skip, limit)
	if err != nil {
		return domain.Collections{}, err
	}
	return colls, err
}
