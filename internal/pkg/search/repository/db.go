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

func (fr *dbSearchRepository) GetSearch(query string, skipFilms int, limitFilms int, skipPersons int, limitPersons int) (domain.SearchResult, error) {
	result, err := fr.dbm.Query(queryCountFilmsBySearch, query)
	if err != nil {
		return domain.SearchResult{}, customErrors.ErrorInternalServer
	}
	dbSizeFilms := int(cast.ToUint64(result[0][0]))
	if skipFilms >= dbSizeFilms && skipFilms != 0 {
		return domain.SearchResult{}, customErrors.ErrorSkip
	}
	moreAvailableFilms := skipFilms+skipFilms < dbSizeFilms

	result, err := fr.dbm.Query(queryCountPersonsBySearch, query)
	if err != nil {
		return domain.SearchResult{}, customErrors.ErrorInternalServer
	}
	dbSizePersons := int(cast.ToUint64(result[0][0]))
	if skipPersons >= dbSizePersons && skipPersons != 0 {
		return domain.SearchResult{}, customErrors.ErrorSkip
	}
	moreAvailablePersons := skipPersons+skipPersons < dbSizePersons


	result, err = fr.dbm.Query(querySearchPersonsByQuery, query, limitFilms, skipFilms)
	if err != nil {
		return domain.SearchResult{}, customErrors.ErrorInternalServer
	}
	films := make([]domain.Film, 0)
	for i := range result {
		timestamp, err := cast.TimestampToString(result[i][6])
		if err != nil {
			return domain.FilmReviews{}, err
		}
		films = append(films, domain.Film{
			Id:              0,
			Title:           "",
			TitleOriginal:   "",
			Rating:          0,
			Description:     "",
			TotalRevenue:    "",
			PosterUrl:       "",
			TrailerUrl:      "",
			ContentType:     "",
			ReleaseYear:     0,
			Duration:        0,
			OriginCountries: nil,
			Cast:            nil,
			Director:        domain.Person{},
			Screenwriter:    domain.Person{},
			Genres:          nil,
		})
		reviews = append(reviews, domain.Review{
			Id:               cast.ToUint64(result[i][0]),
			FilmId:           cast.ToUint64(result[i][1]),
			AuthorName:       cast.ToString(result[i][7]) + " " + cast.ToString(result[i][8]),
			ReviewText:       cast.ToString(result[i][3]),
			AuthorPictureUrl: cast.ToString(result[i][9]),
			ReviewType:       int(cast.ToUint32(result[i][4])),
			Stars:            cast.ToFloat64(result[i][5]),
			Date:             timestamp,
		})
	}
	searchResult := domain.SearchResult{
	}
	return reviewsList, nil
}