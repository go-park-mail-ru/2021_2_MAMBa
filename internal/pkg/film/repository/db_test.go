package repository

import (
	"2021_2_MAMBa/internal/pkg/database"
	"2021_2_MAMBa/internal/pkg/domain"
	mylog "2021_2_MAMBa/internal/pkg/utils/log"
	"encoding/binary"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"math"
	"regexp"
	"testing"
	"time"
)

type testRow struct {
	inQuery    string
	bodyString string
	out        string
	status     int
	name       string
}

func MockDatabase() (*database.DBManager, pgxmock.PgxPoolIface, error) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		mylog.Error(errors.New("failed to create mock"))
	}
	return &database.DBManager{Pool: mock}, mock, err
}

var mockPersonPreview = domain.Person{
	Id:         1,
	NameEn:     "Miley",
	NameRus:    "Сайрус",
	PictureUrl: "/miley.webp",
	Career:     []string{"Актриса"},
}

func TestGetSuccess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()

	mf := domain.Film{
		Id:              1,
		Title:           "Тест",
		TitleOriginal:   "Test",
		Rating:          4.1,
		Description:     "this is a test film",
		TotalRevenue:    "$100",
		PosterUrl:       "/img/poster1.webp",
		TrailerUrl:      "/trailer/mov.mov",
		ContentType:     "film",
		ReleaseYear:     2021,
		Duration:        80,
		OriginCountries: []string{"США", "Канада"},
		Cast:            []domain.Person{mockPersonPreview},
		Director:        mockPersonPreview,
		Screenwriter:    mockPersonPreview,
		Genres:          []domain.Genre{{Id: 1, Name: "Комедия"}},
		PremiereRu:      "2005-07-14",
	}
	rowsFilm := pgxmock.NewRows([]string{"id", "title", "title_original", "rating", "description",
		"poster_url", "trailer", "total_revenue", "release", "duration", "scr", "dir",
		"cont_type", "premiere", "pid", "pnameen", "pnameru", "ppicture", "pcareer",
		"p1id", "p1nameen", "p1nameru", "p1picture", "p1career"})
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, mf.Id)
	byteRating := make([]byte, 8)
	binary.BigEndian.PutUint64(byteRating, math.Float64bits(mf.Rating))
	byteRelease := make([]byte, 4)
	binary.BigEndian.PutUint32(byteRelease, uint32(mf.ReleaseYear))
	byteDuration := make([]byte, 4)
	binary.BigEndian.PutUint32(byteDuration, uint32(mf.Duration))
	byteGenreId := make([]byte, 4)
	binary.BigEndian.PutUint32(byteGenreId, 1)
	rowsFilm.AddRow(byteId, []uint8(mf.Title), []uint8(mf.TitleOriginal), byteRating, []uint8(mf.Description),
		[]uint8(mf.PosterUrl), []uint8(mf.TrailerUrl), []uint8(mf.TotalRevenue), byteRelease, byteDuration, byteId, byteId,
		[]uint8(mf.ContentType), byteRelease, byteId, []uint8(mf.Director.NameEn), []uint8(mf.Director.NameRus), []uint8(mf.Director.PictureUrl), []uint8(mf.Director.Career[0]),
		byteId, []uint8(mf.Director.NameEn), []uint8(mf.Director.NameRus), []uint8(mf.Director.PictureUrl), []uint8(mf.Director.Career[0]))
	rowsRate := pgxmock.NewRows([]string{"rating"}).AddRow(byteRating)
	rowsCountry := pgxmock.NewRows([]string{"country_name"}).AddRow([]uint8("США"))
	rowsCountry.AddRow([]uint8("Канада"))
	rowsGenres := pgxmock.NewRows([]string{"id", "genre_name"}).AddRow(byteGenreId, []uint8(mf.Genres[0].Name))
	rowsCast := pgxmock.NewRows([]string{"pid", "pnameen", "pnameru", "ppicture", "pcareer"})
	rowsCast.AddRow(byteId, []uint8(mf.Director.NameEn), []uint8(mf.Director.NameRus), []uint8(mf.Director.PictureUrl), []uint8(mf.Director.Career[0]))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmDirScr)).WithArgs(mf.Id).WillReturnRows(rowsFilm)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmRating)).WithArgs(mf.Id).WillReturnRows(rowsRate)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmCountries)).WithArgs(mf.Id).WillReturnRows(rowsCountry)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmGenres)).WithArgs(mf.Id).WillReturnRows(rowsGenres)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmCast)).WithArgs(mf.Id).WillReturnRows(rowsCast)
	pool.ExpectCommit()

	actual, err := repository.GetFilm(mf.Id)
	assert.NoError(t, err)
	assert.Equal(t, mf, actual)
}

func TestGetFilmRecommendationsSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()

	recom := domain.FilmRecommendations{
		RecommendationList: []domain.Film{{
			Id:        1,
			Title:     "Test_film",
			PosterUrl: "/pic/TestPoster.webp"}},
		MoreAvailable:       false,
		RecommendationTotal: 1,
		CurrentLimit:        10,
		CurrentSkip:         10,
	}
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, recom.RecommendationList[0].Id)
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))

	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	rowsRecom := pgxmock.NewRows([]string{"id", "title", "poster_url"})
	rowsRecom.AddRow(byteId, []uint8(recom.RecommendationList[0].Title), []uint8(recom.RecommendationList[0].PosterUrl))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFilmRecommendations)).WithArgs(uint64(1)).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmRecommendations)).WithArgs(uint64(1), 10, 0).WillReturnRows(rowsRecom)
	pool.ExpectCommit()

	actual, err := repository.GetFilmRecommendations(1, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, recom, actual)
}

func TestGetFilmReviewsSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()
	timeBuffer := pgtype.Timestamp{}
	timeBuffer.Time = time.Now()
	reviews := domain.FilmReviews{
		ReviewList: []domain.Review{{
			Id:               1,
			FilmId:           1,
			AuthorName:       "Ivan Ivanovich",
			ReviewText:       "Test review on film",
			AuthorId:         1,
			AuthorPictureUrl: "pic1.jopeg",
			ReviewType:       3,
			Stars:            4.0,
			Date:             "",
		}},
		MoreAvailable: false,
		ReviewTotal:   1,
		CurrentLimit:  10,
		CurrentSkip:   10,
	}
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, uint64(1))
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))
	byteTime := make([]byte, 16)
	byteTime, err = timeBuffer.EncodeBinary(nil, byteTime)
	byteType := make([]byte, 4)
	binary.BigEndian.PutUint32(byteType, uint32(reviews.ReviewList[0].ReviewType))
	byteStars := make([]byte, 8)
	binary.BigEndian.PutUint64(byteStars, math.Float64bits(reviews.ReviewList[0].Stars))

	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	rowsRev := pgxmock.NewRows([]string{"review_id", "film_id", "authid", "text", "type", "stars", "date", "fname", "sname", "url"})
	rowsRev.AddRow(byteId, byteId, byteId, []uint8(reviews.ReviewList[0].ReviewText), byteType, byteStars, byteTime, []uint8("Ivan"), []uint8("Ivanovich"), []uint8(reviews.ReviewList[0].AuthorPictureUrl))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFilmReviews)).WithArgs(uint64(1)).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmReviews)).WithArgs(uint64(1), 10, 0).WillReturnRows(rowsRev)
	pool.ExpectCommit()

	actual, err := repository.GetFilmReviews(1, 0, 10)
	assert.NoError(t, err)
	assert.Equal(t, reviews, actual)
}

func TestPostFilmRatingSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()

	countByte := make([]uint8, 8)
	binary.BigEndian.PutUint64(countByte, uint64(1))
	byteRating := make([]uint8, 8)
	binary.BigEndian.PutUint64(byteRating, math.Float64bits(10.0))

	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(countByte)
	rowsRate := pgxmock.NewRows([]string{"rating"}).AddRow(byteRating)
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetReviewByAuthor)).WithArgs(uint64(1), uint64(1)).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectExec(regexp.QuoteMeta(queryUpdateRating)).WithArgs(float64(10), uint64(1), uint64(1)).WillReturnResult(pgconn.CommandTag{})
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmRating)).WithArgs(uint64(1)).WillReturnRows(rowsRate)
	pool.ExpectCommit()

	actual, err := repository.PostRating(uint64(1), uint64(1), 10.0)
	assert.NoError(t, err)
	assert.Equal(t, 10.0, actual)
}

func TestGetFilmPopSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()

	f := domain.FilmList{
		FilmList: []domain.Film{{
			Id:            1,
			Title:         "Test_film",
			TitleOriginal: "bnm",
			ReleaseYear:   2021,
			Description:   "asdf",
			Rating:        2.3,
			PosterUrl:     "/pic/TestPoster.webp",
			PremiereRu:    "2000-01-01"}},
		MoreAvailable: false,
	}
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, f.FilmList[0].Id)
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))
	byteRating := make([]byte, 8)
	binary.BigEndian.PutUint64(byteRating, math.Float64bits(f.FilmList[0].Rating))
	byteRelease := make([]byte, 4)
	binary.BigEndian.PutUint32(byteRelease, uint32(f.FilmList[0].ReleaseYear))
	bytePrem := make([]byte, 4)
	binary.BigEndian.PutUint32(bytePrem, 0)

	rowsRecom := pgxmock.NewRows([]string{"id", "title", "title_or", "release", "desc", "poster_url", "premiere", "rate"})
	rowsRecom.AddRow(byteId, []uint8(f.FilmList[0].Title), []uint8(f.FilmList[0].TitleOriginal), byteRelease, []uint8(f.FilmList[0].Description), []uint8(f.FilmList[0].PosterUrl), bytePrem, byteRating)

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetPopularFilms)).WillReturnRows(rowsRecom)
	pool.ExpectCommit()

	actual, err := repository.GetPopularFilms()
	assert.NoError(t, err)
	assert.Equal(t, f, actual)
}

func TestGetFilmBannersSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()

	f := domain.BannersList{
		BannersList: []domain.Banner{{
			Id:          1,
			Title:       "Test_film",
			Description: "asdf",
			PictureURL:  "/pic/TestPoster.webp",
			Link:        "2000-01-01"}},
	}
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, f.BannersList[0].Id)
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))
	bytePrem := make([]byte, 4)
	binary.BigEndian.PutUint32(bytePrem, 0)

	rowsRecom := pgxmock.NewRows([]string{"id", "title", "desc", "poster_url", "link"})
	rowsRecom.AddRow(byteId, []uint8(f.BannersList[0].Title), []uint8(f.BannersList[0].Description), []uint8(f.BannersList[0].PictureURL), []uint8(f.BannersList[0].Link))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetBanners)).WillReturnRows(rowsRecom)
	pool.ExpectCommit()

	actual, err := repository.GetBanners()
	assert.NoError(t, err)
	assert.Equal(t, f, actual)
}

func TestGetGenresSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()

	f := domain.GenresList{
		GenresList: []domain.Genre{{
			Id:         1,
			Name:       "Test_film",
			PictureURL: "/pic/TestPoster.webp"}},
	}
	byteId := make([]byte, 4)
	binary.BigEndian.PutUint32(byteId, uint32(f.GenresList[0].Id))

	rowsRecom := pgxmock.NewRows([]string{"id", "title", "poster_url"})
	rowsRecom.AddRow(byteId, []uint8(f.GenresList[0].Name), []uint8(f.GenresList[0].PictureURL))

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetGenres)).WillReturnRows(rowsRecom)
	pool.ExpectCommit()

	actual, err := repository.GetGenres()
	assert.NoError(t, err)
	assert.Equal(t, f, actual)
}

func TestGetFilmByGenresSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()

	f := domain.FilmList{
		FilmList: []domain.Film{{
			Id:            1,
			Title:         "Test_film",
			TitleOriginal: "bnm",
			ReleaseYear:   2021,
			Description:   "asdf",
			Rating:        2.3,
			PosterUrl:     "/pic/TestPoster.webp",
			PremiereRu:    "2000-01-01"}},
		MoreAvailable: false,
		FilmTotal:     1,
		CurrentLimit:  10,
		CurrentSkip:   10,
	}
	g := domain.GenreFilmList{
		Id:        1,
		Name:      "genre",
		FilmsList: f,
	}
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, f.FilmList[0].Id)
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))
	byteRating := make([]byte, 8)
	binary.BigEndian.PutUint64(byteRating, math.Float64bits(f.FilmList[0].Rating))
	byteRelease := make([]byte, 4)
	binary.BigEndian.PutUint32(byteRelease, uint32(f.FilmList[0].ReleaseYear))
	bytePrem := make([]byte, 4)
	binary.BigEndian.PutUint32(bytePrem, 0)

	rowsGenre := pgxmock.NewRows([]string{"name"}).AddRow([]uint8("genre"))
	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	rowsRate := pgxmock.NewRows([]string{"rating"}).AddRow(byteRating)
	rowsRecom := pgxmock.NewRows([]string{"id", "title", "title_or", "release", "desc", "poster_url", "premiere"})
	rowsRecom.AddRow(byteId, []uint8(f.FilmList[0].Title), []uint8(f.FilmList[0].TitleOriginal), byteRelease, []uint8(f.FilmList[0].Description), []uint8(f.FilmList[0].PosterUrl), bytePrem)

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetGenreName)).WithArgs(uint64(1)).WillReturnRows(rowsGenre)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFilmsByGenreID)).WithArgs(uint64(1)).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmsByGenreID)).WithArgs(uint64(1), 10, 0).WillReturnRows(rowsRecom)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmRating)).WillReturnRows(rowsRate)
	pool.ExpectCommit()

	actual, err := repository.GetFilmsByGenre(uint64(1), 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, g, actual)
}

func TestGetFilmByMYSucess(t *testing.T) {
	mdb, pool, err := MockDatabase()
	assert.Equal(t, nil, err, "create a mock")
	repository := NewFilmRepository(mdb)
	defer pool.Close()

	f := domain.FilmList{
		FilmList: []domain.Film{{
			Id:            1,
			Title:         "Test_film",
			TitleOriginal: "bnm",
			ReleaseYear:   2021,
			Description:   "asdf",
			Rating:        2.3,
			PosterUrl:     "/pic/TestPoster.webp",
			PremiereRu:    "2000-01-01"}},
		MoreAvailable: false,
		FilmTotal:     1,
		CurrentLimit:  10,
		CurrentSkip:   10,
	}
	byteId := make([]byte, 8)
	binary.BigEndian.PutUint64(byteId, f.FilmList[0].Id)
	byteCount := make([]byte, 8)
	binary.BigEndian.PutUint64(byteCount, uint64(1))
	byteRating := make([]byte, 8)
	binary.BigEndian.PutUint64(byteRating, math.Float64bits(f.FilmList[0].Rating))
	byteRelease := make([]byte, 4)
	binary.BigEndian.PutUint32(byteRelease, uint32(f.FilmList[0].ReleaseYear))
	bytePrem := make([]byte, 4)
	binary.BigEndian.PutUint32(bytePrem, 0)

	rowsCount := pgxmock.NewRows([]string{"count"}).AddRow(byteCount)
	rowsRate := pgxmock.NewRows([]string{"rating"}).AddRow(byteRating)
	rowsRecom := pgxmock.NewRows([]string{"id", "title", "title_or", "release", "desc", "poster_url", "premiere"})
	rowsRecom.AddRow(byteId, []uint8(f.FilmList[0].Title), []uint8(f.FilmList[0].TitleOriginal), byteRelease, []uint8(f.FilmList[0].Description), []uint8(f.FilmList[0].PosterUrl), bytePrem)

	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryCountFilmsByMonthYear)).WithArgs(2, 2012).WillReturnRows(rowsCount)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmsByMonthYear)).WithArgs(2, 2012, 10, 0).WillReturnRows(rowsRecom)
	pool.ExpectCommit()
	pool.ExpectBegin()
	pool.ExpectQuery(regexp.QuoteMeta(queryGetFilmRating)).WillReturnRows(rowsRate)
	pool.ExpectCommit()

	actual, err := repository.GetFilmsByMonthYear(2, 2012, 10, 0)
	assert.NoError(t, err)
	assert.Equal(t, f, actual)
}
