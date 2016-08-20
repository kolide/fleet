package server

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoginAndLogout(t *testing.T) {
	// create the test datastore and server
	ds := createTestUsers(t, createTestDatastore(t))
	server := createTestServer(ds)

	admin, err := ds.User("admin1")

	// ensure that there are no sessions in the database for our test user
	sessions, err := ds.FindAllSessionsForUser(admin.ID)
	assert.Nil(t, err)
	assert.Len(t, sessions, 0)

	////////////////////////////////////////////////////////////////////////////
	// Test logging in
	////////////////////////////////////////////////////////////////////////////

	// log in with test user created above
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
	cookie := response.Header().Get("Set-Cookie")
	assert.NotEmpty(t, cookie)

	// ensure that a session was created for our test user and stored
	sessions, err = ds.FindAllSessionsForUser(admin.ID)
	assert.Nil(t, err)
	assert.Len(t, sessions, 1)

	// ensure the session key is not blank
	assert.NotEqual(t, "", sessions[0].Key)

	////////////////////////////////////////////////////////////////////////////
	// Test logging out
	////////////////////////////////////////////////////////////////////////////

	// log out our test user
	response = makeRequest(
		t,
		server,
		"GET",
		"/api/v1/kolide/logout",
		nil,
		cookie,
	)
	assert.Equal(t, http.StatusOK, response.Code)

	// ensure that a cookie was actually set to erase the current cookie
	assert.Equal(t, "KolideSession=", response.Header().Get("Set-Cookie"))

	// ensure that our user's session was deleted from the store
	sessions, err = ds.FindAllSessionsForUser(admin.ID)
	assert.Nil(t, err)
	assert.Len(t, sessions, 0)
}

func TestNeedsPasswordReset(t *testing.T) {
	// create the test datastore and server
	ds := createTestUsers(t, createTestDatastore(t))
	server := createTestServer(ds)

	////////////////////////////////////////////////////////////////////////////
	// log-in with a user which needs a password reset
	////////////////////////////////////////////////////////////////////////////

	// log in with admin test user
	response := makeRequest(
		t,
		server,
		"POST",
		"/api/v1/kolide/login",
		CreateUserRequestBody{
			Username: "user2",
			Password: "foobar",
		},
		"",
	)
	assert.Equal(t, http.StatusOK, response.Code)

	// ensure that a non-empty cookie was in-fact set
	userCookie := response.Header().Get("Set-Cookie")
	assert.NotEmpty(t, userCookie)

	////////////////////////////////////////////////////////////////////////////
	// get the info of user1 from user2's account
	////////////////////////////////////////////////////////////////////////////

	user1, err := ds.User("user1")
	assert.Nil(t, err)

	response = makeRequest(
		t,
		server,
		"POST",
		"/api/v1/kolide/user",
		GetUserRequestBody{
			ID: user1.ID,
		},
		userCookie,
	)
	assert.Equal(t, http.StatusUnauthorized, response.Code)
}

func TestDeleteSession(t *testing.T) {
}

func TestDeleteSessionForUser(t *testing.T) {
}

func TestGetInfoAboutSession(t *testing.T) {
}

func TestGetInfoAboutSessionsForUser(t *testing.T) {
}
