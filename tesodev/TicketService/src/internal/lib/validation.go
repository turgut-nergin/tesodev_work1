package lib

import (
	"errors"

	"github.com/turgut-nergin/tesodev_work1/internal/models"
)

func validateSubject(subject string) error {
	if subject == "" {
		return errors.New("The subject field is required!")

	}

	return nil
}

func validateStatus(status string) error {
	if status == "" {
		return errors.New("The status field is required!")

	}
	return nil
}

func validateBody(body string) error {
	if body == "" {
		return errors.New("The status field is required!")
	}
	return nil
}

func Validate(ticket models.TicketRequest) error {

	if err := validateBody(ticket.Body); err != nil {
		return err
	}

	if err := validateStatus(ticket.Status); err != nil {
		return err
	}

	if err := validateSubject(ticket.Subject); err != nil {
		return err
	}

	return nil

}
