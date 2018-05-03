package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

// Logout attempts to logout to the current Fleet instance.
func (c *Client) Logout() error {
	response, err := c.AuthenticatedDo("POST", "/api/v1/kolide/logout", nil)
	if err != nil {
		return errors.Wrap(err, "error making request")
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Received HTTP %d instead of HTTP 200", response.StatusCode)
	}

	responeBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "error reading response body")
	}

	var responseBody logoutResponse
	err = json.Unmarshal(responeBytes, &responseBody)
	if err != nil {
		return errors.Wrap(err, "error decoding HTTP response body")
	}

	if responseBody.Err != nil {
		return errors.Wrap(err, "error logging out of Fleet")
	}

	return nil
}
