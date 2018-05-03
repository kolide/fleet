package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Login attempts to login to the current Fleet instance. If setup is successful,
// an auth token is returned.
func (c *Client) Login(email, password string) (string, error) {
	params := loginRequest{
		Username: email,
		Password: password,
	}

	response, err := c.Do("POST", "/api/v1/kolide/login", params)
	if err != nil {
		return "", errors.Wrap(err, "error making request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Received HTTP %d instead of HTTP 200", response.StatusCode)
	}

	responeBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading response body")
	}

	var responseBody loginResponse
	err = json.Unmarshal(responeBytes, &responseBody)
	if err != nil {
		return "", errors.Wrap(err, "error decoding HTTP response body")
	}

	if responseBody.Err != nil {
		return "", errors.Wrap(err, "error setting up fleet instance")
	}

	return responseBody.Token, nil
}
