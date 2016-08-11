package app

import (
	"fmt"

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

func GetEmailBody(t EmailType, params interface{}) (html []byte, text []byte, err error) {
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
		err = errors.New("Couldn't get email body", "Email type unknown")
	}
	return
}

func GetEmailSubject(t EmailType) (string, error) {
	switch t {
	case PasswordResetEmail:
		return "Your Kolide Password Reset Request", nil
	default:
		return "", errors.New("Couldn't get email subject", "Email type unknown")
	}
}
