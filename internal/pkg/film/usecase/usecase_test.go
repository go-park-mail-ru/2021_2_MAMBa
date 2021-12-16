package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	mockF "2021_2_MAMBa/internal/pkg/film/repository/mock"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetPopular(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilm := mockF.NewMockFilmRepository(ctrl)
	usecase := NewFilmUsecase(mockFilm)
	testErr := errors.New("test")

	res := domain.FilmList{FilmList: []domain.Film{domain.Film{Id: 1}}, MoreAvailable:true, FilmTotal:1, CurrentLimit:10, CurrentSkip:10}
	mockFilm.EXPECT().GetPopularFilms().Return(res, nil)
	actual, err := usecase.GetPopularFilms()
	assert.Equal(t, res, actual)
	assert.Equal(t, nil, err)

	mockFilm.EXPECT().GetPopularFilms().Return(domain.FilmList{}, testErr)
	actual, err = usecase.GetPopularFilms()
	assert.Equal(t, domain.FilmList{}, actual)
	assert.Equal(t, testErr, err)
}

func TestGetBanners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilm := mockF.NewMockFilmRepository(ctrl)
	usecase := NewFilmUsecase(mockFilm)
	testErr := errors.New("test")

	res := domain.BannersList{BannersList: []domain.Banner{domain.Banner{Id: 1}}}
	mockFilm.EXPECT().GetBanners().Return(res, nil)
	actual, err := usecase.GetBanners()
	assert.Equal(t, res, actual)
	assert.Equal(t, nil, err)

	mockFilm.EXPECT().GetBanners().Return(domain.BannersList{}, testErr)
	actual, err = usecase.GetBanners()
	assert.Equal(t, testErr, err)
}

func TestGetGenres(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilm := mockF.NewMockFilmRepository(ctrl)
	usecase := NewFilmUsecase(mockFilm)
	testErr := errors.New("test")

	res := domain.GenresList{GenresList: []domain.Genre{domain.Genre{Id: 1}}}
	mockFilm.EXPECT().GetGenres().Return(res, nil)
	actual, err := usecase.GetGenres()
	assert.Equal(t, res, actual)
	assert.Equal(t, nil, err)

	mockFilm.EXPECT().GetGenres().Return(domain.GenresList{}, testErr)
	actual, err = usecase.GetGenres()
	assert.Equal(t, domain.GenresList{}, actual)
	assert.Equal(t, testErr, err)
}

func TestGetFilmsByGenre(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilm := mockF.NewMockFilmRepository(ctrl)
	usecase := NewFilmUsecase(mockFilm)
	testErr := errors.New("test")

	res := domain.GenreFilmList{Id: 1}
	mockFilm.EXPECT().GetFilmsByGenre(uint64(1),10,0).Return(res, nil)
	actual, err := usecase.GetFilmsByGenre(uint64(1),10,0)
	assert.Equal(t, res, actual)
	assert.Equal(t, nil, err)

	mockFilm.EXPECT().GetFilmsByGenre(uint64(1),10,0).Return(domain.GenreFilmList{}, testErr)
	actual, err = usecase.GetFilmsByGenre(uint64(1),10,0)
	assert.Equal(t, domain.GenreFilmList{}, actual)
	assert.Equal(t, testErr, err)
}

func TestGetFilmMY(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilm := mockF.NewMockFilmRepository(ctrl)
	usecase := NewFilmUsecase(mockFilm)
	testErr := errors.New("test")

	res := domain.FilmList{FilmList: []domain.Film{domain.Film{Id: 1}}, MoreAvailable:true, FilmTotal:1, CurrentLimit:10, CurrentSkip:10}
	mockFilm.EXPECT().GetFilmsByMonthYear(2,2010,0,10).Return(res, nil)
	actual, err := usecase.GetFilmsByMonthYear(2,2010,0,10)
	assert.Equal(t, res, actual)
	assert.Equal(t, nil, err)

	mockFilm.EXPECT().GetFilmsByMonthYear(2,2010,0,10).Return(domain.FilmList{}, testErr)
	actual, err = usecase.GetFilmsByMonthYear(2,2010,0,10)
	assert.Equal(t, domain.FilmList{}, actual)
	assert.Equal(t, testErr, err)
}


func TestLoadBMs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilm := mockF.NewMockFilmRepository(ctrl)
	usecase := NewFilmUsecase(mockFilm)
	testErr := errors.New("test")

	res := domain.FilmList{FilmList: []domain.Film{domain.Film{Id: 1}}, MoreAvailable:true, FilmTotal:1, CurrentLimit:10, CurrentSkip:10}
	mockFilm.EXPECT().LoadUserBookmarkedFilmsID(uint64(1), 0,10).Return([]uint64{1}, nil)
	mockFilm.EXPECT().CountBookmarks(uint64(1)).Return(1, nil)
	mockFilm.EXPECT().GetFilm(uint64(1)).Return(res.FilmList[0], nil)
	actual, err := usecase.LoadUserBookmarks(uint64(1), 0, 10)
	bms:= domain.FilmBookmarks{
		FilmsList:     res.FilmList,
		MoreAvailable: false,
		FilmsTotal:    1,
		CurrentSort:   "",
		CurrentLimit:  10,
		CurrentSkip:   10,
	}
	assert.Equal(t, bms, actual)
	assert.Equal(t, nil, err)

	mockFilm.EXPECT().LoadUserBookmarkedFilmsID(uint64(1), 0,10).Return([]uint64{1}, testErr)
	bms= domain.FilmBookmarks{}
	actual, err = usecase.LoadUserBookmarks(uint64(1), 0, 10)
	assert.Equal(t, bms, actual)
	assert.Equal(t, testErr, err)

	mockFilm.EXPECT().LoadUserBookmarkedFilmsID(uint64(1), 0,10).Return([]uint64{1}, nil)
	mockFilm.EXPECT().GetFilm(uint64(1)).Return(res.FilmList[0], testErr)
	bms= domain.FilmBookmarks{}
	actual, err = usecase.LoadUserBookmarks(uint64(1), 0, 10)
	assert.Equal(t, bms, actual)
	assert.Equal(t, testErr, err)

	mockFilm.EXPECT().LoadUserBookmarkedFilmsID(uint64(1), 0,10).Return([]uint64{1}, nil)
	mockFilm.EXPECT().GetFilm(uint64(1)).Return(res.FilmList[0], nil)
	mockFilm.EXPECT().CountBookmarks(uint64(1)).Return(1, testErr)
	bms= domain.FilmBookmarks{}
	actual, err = usecase.LoadUserBookmarks(uint64(1), 0, 10)
	assert.Equal(t, bms, actual)
	assert.Equal(t, testErr, err)
}

