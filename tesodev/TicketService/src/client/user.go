package client

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/turgut-nergin/tesodev_work1/internal/errors"
	"github.com/valyala/fasthttp"
)

func (c Client) UserIsExist(id string) (*bool, *errors.Error) {
	params := map[string]string{
		"userId": id,
	}

	response, err := c.do(fasthttp.MethodGet, "user", params, nil)
	if err != nil {
		error := errors.UnknownError.WrapErrorCode(3024).WrapDesc(err.Error())
		return nil, error
	}

	body := response.Body()

	isExist := strings.TrimSpace(string(body)) == "true"
	if isExist {
		return &isExist, nil
	}
	fmt.Println("%s", body)
	var result errors.Error
	if err := json.Unmarshal(body, &result); err != nil {
		error := errors.UnknownError.WrapErrorCode(3025).WrapDesc(err.Error())
		return nil, error
	}

	return nil, &result
}
