package domain

type CollectionPreview struct {
	Id         uint64   `json:"id"`
	Title      string `json:"title"`
	PictureUrl string `json:"picture_url"`
}

type Collections struct {
	CollArray       []CollectionPreview `json:"collections_list"`
	MoreAvailable   bool                `json:"more_available"`
	CollectionTotal int                 `json:"collection_total"`
	CurrentSort     string              `json:"current_sort"`
	CurrentLimit    int                 `json:"current_limit"`
	CurrentSkip     int                 `json:"current_skip"`
}


type CollectionsRepository interface {
	GetCollections(skip int, limit int) (Collections, error)
}
//go:generate mockgen -destination=../collections/usecase/mock/usecase_mock.go  -package=mock 2021_2_MAMBa/internal/pkg/domain CollectionsUsecase
type CollectionsUsecase interface {
	GetCollections(skip int, limit int) (Collections, error)
}
