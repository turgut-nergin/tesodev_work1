package models

import (
	"errors"
	"regexp"
)

func validateUserName(userName string) error {
	if len(userName) < 2 || len(userName) > 40 {
		return errors.New("The name field must be between 2-40 chars!")
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 2 || len(password) > 40 {
		return errors.New("The password field must be between 2-40 chars!")
	}
	return nil

}

func validateEmail(email string) error {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if isEmaailValid := emailRegex.MatchString(email); !isEmaailValid {
		return errors.New("invalid Email")
	}
	return nil
}

func validateType(userType string) error {
	if len(userType) < 2 || len(userType) > 5 {
		return errors.New("The password field must be between 2-5 chars!")
	}
	return nil
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
