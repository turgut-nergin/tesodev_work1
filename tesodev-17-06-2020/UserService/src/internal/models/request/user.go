package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Type     string `json:"type"`
}

func (user User) validateUser() error {
	return validation.ValidateStruct(&user,
		validation.Field(&user.UserName, validation.Required),
		validation.Field(&user.Password, validation.Required),
		validation.Field(&user.Email, validation.Required, is.Email),
		validation.Field(&user.Type, validation.Required),
	)
}

func (user User) Validate() error {
	err := user.validateUser()
	return err
}
