package server

import (
	"errors"
	"testing"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/kolide/kolide-ose/config"
	"github.com/kolide/kolide-ose/datastore"
	"github.com/kolide/kolide-ose/kolide"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
)

func TestRequestPasswordReset(t *testing.T) {
	ds, _ := datastore.New("inmem", "")
	createTestUsers(t, ds)
	admin1, _ := ds.User("admin1")
	user1, _ := ds.User("user1")
	// user2, _ := ds.User("user2")
	var defaultEmailFn = func(e kolide.Email) error {
		return nil
	}
	var errEmailFn = func(e kolide.Email) error {
		return errors.New("test err")
	}
	svc := service{ds: ds,
		smtpTokenKeySize: 10,
	}

	var requestPasswordResetTests = []struct {
		email   string
		emailFn func(e kolide.Email) error
		wantErr interface{}
		user    *kolide.User
		vc      *viewerContext
	}{
		{
			email:   admin1.Email,
			emailFn: defaultEmailFn,
			user:    admin1,
			vc:      &viewerContext{user: admin1},
		},
		{
			email:   admin1.Email,
			emailFn: defaultEmailFn,
			user:    admin1,
			vc:      emptyVC(),
		},
		{
			email:   user1.Email,
			emailFn: defaultEmailFn,
			user:    user1,
			vc:      &viewerContext{user: admin1},
		},
		{
			email:   admin1.Email,
			emailFn: errEmailFn,
			user:    user1,
			vc:      emptyVC(),
			wantErr: "test err",
		},
	}

	for i, tt := range requestPasswordResetTests {
		ctx := context.Background()
		if tt.vc != nil {
			ctx = context.WithValue(ctx, "viewerContext", tt.vc)
		}
		mailer := &mockMailSvc{SendEmailFn: tt.emailFn}
		svc.mailService = mailer
		serr := svc.RequestPasswordReset(ctx, tt.email)
		if err := matchErr(serr, tt.wantErr); err != nil {
			t.Errorf("test id %d failed with %v", i, err)
		}
		if tt.vc.IsAdmin() {
			if have, want := mailer.Invoked, false; have != want {
				t.Errorf("test id %d sentEmail: have: %v, want: %v", i, have, want)
			}
			if have, want := tt.user.AdminForcedPasswordReset, true; have != want {
				t.Errorf("test id %d AdminForcedPasswordReset: have: %v, want: %v", i, have, want)
			}
		} else {
			if have, want := mailer.Invoked, true; have != want {
				t.Errorf("test id %d sentEmail: have: %v, want: %v", i, have, want)
			}

			if serr == nil {
				req, err := ds.FindPassswordResetsByUserID(tt.user.ID)
				if err != nil {
					t.Fatalf("expected password reset request, got err: %q", err)
				}
				if len(req) == 0 {
					t.Errorf("expected at least one password request token, got 0")
				}
			}
		}
	}
}

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
	// admin1, _ := ds.User("admin1")
	user1, _ := ds.User("user1")
	user2, _ := ds.User("user2")
	request := &kolide.PasswordResetRequest{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24),
		UserID:    user1.ID,
		Token:     "abcd",
	}
	_, err := ds.NewPasswordResetRequest(request)
	if err != nil {
		t.Fatal(err)
	}

	var passwordChangeTests = []struct {
		token       string
		newPassword string
		vc          *viewerContext
		wantErr     interface{}
		user        *kolide.User
		// what resource are we editing
		requestID uint
	}{
		{ // all good
			requestID:   user1.ID,
			token:       "abcd",
			vc:          emptyVC(),
			newPassword: "123cat!",
		},
		{ // bad token
			requestID:   user1.ID,
			token:       "dcba",
			vc:          emptyVC(),
			newPassword: "123cat!",
			wantErr:     "resource not found",
		},
		{ // missing token
			requestID:   user1.ID,
			vc:          emptyVC(),
			newPassword: "123cat!",
			wantErr:     "argument invalid or missing: token",
		},
		{
			requestID:   user1.ID,
			vc:          &viewerContext{user: user1},
			newPassword: "123cat!",
		},
		{ // missing password
			requestID: user2.ID,
			vc:        &viewerContext{user: user2},
			wantErr:   "argument invalid or missing: password",
		},
		{ // no such user
			requestID:   999,
			vc:          emptyVC(),
			token:       "abcd",
			wantErr:     "resource not found",
			newPassword: "123cat!",
		},
	}

	for _, tt := range passwordChangeTests {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "viewerContext", tt.vc)
		request := &kolide.PasswordResetRequest{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: time.Now().Add(time.Hour * 24),
			UserID:    tt.ID,
			Token:     tt.token,
		}
		_, err = ds.NewPasswordResetRequest(request)
		assert.Nil(t, err)

		serr := svc.ChangePassword(ctx, tt.requestID, tt.token, tt.newPassword)
		assert.Equal(t, tt.wantErr, serr)
	}
}

type mockMailSvc struct {
	SendEmailFn func(e kolide.Email) error
	Invoked     bool
}

func (svc *mockMailSvc) SendEmail(e kolide.Email) error {
	svc.Invoked = true
	return svc.SendEmailFn(e)
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
