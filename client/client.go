package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/service"
	"github.com/pkg/errors"
)

type Client struct {
	addr string
	http *http.Client
}

func New(addr string, insecureSkipVerify bool) (*Client, error) {
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

// Setup attempts to setup the current Fleet instance. If setup is successful,
// an auth token is returned.
func (c *Client) Setup(email, password, org string) (string, error) {
	t := true
	params := service.SetupRequest{
		Admin: &kolide.UserPayload{
			Admin:    &t,
			Username: &email,
			Email:    &email,
			Password: &password,
		},
		OrgInfo: &kolide.OrgInfo{
			OrgName: &org,
		},
		KolideServerURL: &c.addr,
	}

	response, err := c.Do("POST", "/api/v1/setup", params)
	if err != nil {
		return "", errors.Wrap(err, "error making request")
	}
	defer response.Body.Close()

	// If setup has already been completed, Kolide Fleet will not serve the
	// setup route, which will cause the request to 404
	if response.StatusCode == http.StatusNotFound {
		return "", setupAlready()
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Received HTTP %d instead of HTTP 200", response.StatusCode)
	}

	responeBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading response body")
	}

	var responseBody service.SetupResponse
	err = json.Unmarshal(responeBytes, &responseBody)
	if err != nil {
		return "", errors.Wrap(err, "error decoding HTTP response body")
	}

	if responseBody.Err != nil {
		return "", errors.Wrap(err, "error setting up fleet instance")
	}

	return *responseBody.Token, nil
}
