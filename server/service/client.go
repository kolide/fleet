package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
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

func (c *Client) Do(verb, path string, params interface{}) (*http.Response, error) {
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
	request.Header.Set("content-type", "application/json")
	request.Header.Set("accept", "application/json")

	return c.http.Do(request)
}

func (c *Client) SetToken(t string) {
	c.token = t
}
