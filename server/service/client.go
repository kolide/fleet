package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	addr  string
	token string
	http  *http.Client
}

func NewClient(addr string, insecureSkipVerify bool) (*Client, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
		},
	}

	return &Client{
		addr: addr,
		http: httpClient,
	}, nil
}

func (c *Client) DoWithHeaders(verb, path string, params interface{}, headers map[string]string) (*http.Response, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, errors.Wrap(err, "error marshaling json")
	}

	request, err := http.NewRequest(
		verb,
		c.addr+path,
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request object")
	}
	for k, v := range headers {
		request.Header.Set(k, v)
	}

	return c.http.Do(request)
}

func (c *Client) Do(verb, path string, params interface{}) (*http.Response, error) {
	headers := map[string]string{
		"Content-type": "application/json",
		"Accept":       "application/json",
	}

	return c.DoWithHeaders(verb, path, params, headers)
}

func (c *Client) AuthenticatedDo(verb, path string, params interface{}) (*http.Response, error) {
	if c.token == "" {
		return nil, errors.New("authentication token is empty")
	}

	headers := map[string]string{
		"Content-Type":  "application/json",
		"Accept":        "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", c.token),
	}

	return c.DoWithHeaders(verb, path, params, headers)
}

func (c *Client) SetToken(t string) {
	c.token = t
}
