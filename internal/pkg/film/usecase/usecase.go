package usecase

import (
	"2021_2_MAMBa/internal/pkg/domain"
	customErrors "2021_2_MAMBa/internal/pkg/domain/errors"
)

type FilmUsecase struct {
	FilmRepo domain.FilmRepository
}

func NewFilmUsecase(u domain.FilmRepository) domain.FilmUsecase {
	return &FilmUsecase{
		FilmRepo: u,
	}
}

func (uc *FilmUsecase) GetPopularFilms() (domain.FilmList, error) {
	filmsList, err := uc.FilmRepo.GetPopularFilms()
	if err != nil {
		return domain.FilmList{}, err
	}
	return filmsList, nil
}

func (uc *FilmUsecase) GetBanners() (domain.BannersList, error) {
	bannersList, err := uc.FilmRepo.GetBanners()
	if err != nil {
		return domain.BannersList{}, err
	}
	return bannersList, nil
}

func (uc *FilmUsecase) GetGenres() (domain.GenresList, error) {
	genresList, err := uc.FilmRepo.GetGenres()
	if err != nil {
		return domain.GenresList{}, err
	}
	return genresList, nil
}

func (uc *FilmUsecase) GetFilmsByGenre(genreID uint64, limit int, skip int) (domain.GenreFilmList, error) {
	genreFilmList, err := uc.FilmRepo.GetFilmsByGenre(genreID, limit, skip)
	if err != nil {
		return domain.GenreFilmList{}, err
	}
	return genreFilmList, nil
}

func (uc *FilmUsecase) BookmarkFilm(userID uint64, filmID uint64, bookmarked bool) error {
	err := uc.FilmRepo.BookmarkFilm(userID, filmID, bookmarked)
	if err != nil {
		return err
	}
	return nil
}

func (uc *FilmUsecase) GetFilmsByMonthYear(month int, year int, limit int, skip int) (domain.FilmList, error) {
	filmList, err := uc.FilmRepo.GetFilmsByMonthYear(month, year, limit, skip)
	if err != nil {
		return domain.FilmList{}, err
	}
	return filmList, nil
}

func (uc *FilmUsecase) LoadUserBookmarks(userID uint64, skip int, limit int) (domain.FilmBookmarks, error) {
	filmIdxList, err := uc.FilmRepo.LoadUserBookmarkedFilmsID(userID, skip, limit)
	if err != nil {
		return domain.FilmBookmarks{}, err
	}
	filmsList := make([]domain.Film, 0)
	for _, filmID := range filmIdxList {
		film, err := uc.FilmRepo.GetFilm(filmID)
		if err != nil {
			return domain.FilmBookmarks{}, err
		}
		filmsList = append(filmsList, film)
	}
	countBookmarks, err := uc.FilmRepo.CountBookmarks(userID)
	if err != nil {
		return domain.FilmBookmarks{}, err
	}
	if skip >= countBookmarks && skip != 0 {
		return domain.FilmBookmarks{}, customErrors.ErrorSkip
	}

	moreAvailable := skip+limit < countBookmarks

	bookmarks := domain.FilmBookmarks{
		FilmsList:     filmsList,
		MoreAvailable: moreAvailable,
		FilmsTotal:    countBookmarks,
		CurrentSort:   "",
		CurrentLimit:  limit,
		CurrentSkip:   skip + limit,
	}
	return bookmarks, nil
}

func (uc *FilmUsecase) GetFilm(userID, filmID uint64, skipReviews int, limitReviews int, skipRecommend int, limitRecommend int) (domain.FilmPageInfo, error) {
	film, err := uc.FilmRepo.GetFilm(filmID)
	if err != nil {
		return domain.FilmPageInfo{}, err
	}
	Reviews, err := uc.FilmRepo.GetFilmReviews(filmID, skipReviews, limitReviews)
	if err != nil {
		return domain.FilmPageInfo{}, err
	}
	Recommended, err := uc.FilmRepo.GetFilmRecommendations(filmID, skipReviews, limitRecommend)
	if err != nil {
		return domain.FilmPageInfo{}, err
	}

	myReview := domain.Review{}
	if userID != 0 {
		myReview, err = uc.FilmRepo.GetMyReview(filmID, userID)
		if err != nil && err != customErrors.ErrorNoReviewForFilm {
			return domain.FilmPageInfo{}, err
		}
	}

	bookmarked := false
	if userID != 0 {
		bookmarked, err = uc.FilmRepo.CheckFilmBookmarked(userID, filmID)
		if err != nil {
			return domain.FilmPageInfo{}, err
		}
	}

	result := domain.FilmPageInfo{
		FilmMain:        &film,
		Reviews:         Reviews,
		Recommendations: Recommended,
		MyReview:        myReview,
		Bookmarked:      bookmarked,
	}
	return result, nil
}

func (uc *FilmUsecase) PostRating(id uint64, authorId uint64, rating float64) (float64, error) {
	rating, err := uc.FilmRepo.PostRating(id, authorId, rating)
	if err != nil {
		return 0, err
	}
	return rating, nil
}

func (uc *FilmUsecase) LoadFilmReviews(id uint64, skip int, limit int) (domain.FilmReviews, error) {
	Reviews, err := uc.FilmRepo.GetFilmReviews(id, skip, limit)
	if err != nil {
		return domain.FilmReviews{}, err
	}
	return Reviews, nil
}
func (uc *FilmUsecase) LoadFilmRecommendations(id uint64, skip int, limit int) (domain.FilmRecommendations, error) {
	Recommendations, err := uc.FilmRepo.GetFilmRecommendations(id, skip, limit)
	if err != nil {
		return domain.FilmRecommendations{}, err
	}
	return Recommendations, nil
}
func (uc *FilmUsecase) LoadMyReview(id uint64, authorId uint64) (domain.Review, error) {
	myRev, err := uc.FilmRepo.GetMyReview(id, authorId)
	if err != nil {
		return domain.Review{}, err
	}
	return myRev, nil
}

func (uc *FilmUsecase) GetRandomFilms(genre1 uint64, genre2 uint64, genre3 uint64, dateStart int, dateEnd int) (domain.FilmList, error) {
	list, err := uc.FilmRepo.GetRandomFilms(genre1, genre2, genre3, dateStart, dateEnd)
	if err != nil {
		return domain.FilmList{}, err
	}
	return list, nil
}
