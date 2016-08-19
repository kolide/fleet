package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
	// create the test datastore and server
	ds := createTestUsers(t, createTestDatastore(t))
	server := createTestServer(ds)

	////////////////////////////////////////////////////////////////////////////
	// log-in with a user
	////////////////////////////////////////////////////////////////////////////

	// log in with admin test user
	response := makeRequest(
		t,
		server,
		"POST",
		"/api/v1/kolide/login",
		CreateUserRequestBody{
			Username: "user1",
			Password: "foobar",
		},
		"",
	)
	assert.Equal(t, http.StatusOK, response.Code)

	// ensure that a non-empty cookie was in-fact set
	userCookie := response.Header().Get("Set-Cookie")
	assert.NotEmpty(t, userCookie)

	////////////////////////////////////////////////////////////////////////////
	// get the info of user2 from user1's account
	////////////////////////////////////////////////////////////////////////////

	user2, err := ds.User("user2")
	assert.Nil(t, err)

	response = makeRequest(
		t,
		server,
		"POST",
		"/api/v1/kolide/user",
		GetUserRequestBody{
			ID: user2.ID,
		},
		userCookie,
	)
	assert.Equal(t, http.StatusOK, response.Code)

	var user2Info GetUserResponseBody
	err = json.NewDecoder(response.Body).Decode(&user2Info)
	assert.Nil(t, err)

	assert.True(t, user2Info.NeedsPasswordReset)
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

func TestModifyUser(t *testing.T) {
}

func TestChangePassword(t *testing.T) {
}

func TestSetUserAdminState(t *testing.T) {
}

func TestSetUserEnabledState(t *testing.T) {
}

func TestDeleteSession(t *testing.T) {
}

func TestDeleteSessionForUser(t *testing.T) {
}

func TestGetInfoAboutSession(t *testing.T) {
}

func TestGetInfoAboutSessionsForUser(t *testing.T) {
}

func TestResetUserPassword(t *testing.T) {
}

func TestVerifyPasswordRequest(t *testing.T) {
}

func TestDeletePasswordResetRequest(t *testing.T) {
}
