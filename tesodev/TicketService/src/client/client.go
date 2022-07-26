package client

import (
	"fmt"
	"os"
	"time"

	"github.com/valyala/fasthttp"
)

type Client struct {
	Host   string
	Client *fasthttp.Client
}

func NewClient(host string) *Client {

	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")

	client := &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	return &Client{
		Host:   host,
		Client: client,
	}
}

func (c *Client) do(method string, endpoint string, params map[string]string, headers map[string]string) (*fasthttp.Response, error) {
	baseURL := fmt.Sprintf("%s/%s", c.Host, endpoint)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(baseURL)
	req.Header.SetMethod(method)
	q := req.URI().QueryArgs()
	req.Header.Add("Content-Type", "application/json")
	for key, val := range params {
		q.Set(key, val)
	}

	for key, val := range headers {
		req.Header.Add(key, val)
	}
	resp := fasthttp.AcquireResponse()

	err := c.Client.Do(req, resp)

	if err == nil {
		fmt.Printf("DEBUG Response: %s\n", resp.Body())
	} else {
		fmt.Fprintf(os.Stderr, "ERR Connection error: %v\n", err)
	}

	return resp, nil

}
