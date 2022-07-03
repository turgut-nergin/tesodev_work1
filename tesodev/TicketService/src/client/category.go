package client

import (
	"encoding/json"
	"net/http"

	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
)

func (c Client) GetCategory(id string) (*models.Category, *errors.Error) {
	params := map[string]string{
		"categoryId": id,
	}

	response, err := c.do(http.MethodGet, "category", params)
	if err != nil {
		error := errors.UnknownError.WrapErrorCode(3035).WrapDesc(err.Error())
		return nil, error
	}

	defer response.Body.Close()

	if err != nil {
		error := errors.UnknownError.WrapErrorCode(3036).WrapDesc(err.Error())
		return nil, error
	}

	body := response.Body
	var category models.Category

	if err := json.NewDecoder(body).Decode(&category); err != nil {

		return nil, errors.UnknownError.WrapErrorCode(3037).WrapDesc(err.Error())
	}

	return &category, nil
}
