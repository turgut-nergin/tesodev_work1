package user

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (c Client) UserIsExist(id string) (*bool, error) {
	params := map[string]string{
		"userId": id,
	}

	response, err := c.do(http.MethodGet, "user", params)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	isExist, err := strconv.ParseBool(strings.TrimSpace(string(body)))

	if err != nil {
		return nil, err
	}

	return &isExist, nil
}