func TestGetFilm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockFilm := mockF.NewMockFilmRepository(ctrl)
	usecase := NewFilmUsecase(mockFilm)
	testErr := errors.New("test")

	res := domain.FilmPageInfo{
		FilmMain:        &domain.Film{Id: 1},
		Reviews:         domain.FilmReviews{ReviewTotal: 0},
		Recommendations: domain.FilmRecommendations{RecommendationTotal: 0},
		MyReview:        domain.Review{Id: 1},
		Bookmarked:      false,
	}

	mockFilm.EXPECT().GetFilm(uint64(1)).Return(domain.Film{Id: 1}, nil)
	mockFilm.EXPECT().GetFilmReviews(uint64(1), 0,10).Return(domain.FilmReviews{ReviewTotal: 0},nil)
	mockFilm.EXPECT().GetFilmRecommendations(uint64(1), 0,10).Return(domain.FilmRecommendations{RecommendationTotal: 0},nil)
	mockFilm.EXPECT().GetMyReview(uint64(1), uint64(1)).Return(domain.Review{Id: 1}, nil)
	mockFilm.EXPECT().CheckFilmBookmarked(uint64(1), uint64(1)).Return(false, nil)
	actual, err := usecase.GetFilm(uint64(1), uint64(1), 0,10,0,10)
	assert.Equal(t, res, actual)
	assert.Equal(t, nil, err)

	mockFilm.EXPECT().GetFilm(uint64(1)).Return(domain.Film{Id: 1}, testErr)
	actual, err = usecase.GetFilm(uint64(1), uint64(1), 0,10,0,10)
	assert.Equal(t, domain.FilmPageInfo{}, actual)
	assert.Equal(t, testErr, err)

	mockFilm.EXPECT().GetFilm(uint64(1)).Return(domain.Film{Id: 1}, nil)
	mockFilm.EXPECT().GetFilmReviews(uint64(1), 0,10).Return(domain.FilmReviews{ReviewTotal: 0},testErr)
	actual, err = usecase.GetFilm(uint64(1), uint64(1), 0,10,0,10)
	assert.Equal(t, domain.FilmPageInfo{}, actual)
	assert.Equal(t, testErr, err)

	mockFilm.EXPECT().GetFilm(uint64(1)).Return(domain.Film{Id: 1}, nil)
	mockFilm.EXPECT().GetFilmReviews(uint64(1), 0,10).Return(domain.FilmReviews{ReviewTotal: 0},nil)
	mockFilm.EXPECT().GetFilmRecommendations(uint64(1), 0,10).Return(domain.FilmRecommendations{RecommendationTotal: 0},testErr)
	actual, err = usecase.GetFilm(uint64(1), uint64(1), 0,10,0,10)
	assert.Equal(t, domain.FilmPageInfo{}, actual)
	assert.Equal(t, testErr, err)

	mockFilm.EXPECT().GetFilm(uint64(1)).Return(domain.Film{Id: 1}, nil)
	mockFilm.EXPECT().GetFilmReviews(uint64(1), 0,10).Return(domain.FilmReviews{ReviewTotal: 0},nil)
	mockFilm.EXPECT().GetFilmRecommendations(uint64(1), 0,10).Return(domain.FilmRecommendations{RecommendationTotal: 0},nil)
	mockFilm.EXPECT().GetMyReview(uint64(1), uint64(1)).Return(domain.Review{Id: 1}, testErr)
	actual, err = usecase.GetFilm(uint64(1), uint64(1), 0,10,0,10)
	assert.Equal(t, domain.FilmPageInfo{}, actual)
	assert.Equal(t, testErr, err)

	mockFilm.EXPECT().GetFilm(uint64(1)).Return(domain.Film{Id: 1}, nil)
	mockFilm.EXPECT().GetFilmReviews(uint64(1), 0,10).Return(domain.FilmReviews{ReviewTotal: 0},nil)
	mockFilm.EXPECT().GetFilmRecommendations(uint64(1), 0,10).Return(domain.FilmRecommendations{RecommendationTotal: 0},nil)
	mockFilm.EXPECT().GetMyReview(uint64(1), uint64(1)).Return(domain.Review{Id: 1}, nil)
	mockFilm.EXPECT().CheckFilmBookmarked(uint64(1), uint64(1)).Return(false, testErr)
	actual, err = usecase.GetFilm(uint64(1), uint64(1), 0,10,0,10)
	assert.Equal(t, domain.FilmPageInfo{}, actual)
	assert.Equal(t, testErr, err)
}