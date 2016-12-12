// Package mail provides implementations of the Kolide MailService
package mail

import (
	"fmt"
	"net/smtp"

	"github.com/kolide/kolide-ose/server/kolide"
)

func NewService() kolide.MailService {
	return simple{}
}

type simple struct {
}

func (m simple) SendEmail(e kolide.Email) error {
	if !e.Config.Disabled {
		if e.Config.Configured {

			body, err := e.Mailer.Message()
			if err != nil {
				return err
			}

			mime := `MIME-version: 1.0;
			Content-Type: text/html; charset="UTF-8";`
			subject := "Subject: " + e.Subject + "\r\n"
			msg := []byte(subject + mime + "\r\n\r\n" + string(body) + "\r\n")
			smtpHost := fmt.Sprintf("%s:%d", e.Config.Server, e.Config.Port)

			if e.Config.AuthenticationType != kolide.AuthTypeUserNamePassword {
				err = smtp.SendMail(smtpHost, nil, e.Config.SenderAddress, e.To, msg)
				if err != nil {
					return err
				}

			} else {
				var auth smtp.Auth
				switch e.Config.AuthenticationMethod {
				case kolide.AuthMethodCramMD5:
					auth = smtp.CRAMMD5Auth(e.Config.UserName, e.Config.Password)
				case kolide.AuthMethodPlain:
					auth = smtp.PlainAuth("", e.Config.UserName, e.Config.Password, e.Config.Server)
				default:
					return fmt.Errorf("Unknown SMTP auth type '%s'", e.Config.AuthenticationMethod)
				}

				if e.Config.EnableSSLTLS && e.Config.EnableStartTLS {
					err = smtp.SendMail(smtpHost, auth, e.Config.SenderAddress, e.To, msg)
					if err != nil {
						return err
					}
				}
			}

		}
	}

	return nil
}
