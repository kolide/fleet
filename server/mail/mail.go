// Package mail provides implementations of the Kolide MailService
package mail

import (
	"net"
	"net/smtp"

	"github.com/kolide/kolide-ose/server/config"
	"github.com/kolide/kolide-ose/server/kolide"
)

func NewService(config config.SMTPConfig) kolide.MailService {
	host, _, _ := net.SplitHostPort(config.Address)
	auth := smtp.PlainAuth("", config.Username, config.Password, host)
	return simple{Auth: auth, Conn: config.Address}
}

type simple struct {
	Auth smtp.Auth
	// Conn includes the email server and port
	Conn string
}

func (m simple) SendEmail(e kolide.Email) error {
	body, err := e.Msg.Message()
	if err != nil {
		return err
	}
	err = smtp.SendMail(m.Conn, m.Auth, e.From, e.To, body)
	if err != nil {
		return err
	}
	return nil
}
