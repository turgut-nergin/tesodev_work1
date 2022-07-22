package client

import (
	"encoding/json"

	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/turgut-nergin/tesodev_work1/internal/models"
	"github.com/valyala/fasthttp"
)

func (c Client) GetCategory(id string) (*models.Category, *errors.Error) {
	params := map[string]string{
		"categoryId": id,
	}

	response, err := c.do(fasthttp.MethodGet, "category", params)
	defer fasthttp.ReleaseResponse(response)
	if err != nil {
		error := errors.UnknownError.WrapErrorCode(3035).WrapDesc(err.Error())
		return nil, error
	}

	body := response.Body()

	if response.StatusCode() == fasthttp.StatusOK {
		var category models.Category
		if err := json.Unmarshal(body, &category); err != nil {
			return nil, errors.UnknownError.WrapOperation("category client").WrapErrorCode(3037).WrapDesc(err.Error())
		}
		return &category, nil
	}

	var errResult errors.Error
	if err := json.Unmarshal(body, &errResult); err != nil {
		return nil, errors.UnknownError.WrapOperation("category client").WrapErrorCode(3077).WrapDesc(err.Error())
	}
	return nil, &errResult
}
