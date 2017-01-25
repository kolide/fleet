package kolide

import (
	"bytes"
	"html/template"

	"golang.org/x/net/context"
)

type EmailChangesStore interface {
	// PendingEmailChange creates a record with a pending email change for a user identified
	// by uid. The change record is keyed by a unique token. The token is emailed to the user
	// with a link that they can use to confirm the change.
	PendingEmailChange(uid uint, newEmail, token string) error
	// CommitEmailChange will confirm new email address identified by token is valid.
	// The new email will be written to user record.
	CommitEmailChange(token string) (string, error)
}

type EmailChangeService interface {
	// CommitEmailChange is used to confirm new email address and if confirmed,
	// write the new email address to user.
	CommitEmailChange(ctx context.Context, token string) (string, error)
}

type ChangeEmailMailer struct {
	KolideServerURL template.URL
	Token           string
}

func (cem *ChangeEmailMailer) Message() ([]byte, error) {
	t, err := getTemplate("server/mail/templates/change_email_confirmation.html")
	if err != nil {
		return nil, err
	}
	var msg bytes.Buffer
	err = t.Execute(&msg, cem)
	if err != nil {
		return nil, err
	}
	return msg.Bytes(), nil
}
