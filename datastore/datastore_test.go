package datastore

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/app"
)

// TestUser tests the UserStore interface
// this test uses the default testing backend
func TestUser(t *testing.T) {
	db := setup(t)
	defer teardown(t, db)

	testUser(t, db)
}

func testUser(t *testing.T, db app.UserStore) {
	var userTests = []struct {
		username, password, email string
		isAdmin, passwordReset    bool
	}{
		{"marpaia", "foobar", "mike@kolide.co", true, false},
		{"jason", "foobar", "jason@kolide.co", true, false},
	}

	for _, tt := range userTests {
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
