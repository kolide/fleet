package service

import (
	"context"
	"testing"

	"github.com/WatchBeam/clock"
	"github.com/kolide/kolide-ose/server/config"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/kolide/kolide-ose/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInviteNewUserMock(t *testing.T) {
	svc, mockStore := setupInviteTest(t)
	ctx := context.Background()

	payload := kolide.InvitePayload{
		Email:     stringPtr("user@acme.co"),
		InvitedBy: &adminUser.ID,
		Admin:     boolPtr(false),
	}

	invite, err := svc.InviteNewUser(ctx, payload)
	require.Nil(t, err)
	assert.Equal(t, invite.ID, validInvite.ID)
	assert.True(t, mockStore.NewInviteFuncInvoked)

	mockStore.UserByEmailFunc = mock.UserByEmailWithUser(existingUser)

}

func setupInviteTest(t *testing.T) (kolide.Service, *mock.Store) {
	ms := new(mock.Store)
	ms.UserByEmailFunc = mock.UserWithEmailNotFound()
	ms.UserByIDFunc = mock.UserWithID(adminUser)
	ms.NewInviteFunc = mock.ReturnNewInivite(validInvite)
	mailer := &mockMailService{SendEmailFn: func(e kolide.Email) error { return nil }}

	svc := validationMiddleware{service{
		ds:          ms,
		config:      config.TestConfig(),
		mailService: mailer,
		clock:       clock.NewMockClock(),
	}}
	return svc, ms
}

var adminUser = &kolide.User{
	ID:       1,
	Email:    "admin@acme.co",
	Username: "admin",
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
