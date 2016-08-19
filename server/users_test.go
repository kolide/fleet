package server

import (
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.Set("session.cookie_name", "KolideSession")
	viper.Set("session.key_size", 24)
}

func TestGetUser(t *testing.T) {

}

func TestCreateUser(t *testing.T) {
	// create the test datastore and server
	ds := createTestUsers(t, createTestDatastore(t))
	server := createTestServer(ds)

	////////////////////////////////////////////////////////////////////////////
	// log-in with an admin
	////////////////////////////////////////////////////////////////////////////

	// log in with admin test user
	response := makeRequest(
		t,
		server,
		"POST",
		"/api/v1/kolide/login",
		CreateUserRequestBody{
			Username: "admin1",
			Password: "foobar",
		},
		"",
	)
	assert.Equal(t, http.StatusOK, response.Code)

	// ensure that a non-empty cookie was in-fact set
	adminCookie := response.Header().Get("Set-Cookie")
	assert.NotEmpty(t, adminCookie)

	////////////////////////////////////////////////////////////////////////////
	// create test user account
	////////////////////////////////////////////////////////////////////////////

	// make the request to create the new user and verify that it succeeded
	response = makeRequest(
		t,
		server,
		"PUT",
		"/api/v1/kolide/user",
		CreateUserRequestBody{
			Username:           "tester",
			Password:           "temp",
			Email:              "tester@kolide.co",
			Admin:              false,
			NeedsPasswordReset: true,
		},
		adminCookie,
	)
	assert.Equal(t, http.StatusOK, response.Code)

	// ensure that the new user was created in the database
	_, err := ds.User("tester")
	assert.Nil(t, err)
}
