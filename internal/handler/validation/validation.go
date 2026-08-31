package validation

import (
	"ladno-teams/internal/entity"
	"regexp"

	"github.com/go-playground/validator/v10"
)

const (
	ValidLogin    = "valid_login"
	ValidPassword = "valid_password"
	ValidRole     = "valid_role"
)

var loginRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)
var passwordRegex = regexp.MustCompile(`^.{6,128}$`)

func LoginValidation(fl validator.FieldLevel) bool {
	return loginRegex.MatchString(fl.Field().String())
}

func PasswordValidation(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}

func RoleValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	switch entity.Role(value) {
	case entity.RoleMaintainer, entity.RoleDeveloper, entity.RoleGuest:
		return true
	default:
		return false
	}
}
