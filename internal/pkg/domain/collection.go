package domain

type CollectionPreview struct {
	Id         uint64 `json:"id"`
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

type Collection struct {
	Id           uint64 `json:"id"`
	AuthId       uint64 `json:"auth_id"`
	CollName     string `json:"collection_name"`
	Description  string `json:"description"`
	CreationTime string `json:"creation_time"`
	PicUrl       string `json:"picture_url"`
}

type CollectionPage struct {
	Films []Film     `json:"films"`
	Coll  Collection `json:"collection"`
}

type CollectionsRepository interface {
	GetCollections(skip int, limit int) (Collections, error)
	GetCollectionFilms(id uint64) ([]Film, error)
	GetCollectionInfo(id uint64) (Collection, error)
}

//go:generate mockgen -destination=../collections/usecase/mock/usecase_mock.go  -package=mock 2021_2_MAMBa/internal/pkg/domain CollectionsUsecase
type CollectionsUsecase interface {
	GetCollections(skip int, limit int) (Collections, error)
	GetCollectionPage(collId uint64) (CollectionPage, error)
}
