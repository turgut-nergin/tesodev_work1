package user

import (
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	Host   string
	Client *http.Client
}

func NewClient(host string) *Client {
	client := &http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			IdleConnTimeout:     10 * time.Second,
			MaxIdleConnsPerHost: 10,
		},
	}
	return &Client{
		Host:   host,
		Client: client,
	}
}

func (c *Client) do(method, endpoint string, params map[string]string) (*http.Response, error) {
	baseURL := fmt.Sprintf("%s/%s", c.Host, endpoint)
	fmt.Println(baseURL)
	req, err := http.NewRequest(method, baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	req.URL.RawQuery = q.Encode()
	return c.Client.Do(req)
}
