package datastore

import (
	"fmt"
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
)

// TestUser tests the UserStore interface
// this test uses the default testing backend
func TestCreateUser(t *testing.T) {
	db := setup(t)
	defer teardown(t, db)

	testCreateUser(t, db)
}

func TestSaveUser(t *testing.T) {
	db := setup(t)
	defer teardown(t, db)

	testSaveUser(t, db)
}

func testCreateUser(t *testing.T, db kolide.UserStore) {
	var createTests = []struct {
		username, password, email string
		isAdmin, passwordReset    bool
	}{
		{"marpaia", "foobar", "mike@kolide.co", true, false},
		{"jason", "foobar", "jason@kolide.co", true, false},
	}

	for _, tt := range createTests {
		u := &kolide.User{
			Username: tt.username,
			Password: []byte(tt.password),
			Admin:    tt.isAdmin,
			AdminForcedPasswordReset: tt.passwordReset,
			Email: tt.email,
		}
		user, err := db.NewUser(u)
		assert.Nil(t, err)

		verify, err := db.User(tt.username)
		assert.Nil(t, err)

		assert.Equal(t, user.ID, verify.ID)
		assert.Equal(t, tt.username, verify.Username)
		assert.Equal(t, tt.email, verify.Email)
		assert.Equal(t, tt.email, verify.Email)
	}
}

func TestUserByID(t *testing.T) {
	db := setup(t)
	defer teardown(t, db)

	testUserByID(t, db)
}

func testUserByID(t *testing.T, db kolide.UserStore) {
	users := createTestUsers(t, db)
	for _, tt := range users {
		returned, err := db.UserByID(tt.ID)
		assert.Nil(t, err)
		assert.Equal(t, tt.ID, returned.ID)
	}

	// test missing user
	_, err := db.UserByID(10000000000)
	assert.NotNil(t, err)
}

func createTestUsers(t *testing.T, db kolide.UserStore) []*kolide.User {
	var createTests = []struct {
		username, password, email string
		isAdmin, passwordReset    bool
	}{
		{"marpaia", "foobar", "mike@kolide.co", true, false},
		{"jason", "foobar", "jason@kolide.co", false, false},
	}

	var users []*kolide.User
	for _, tt := range createTests {
		u := &kolide.User{
			Username: tt.username,
			Password: []byte(tt.password),
			Admin:    tt.isAdmin,
			AdminForcedPasswordReset: tt.passwordReset,
			Email: tt.email,
		}

		user, err := db.NewUser(u)
		assert.Nil(t, err)

		users = append(users, user)
	}
	assert.NotEmpty(t, users)
	return users
}

func testSaveUser(t *testing.T, db kolide.UserStore) {
	users := createTestUsers(t, db)
	testAdminAttribute(t, db, users)
	testEmailAttribute(t, db, users)
	testPasswordAttribute(t, db, users)
}

func testPasswordAttribute(t *testing.T, db kolide.UserStore, users []*kolide.User) {
	for _, user := range users {
		randomText, err := generateRandomText(8)
		assert.Nil(t, err)
		user.Password = []byte(randomText)
		err = db.SaveUser(user)
		assert.Nil(t, err)

		verify, err := db.User(user.Username)
		assert.Nil(t, err)
		assert.Equal(t, user.Password, verify.Password)
	}
}

func testEmailAttribute(t *testing.T, db kolide.UserStore, users []*kolide.User) {
	for _, user := range users {
		user.Email = fmt.Sprintf("test.%s", user.Email)
		err := db.SaveUser(user)
		assert.Nil(t, err)

		verify, err := db.User(user.Username)
		assert.Nil(t, err)
		assert.Equal(t, user.Email, verify.Email)
	}
}

func testAdminAttribute(t *testing.T, db kolide.UserStore, users []*kolide.User) {
	for _, user := range users {
		user.Admin = false
		err := db.SaveUser(user)
		assert.Nil(t, err)

		verify, err := db.User(user.Username)
		assert.Nil(t, err)
		assert.Equal(t, user.Admin, verify.Admin)
	}
}
