package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/turgut-nergin/tesodev_work1/internal/errors"
)

func (c Client) UserIsExist(id string) (*bool, *errors.Error) {
	params := map[string]string{
		"userId": id,
	}

	response, err := c.do(http.MethodGet, "user", params)
	if err != nil {
		error := errors.UnknownError.WrapErrorCode(3024).WrapDesc(err.Error())
		return nil, error
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		error := errors.UnknownError.WrapErrorCode(3025).WrapDesc(err.Error())
		return nil, error
	}

	isExist := strings.TrimSpace(string(body)) == "true"
	if isExist {
		return &isExist, nil
	}

	var result errors.Error
	if err := json.Unmarshal(body, &result); err != nil {
		error := errors.UnknownError.WrapErrorCode(3025).WrapDesc(err.Error())
		return nil, error
	}

	return nil, &result
}
