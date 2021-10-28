package repository

import (
	"2021_2_MAMBa/internal/pkg/collections"
	"2021_2_MAMBa/internal/pkg/domain"
	"sync"
)

type dbCollectionsRepository struct {
	sync.RWMutex
	Previews []domain.CollectionPreview
}

func NewCollectionsRepository() domain.CollectionsRepository {
	return &dbCollectionsRepository{Previews: PreviewMock}
}

var PreviewMock = []domain.CollectionPreview{
	{Id: 1, Title: "Для ценителей Хогвартса", PictureUrl: "server/images/collections1.png"},
	{Id: 2, Title: "Про настоящую любовь", PictureUrl: "server/images/collections2.png"},
	{Id: 3, Title: "Аферы века", PictureUrl: "server/images/collections3.png"},
	{Id: 4, Title: "Про Вторую Мировую", PictureUrl: "server/images/collections4.jpg"},
	{Id: 5, Title: "Осеннее настроение", PictureUrl: "server/images/collections5.png"},
	{Id: 6, Title: "Летняя атмосфера", PictureUrl: "server/images/collections6.png"},
	{Id: 7, Title: "Про дружбу", PictureUrl: "server/images/collections7.png"},
	{Id: 8, Title: "Романтические фильмы", PictureUrl: "server/images/collections8.jpg"},
	{Id: 9, Title: "Джунгли зовут", PictureUrl: "server/images/collections9.jpg"},
	{Id: 10, Title: "Фантастические фильмы", PictureUrl: "server/images/collections10.jpg"},
	{Id: 11, Title: "Про петлю времени", PictureUrl: "server/images/collections11.png"},
	{Id: 12, Title: "Классика на века", PictureUrl: "server/images/collections12.jpg"},
}

func (cr *dbCollectionsRepository) GetCollections(skip int, limit int) (domain.Collections, error) {
	cr.RLock()
	dbSize := len(cr.Previews)
	cr.RUnlock()
	if skip >= dbSize {
		return domain.Collections{}, collections.ErrorSkip
	}
	moreAvailable := skip+limit < dbSize
	next := skip + limit
	if !moreAvailable {
		next = dbSize
	}
	cr.RLock()
	previews := cr.Previews[skip:next]
	cr.RUnlock()

	collect := domain.Collections{
		CollArray:       previews,
		MoreAvailable:   moreAvailable,
		CollectionTotal: dbSize,
		CurrentLimit:    limit,
		CurrentSkip:     skip + limit,
	}
	return collect, nil
}
