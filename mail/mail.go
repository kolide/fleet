// Package mail provides implementations of the Kolide MailService
package mail

import (
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/kolide/kolide-ose/config"
	"github.com/kolide/kolide-ose/kolide"
)

func NewService(config config.SMTPConfig) kolide.MailService {
	var auth smtp.Auth
	conn := fmt.Sprintf("%s:%s", config.Server, strconv.Itoa(587))
	return simple{Auth: auth, Conn: conn}
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
