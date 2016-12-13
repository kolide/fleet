// Package mail provides implementations of the Kolide MailService
package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"github.com/kolide/kolide-ose/server/kolide"
)

func NewService() kolide.MailService {
	return &mailService{}
}

func NewDevService() kolide.MailService {
	return &devMailService{}
}

type mailService struct{}
type devMailService struct{}

const (
	PortSSL = 465
	PortTLS = 587
)

func (dm devMailService) SendEmail(e kolide.Email) error {
	if !e.Config.Disabled {
		if e.Config.Configured {

			body, err := e.Mailer.Message()
			if err != nil {
				return err
			}

			mime := `MIME-version: 1.0;` + "\r\n"
			content := `Content-Type: text/html; charset="UTF-8";` + "\r\n"
			subject := "Subject: " + e.Subject + "\r\n"
			msg := subject + mime + content + "\r\n" + string(body) + "\r\n"
			fmt.Printf(msg)
		}
	}
	return nil
}

func (m mailService) SendEmail(e kolide.Email) error {
	if !e.Config.Disabled {
		if e.Config.Configured {

			body, err := e.Mailer.Message()
			if err != nil {
				return err
			}

			mime := `MIME-version: 1.0;` + "\r\n"
			content := `Content-Type: text/html; charset="UTF-8";` + "\r\n"
			subject := "Subject: " + e.Subject + "\r\n"
			msg := []byte(subject + mime + content + "\r\n" + string(body) + "\r\n")
			smtpHost := fmt.Sprintf("%s:%d", e.Config.Server, e.Config.Port)

			var skipVerify bool
			if e.Config.Port == PortTLS {
				skipVerify = !e.Config.EnableStartTLS
			}

			if e.Config.Port == PortSSL {
				skipVerify = !e.Config.EnableSSLTLS
			}

			if !e.Config.EnableSSLTLS && !e.Config.EnableStartTLS {
				skipVerify = true
			}

			var auth smtp.Auth
			if e.Config.AuthenticationType == kolide.AuthTypeUserNamePassword {
				switch e.Config.AuthenticationMethod {
				case kolide.AuthMethodCramMD5:
					auth = smtp.CRAMMD5Auth(e.Config.UserName, e.Config.Password)
					return smtp.SendMail(smtpHost, auth, e.Config.SenderAddress, e.To, msg)
				case kolide.AuthMethodPlain:
					auth = smtp.PlainAuth("", e.Config.UserName, e.Config.Password, e.Config.Server)

				default:
					return fmt.Errorf("Unknown SMTP auth type '%s'", e.Config.AuthenticationMethod)
				}
			} else {
				auth = nil // No Auth
			}

			if skipVerify {
				return m.sendMailWithoutTLS(auth, smtpHost, e, msg)
			}

			return smtp.SendMail(smtpHost, auth, e.Config.SenderAddress, e.To, msg)
		}
	}
	return nil
}

func (m mailService) sendMailWithoutTLS(auth smtp.Auth, smtpHost string, e kolide.Email, msg []byte) error {
	client, err := smtp.Dial(smtpHost)
	if err != nil {
		return err
	}
	defer client.Close()
	if err = client.Hello(""); err != nil {
		return err
	}
	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{
			ServerName:         e.Config.Server,
			InsecureSkipVerify: true,
		}
		if err = client.StartTLS(config); err != nil {
			return err
		}
	}
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}
	if err = client.Mail(e.Config.SenderAddress); err != nil {
		return err
	}
	for _, recip := range e.To {
		if err = client.Rcpt(recip); err != nil {
			return err
		}
	}
	writer, err := client.Data()
	if err != nil {
		return nil
	}
	_, err = writer.Write(msg)
	if err = writer.Close(); err != nil {
		return err
	}
	return client.Quit()
}
