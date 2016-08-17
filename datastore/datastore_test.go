package datastore

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/app"
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

func testCreateUser(t *testing.T, db app.UserStore) {
	var createTests = []struct {
		username, password, email string
		isAdmin, passwordReset    bool
	}{
		{"marpaia", "foobar", "mike@kolide.co", true, false},
		{"jason", "foobar", "jason@kolide.co", true, false},
	}

	for _, tt := range createTests {
		u, err := app.NewUser(tt.username, tt.password, tt.email, tt.isAdmin, tt.passwordReset)
		if err != nil {
			t.Fatal(err)
		}

		user, err := db.NewUser(u)
		if err != nil {
			t.Fatal(err)
		}

		verify, err := db.User(tt.username)
		if err != nil {
			t.Fatal(err)
		}

		if verify.ID != user.ID {
			t.Fatalf("expected %q, got %q", user.ID, verify.ID)
		}

		if verify.Username != tt.username {
			t.Errorf("expected username: %s, got %s", tt.username, verify.Username)
		}

		if verify.Email != tt.email {
			t.Errorf("expected email: %s, got %s", tt.email, verify.Email)
		}

		if verify.Admin != tt.isAdmin {
			t.Errorf("expected email: %s, got %s", tt.email, verify.Email)
		}
	}
}

func createTestUsers(t *testing.T, db app.UserStore) []*app.User {
	var createTests = []struct {
		username, password, email string
		isAdmin, passwordReset    bool
	}{
		{"marpaia", "foobar", "mike@kolide.co", true, false},
		{"jason", "foobar", "jason@kolide.co", false, false},
	}

	var users []*app.User
	for _, tt := range createTests {
		u, err := app.NewUser(tt.username, tt.password, tt.email, tt.isAdmin, tt.passwordReset)
		if err != nil {
			t.Fatal(err)
		}

		user, err := db.NewUser(u)
		if err != nil {
			t.Fatal(err)
		}

		users = append(users, user)
	}
	if len(users) == 0 {
		t.Fatal("expected a list of users, got 0")
	}
	return users
}

func testSaveUser(t *testing.T, db app.UserStore) {
	users := createTestUsers(t, db)
	testAdminAttribute(t, db, users)
	testEmailAttribute(t, db, users)
	testPasswordAttribute(t, db, users)
}

func testPasswordAttribute(t *testing.T, db app.UserStore, users []*app.User) {
	for _, user := range users {
		user.Password = []byte(randomString(8))
		err := db.SaveUser(user)
		if err != nil {
			t.Fatalf("failed to save user %s", user.Name)
		}

		verify, err := db.User(user.Username)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(verify.Password, user.Password) {
			t.Errorf("expected password attribute to be %v, got %v", user.Password, verify.Password)
		}

	}
}

func testEmailAttribute(t *testing.T, db app.UserStore, users []*app.User) {
	for _, user := range users {
		user.Email = fmt.Sprintf("test.%s", user.Email)
		err := db.SaveUser(user)
		if err != nil {
			t.Fatalf("failed to save user %s", user.Name)
		}

		verify, err := db.User(user.Username)
		if err != nil {
			t.Fatal(err)
		}

		if verify.Email != user.Email {
			t.Errorf("expected admin attribute to be %v, got %v", user.Email, verify.Email)
		}
	}
}

func testAdminAttribute(t *testing.T, db app.UserStore, users []*app.User) {
	for _, user := range users {
		user.Admin = false
		err := db.SaveUser(user)
		if err != nil {
			t.Fatalf("failed to save user %s", user.Name)
		}

		verify, err := db.User(user.Username)
		if err != nil {
			t.Fatal(err)
		}

		if verify.Admin != user.Admin {
			t.Errorf("expected admin attribute to be %v, got %v", user.Admin, verify.Admin)
		}
	}
}

// setup creates a datastore for testing
func setup(t *testing.T) app.Datastore {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("error opening test db: %s", err)
	}
	ds := gormDB{DB: db}
	if err := ds.migrate(); err != nil {
		t.Fatal(err)
	}
	return ds
}

func teardown(t *testing.T, ds app.Datastore) {
	backend := ds.(gormDB)
	if err := backend.rollback(); err != nil {
		t.Fatal(err)
	}
}

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
