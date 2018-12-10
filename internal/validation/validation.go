package validation

import (
	"SimpleGame/internal/models"
	"github.com/asaskevich/govalidator"
)

func ValidUser(user *models.User) bool {
	validEmail := govalidator.IsEmail(user.Email)
	validPassword := govalidator.HasUpperCase(user.Password) && govalidator.HasLowerCase(user.Password) //&& govalidator.IsByteLength(user.Password, 6, 12)
	validNick := !govalidator.HasWhitespace(user.Nick)

	if validEmail && validPassword && validNick {
		return true
	} else {
		return false
	}
}
