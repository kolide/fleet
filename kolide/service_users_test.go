package kolide_test

import (
	"testing"

	"github.com/kolide/kolide-ose/datastore"
	"github.com/kolide/kolide-ose/kolide"
)

func TestUserService(t *testing.T) {
	ds, _ := datastore.New("mock", "")
	svc, _ := kolide.NewService(ds)

	userTests := []struct {
		user       *kolide.User
		password   string
		changePass string
	}{
		{user: &kolide.User{
			Username: "marpaia",
			Password: []byte("foobar"),
			Admin:    true,
			Enabled:  true,
		},
			password:   "foobar",
			changePass: "bazfoo",
		},
	}

	for _, tt := range userTests {
		u, err := svc.NewUser(tt.user)
		if err != nil {
			t.Fatal(err)
		}

		if u.Username != tt.user.Username {
			t.Errorf("expected %q, got %q", tt.user.Username, u.Username)
		}

		if u.ID == 0 {
			t.Errorf("expected a user ID, got 0")
		}

		if err := u.ValidatePassword(tt.password); err != nil {
			t.Errorf("expected nil, got %q", err)
		}

		if err := u.ValidatePassword("notthepassword"); err == nil {
			t.Errorf("expected error, got nil")
		}

		if err := svc.SetPassword(u.ID, tt.changePass); err != nil {
			t.Fatalf("expected nil, got %q", err)
		}

		u, err = svc.UserByID(u.ID)
		if err != nil {
			t.Fatalf("failed to retrieve user %q", err)
		}

		if err := u.ValidatePassword(tt.password); err == nil {
			t.Log("original password should fail after SetPassword")
			t.Errorf("expected error, got nil")
		}

		if err := u.ValidatePassword(tt.changePass); err != nil {
			t.Errorf("expected nil, got %q", err)
		}

		userByUsername, err := svc.User(tt.user.Username)
		if err != nil {
			t.Fatalf("failed to retrieve user %q", err)
		}

		if userByUsername.ID != u.ID {
			t.Fatalf("expected %q, got %q", u.ID, userByUsername.ID)
		}
	}

	// test user does not exist
	badID := uint(len(userTests) + 100)
	_, err := svc.UserByID(badID)
	if err != datastore.ErrNotFound {
		t.Errorf("expected %q, got %q", datastore.ErrNotFound, err)
	}

	badUsername := "doesnotexist"
	_, err = svc.User(badUsername)
	if err != datastore.ErrNotFound {
		t.Errorf("expected %q, got %q", datastore.ErrNotFound, err)
	}

}
