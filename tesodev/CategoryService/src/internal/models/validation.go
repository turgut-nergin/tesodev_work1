package models

import (
	"errors"
)

func validateUser(name string) error {
	if len(name) < 2 || len(name) > 40 {
		return errors.New("The name field must be between 2-40 chars!")
	}
	return nil
}
