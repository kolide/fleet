package mail

import (
	"os"
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/suite"
)

type mockMailer struct{}

func (m *mockMailer) SendEmail(e kolide.Email) error {
	return nil
}

type EmailTestSuite struct {
	suite.Suite
	mailService kolide.MailService
}

func (s *EmailTestSuite) SetupTest() {

	if os.Getenv("MAIL_TEST") == "" {
		s.mailService = &mockMailer{}
	} else {
		s.mailService = NewService()
	}

}

func (s *EmailTestSuite) TestSMTPPlainAuth() {
	mail := kolide.Email{
		Subject: "smtp plain auth",
		To:      []string{"john@kolide.co"},
		Config: &kolide.AppConfig{
			SMTPConfigured:           true,
			SMTPDisabled:             false,
			SMTPAuthenticationType:   kolide.AuthTypeUserNamePassword,
			SMTPAuthenticationMethod: kolide.AuthMethodPlain,
			SMTPUserName:             "bob",
			SMTPPassword:             "secret",
			SMTPEnableTLS:            true,
			SMTPVerifySSLCerts:       true,
			SMTPEnableStartTLS:       true,
			SMTPPort:                 1025,
			SMTPServer:               "localhost",
			SMTPSenderAddress:        "kolide@kolide.com",
		},
		Mailer: &kolide.SMTPTestMailer{
			KolideServerURL: "https://localhost:8080",
		},
	}

	err := s.mailService.SendEmail(mail)
	s.Nil(err)
}

// TODO: MailHog doesn't support cram md5
// func (s *EmailTestSuite) TestSMTPCramMD5Auth() {
// 	mail := kolide.Email{
// 		Subject: "cram md5 auth",
// 		To:      []string{"john@kolide.co"},
// 		Config: &kolide.SMTPConfig{
// 			Configured:           true,
// 			Disabled:             false,
// 			AuthenticationType:   kolide.AuthTypeUserNamePassword,
// 			AuthenticationMethod: kolide.AuthMethodCramMD5,
// 			UserName:             "bob",
// 			Password:             "secret",
// 			EnableSSLTLS:         true,
// 			VerifySSLCerts:       true,
// 			EnableStartTLS:       true,
// 			Port:                 1025,
// 			Server:               "localhost",
// 			SenderAddress:        "kolide@kolide.com",
// 		},
// 		Mailer: &kolide.SMTPTestMailer{
// 			KolideServerURL: "https://localhost:8080",
// 		},
// 	}
//
// 	err := s.mailService.SendEmail(mail)
// 	s.Nil(err)
// }

func (s *EmailTestSuite) TestSMTPSkipVerify() {
	mail := kolide.Email{
		Subject: "skip verify",
		To:      []string{"john@kolide.co"},
		Config: &kolide.AppConfig{
			SMTPConfigured:           true,
			SMTPDisabled:             false,
			SMTPAuthenticationType:   kolide.AuthTypeUserNamePassword,
			SMTPAuthenticationMethod: kolide.AuthMethodPlain,
			SMTPUserName:             "bob",
			SMTPPassword:             "secret",
			SMTPEnableTLS:            true,
			SMTPVerifySSLCerts:       false,
			SMTPEnableStartTLS:       true,
			SMTPPort:                 1025,
			SMTPServer:               "localhost",
			SMTPSenderAddress:        "kolide@kolide.com",
		},
		Mailer: &kolide.SMTPTestMailer{
			KolideServerURL: "https://localhost:8080",
		},
	}

	err := s.mailService.SendEmail(mail)
	s.Nil(err)
}

func (s *EmailTestSuite) TestSMTPNoAuth() {
	mail := kolide.Email{
		Subject: "no auth",
		To:      []string{"bob@foo.com"},
		Config: &kolide.AppConfig{
			SMTPConfigured:         true,
			SMTPDisabled:           false,
			SMTPAuthenticationType: kolide.AuthTypeNone,
			SMTPEnableTLS:          true,
			SMTPVerifySSLCerts:     true,
			SMTPPort:               1025,
			SMTPServer:             "localhost",
			SMTPSenderAddress:      "kolide@kolide.com",
		},
		Mailer: &kolide.SMTPTestMailer{
			KolideServerURL: "https://localhost:8080",
		},
	}

	err := s.mailService.SendEmail(mail)
	s.Nil(err)
}

func TestEmailSending(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}
