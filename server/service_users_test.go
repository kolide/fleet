package server

import (
	"testing"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/kolide/kolide-ose/config"
	"github.com/kolide/kolide-ose/datastore"
	"github.com/kolide/kolide-ose/kolide"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestCreateUser(t *testing.T) {
	ds, _ := datastore.New("inmem", "")
	svc, _ := NewService(ds, kitlog.NewNopLogger(), config.TestConfig())
	ctx := context.Background()

	var createUserTests = []struct {
		Username           *string
		Password           *string
		Email              *string
		NeedsPasswordReset *bool
		Admin              *bool
		Err                error
	}{
		{
			Username: stringPtr("admin1"),
			Password: stringPtr("foobar"),
			Err:      invalidArgumentError{},
		},
		{
			Username:           stringPtr("admin1"),
			Password:           stringPtr("foobar"),
			Email:              stringPtr("admin1@example.com"),
			NeedsPasswordReset: boolPtr(true),
			Admin:              boolPtr(false),
		},
	}

	for _, tt := range createUserTests {
		payload := kolide.UserPayload{
			Username: tt.Username,
			Password: tt.Password,
			Email:    tt.Email,
			Admin:    tt.Admin,
			AdminForcedPasswordReset: tt.NeedsPasswordReset,
		}
		user, err := svc.NewUser(ctx, payload)
		switch err.(type) {
		case nil:
		case invalidArgumentError:
			continue
		default:
			t.Fatalf("got %q, want %q", err, tt.Err)
		}

		assert.NotZero(t, user.ID)

		err = user.ValidatePassword(*tt.Password)
		assert.Nil(t, err)

		err = user.ValidatePassword("different_password")
		assert.NotNil(t, err)

		assert.Equal(t, user.AdminForcedPasswordReset, *tt.NeedsPasswordReset)
		assert.Equal(t, user.Admin, *tt.Admin)

		// check duplicate creation
		_, err = svc.NewUser(ctx, payload)
		assert.Equal(t, datastore.ErrExists, err)
	}
}

func TestChangeUserPassword(t *testing.T) {
	ds, _ := datastore.New("inmem", "")
	svc, _ := NewService(ds, kitlog.NewNopLogger(), config.TestConfig())
	createTestUsers(t, ds)

	var passwordChangeTests = []struct {
		username    string
		token       string
		newPassword string
		err         error
	}{
		{
			username:    "admin1",
			token:       "abcd",
			newPassword: "123cat!",
		},
	}

	ctx := context.Background()
	vc := &viewerContext{
		user: &kolide.User{
			Username: "admin1",
			Enabled:  true,
			Admin:    true,
			AdminForcedPasswordReset: true,
		},
	}
	ctx = context.WithValue(ctx, "viewerContext", vc)
	for _, tt := range passwordChangeTests {
		token, err := generateRandomText(10)
		if err != nil {
			t.Fatal(err)
		}
		user, err := ds.User(tt.username)
		assert.Nil(t, err)
		request := &kolide.PasswordResetRequest{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Hour * 24),
			UserID:    user.ID,
			Token:     token,
		}
		_, err = ds.NewPasswordResetRequest(request)
		assert.Nil(t, err)

		err = svc.ChangePassword(ctx, user.ID, tt.token, tt.newPassword)
		assert.Nil(t, err)
	}
}

var testUsers = map[string]kolide.UserPayload{
	"admin1": {
		Username: stringPtr("admin1"),
		Password: stringPtr("foobar"),
		Email:    stringPtr("admin1@example.com"),
		Admin:    boolPtr(true),
	},
	"user1": {
		Username: stringPtr("user1"),
		Password: stringPtr("foobar"),
		Email:    stringPtr("user1@example.com"),
	},
	"user2": {
		Username: stringPtr("user2"),
		Password: stringPtr("bazfoo"),
		Email:    stringPtr("user2@example.com"),
	},
}

func stringPtr(s string) *string {
	return &s
}

func boolPtr(b bool) *bool {
	return &b
}
