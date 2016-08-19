package server

import (
	"net/http"
	"testing"

	"github.com/kolide/kolide-ose/kolide"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	viper.Set("session.cookie_name", "KolideSession")
	viper.Set("session.key_size", 24)
}

func TestLoginAndLogout(t *testing.T) {
	// create the test datastore and server
	ds := createTestDatastore()
	server := createTestServer(ds)

	// popuplate the database with test data
	admin, err := kolide.NewUser(
		"admin",
		"foobar",
		"admin@kolide.co",
		true,
		false,
	)
	assert.Nil(t, err)

	admin, err = ds.NewUser(admin)
	assert.Nil(t, err)
	assert.NotZero(t, admin.ID)

	// ensure that there are no sessions in the database for our test user
	sessions, err := ds.FindAllSessionsForUser(admin.ID)
	assert.Nil(t, err)
	assert.Len(t, sessions, 0)

	////////////////////////////////////////////////////////////////////////////
	// Test logging in
	////////////////////////////////////////////////////////////////////////////

	// log in with test user created above
	response := makeRequest(
		server,
		"POST",
		"/api/v1/kolide/login",
		CreateUserRequestBody{
			Username: "admin",
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
