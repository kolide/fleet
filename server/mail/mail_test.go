package mail

import (
	"os"
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
)

func TestSMTPNoAuth(t *testing.T) {
	if os.Getenv("MAIL_TEST") == "" {
		t.Skip("Mail testing disabled")
	}
	mail := kolide.Email{
		Subject: "test",
		To:      []string{"bob@foo.com"},
		Config: &kolide.SMTPConfig{
			Configured:         true,
			Disabled:           false,
			AuthenticationType: kolide.AuthTypeNone,
			EnableSSLTLS:       true,
			VerifySSLCerts:     true,
			EnableStartTLS:     true,
			Port:               1025,
			Server:             "localhost",
			SenderAddress:      "kolide@kolide.com",
		},
		Mailer: &kolide.SMTPTestMailer{
			KolideServerURL: "https://localhost:8080",
		},
	}

	mailService := NewService()
	err := mailService.SendEmail(mail)
	assert.Nil(t, err)
}
