package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func validateUserName(userName string) error {
	return validation.Validate(userName, validation.Required, validation.Length(2, 100))
}

func validatePassword(password string) error {
	return validation.Validate(password, validation.Required, validation.Length(2, 50))
}

func validateEmail(email string) error {
	return validation.Validate(email, validation.Required, is.Email)
}

func validateType(userType string) error {
	return validation.Validate(userType, validation.Required)
}

func (user UserRequest) ValidateUpdate() error {
	if user.Password != "" {
		if err := validatePassword(user.Password); err != nil {
			return err
		}
	}

	if user.Email != "" {
		if err := validateEmail(user.Email); err != nil {
			return err
		}
	}

	if user.UserName != "" {
		if err := validateUserName(user.UserName); err != nil {
			return err
		}
	}

	if user.Type != "" {
		if err := validateType(user.Type); err != nil {
			return err
		}
	}
	return nil
}

func (user UserRequest) ValidateInsert() error {

	if err := validatePassword(user.Password); err != nil {
		return err
	}

	if err := validateEmail(user.Email); err != nil {
		return err
	}

	if err := validateUserName(user.UserName); err != nil {
		return err
	}

	if err := validateType(user.Type); err != nil {
		return err
	}
	return nil

}
