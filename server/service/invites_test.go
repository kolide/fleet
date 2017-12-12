package service

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/gorilla/mux"
	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func TestInviteNewUserMock(t *testing.T) {
	svc, mockStore, mailer := setupInviteTest(t)
	ctx := context.Background()

	payload := kolide.InvitePayload{
		Email:     stringPtr("user@acme.co"),
		InvitedBy: &adminUser.ID,
		Admin:     boolPtr(false),
	}

	// happy path
	invite, err := svc.InviteNewUser(ctx, payload)
	require.Nil(t, err)
	assert.Equal(t, invite.ID, validInvite.ID)
	assert.True(t, mockStore.NewInviteFuncInvoked)
	assert.True(t, mockStore.AppConfigFuncInvoked)
	assert.True(t, mailer.Invoked)

	mockStore.UserByEmailFunc = mock.UserByEmailWithUser(new(kolide.User))
	_, err = svc.InviteNewUser(ctx, payload)
	require.NotNil(t, err, "should err if the user we're inviting already exists")
}

func TestVerifyInvite(t *testing.T) {
	ms := new(mock.Store)
	svc := service{
		ds:     ms,
		config: config.TestConfig(),
		clock:  clock.NewMockClock(),
	}
	ctx := context.Background()

	ms.InviteByTokenFunc = mock.ReturnFakeInviteByToken(expiredInvite)
	wantErr := &invalidArgumentError{{name: "invite_token", reason: "Invite token has expired."}}
	_, err := svc.VerifyInvite(ctx, expiredInvite.Token)
	assert.Equal(t, err, wantErr)

	wantErr = &invalidArgumentError{{name: "invite_token",
		reason: "Invite Token does not match Email Address."}}

	_, err = svc.VerifyInvite(ctx, "bad_token")
	assert.Equal(t, err, wantErr)
}

func TestDeleteInvite(t *testing.T) {
	ms := new(mock.Store)
	svc := service{ds: ms}

	ms.DeleteInviteFunc = func(uint) error { return nil }
	err := svc.DeleteInvite(context.Background(), 1)
	require.Nil(t, err)
	assert.True(t, ms.DeleteInviteFuncInvoked)
}

func TestListInvites(t *testing.T) {
	ms := new(mock.Store)
	svc := service{ds: ms}

	ms.ListInvitesFunc = func(kolide.ListOptions) ([]*kolide.Invite, error) {
		return nil, nil
	}
	_, err := svc.ListInvites(context.Background(), kolide.ListOptions{})
	require.Nil(t, err)
	assert.True(t, ms.ListInvitesFuncInvoked)
}

func setupInviteTest(t *testing.T) (kolide.Service, *mock.Store, *mockMailService) {

	ms := new(mock.Store)
	ms.UserByEmailFunc = mock.UserWithEmailNotFound()
	ms.UserByIDFunc = mock.UserWithID(adminUser)
	ms.NewInviteFunc = mock.ReturnNewInivite(validInvite)
	ms.AppConfigFunc = mock.ReturnFakeAppConfig(&kolide.AppConfig{
		KolideServerURL: "https://acme.co",
	})
	mailer := &mockMailService{SendEmailFn: func(e kolide.Email) error { return nil }}
	svc := validationMiddleware{service{
		ds:          ms,
		config:      config.TestConfig(),
		mailService: mailer,
		clock:       clock.NewMockClock(),
	}, ms, nil}
	return svc, ms, mailer
}

var adminUser = &kolide.User{
	ID:       1,
	Email:    "admin@acme.co",
	Username: "admin",
	Name:     "Administrator",
}

var existingUser = &kolide.User{
	ID:       2,
	Email:    "user@acme.co",
	Username: "user",
}

var validInvite = &kolide.Invite{
	ID:    1,
	Token: "abcd",
}

var expiredInvite = &kolide.Invite{
	ID:    1,
	Token: "abcd",
	UpdateCreateTimestamps: kolide.UpdateCreateTimestamps{
		CreateTimestamp: kolide.CreateTimestamp{
			CreatedAt: time.Now().AddDate(-1, 0, 0),
		},
	},
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func TestDecodeCreateInviteRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/invites", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeCreateInviteRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(createInviteRequest)
		assert.Equal(t, uint(1), *params.payload.InvitedBy)
	}).Methods("POST")

	t.Run("lowercase email", func(t *testing.T) {
		var body bytes.Buffer
		body.Write([]byte(`{
        "name": "foo",
        "email": "foo@kolide.co",
        "invited_by": 1
    }`))

		router.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest("POST", "/api/v1/kolide/invites", &body),
		)
	})

	t.Run("uppercase email", func(t *testing.T) {
		// email string should be lowerased after decode.
		var body bytes.Buffer
		body.Write([]byte(`{
        "name": "foo",
        "email": "Foo@Kolide.co",
        "invited_by": 1
    }`))

		router.ServeHTTP(
			httptest.NewRecorder(),
			httptest.NewRequest("POST", "/api/v1/kolide/invites", &body),
		)
	})

}

func TestDecodeVerifyInviteRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/invites/{token}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeCreateInviteRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(verifyInviteRequest)
		assert.Equal(t, "test_token", params.Token)
	}).Methods("GET")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/api/v1/kolide/tokens/test_token", nil),
	)

}
