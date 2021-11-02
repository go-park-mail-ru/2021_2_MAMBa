package xss

import (
	"2021_2_MAMBa/internal/pkg/domain"
	"github.com/microcosm-cc/bluemonday"
)

func SanitizeUser(user *domain.User) {
	sanitizer := bluemonday.UGCPolicy()
	user.Email = sanitizer.Sanitize(user.Email)
	user.FirstName = sanitizer.Sanitize(user.FirstName)
	user.Surname = sanitizer.Sanitize(user.Surname)
	user.ProfilePic = sanitizer.Sanitize(user.ProfilePic)
}

func SanitizeProfile(profile *domain.Profile) {
	sanitizer := bluemonday.UGCPolicy()
	profile.Email = sanitizer.Sanitize(profile.Email)
	profile.FirstName = sanitizer.Sanitize(profile.FirstName)
	profile.Surname = sanitizer.Sanitize(profile.Surname)
	profile.PictureUrl = sanitizer.Sanitize(profile.PictureUrl)
}

func SanitizeReview(review *domain.Review) {
	sanitizer := bluemonday.UGCPolicy()
	review.ReviewText = sanitizer.Sanitize(review.ReviewText)
}
