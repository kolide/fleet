package app

import (
	"fmt"
	"time"

	"github.com/jordan-wright/email"
	"github.com/kolide/kolide-ose/config"
	"github.com/kolide/kolide-ose/errors"
)

type EmailType int

const (
	PasswordResetEmail EmailType = iota
)

type PasswordResetRequestEmailParameters struct {
	Name  string
	Token string
}

const (
	NoReplyEmailAddress = "no-reply@kolide.co"
)

type SMTPConnectionPool interface {
	Send(e *email.Email, timeout time.Duration) error
	Close()
}

func SendEmail(pool SMTPConnectionPool, to, subject string, html, text []byte) *errors.KolideError {
	e := email.Email{
		From:    fmt.Sprintf("Kolide <%s>", NoReplyEmailAddress),
		To:      []string{to},
		Subject: subject,
		HTML:    html,
		Text:    text,
	}

	err := pool.Send(&e, time.Second*10)
	if err != nil {
		return errors.New("Mail error", "Error attempting to send email on the SMTP pool")
	}

	return nil
}

func GetEmailBody(t EmailType, params interface{}) (html []byte, text []byte, err *errors.KolideError) {
	switch t {
	case PasswordResetEmail:
		resetParams, ok := params.(*PasswordResetRequestEmailParameters)
		if !ok {
			err = errors.New("Couldn't get email body", "Parameters were of incorrect type")
			return
		}

		html = []byte(fmt.Sprintf(
			"Hi %s! <a href=\"%s/password/reset?token=%s\">Reset your password!</a>",
			resetParams.Name,
			config.App.WebAddress,
			resetParams.Token,
		))
		text = []byte(fmt.Sprintf(
			"Hi %s! Reset your password: %s/password/reset?token=%s",
			resetParams.Name,
			config.App.WebAddress,
			resetParams.Token,
		))
	default:
		err = errors.New(
			"Couldn't get email body",
			fmt.Sprintf("Email type unknown: %d", t),
		)
	}
	return
}

func GetEmailSubject(t EmailType) (string, *errors.KolideError) {
	switch t {
	case PasswordResetEmail:
		return "Your Kolide Password Reset Request", nil
	default:
		return "", errors.New(
			"Couldn't get email subject",
			fmt.Sprintf("Email type unknown: %d", t),
		)
	}
}
