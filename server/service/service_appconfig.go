package service

import (
	"fmt"
	"net/smtp"

	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) NewAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {

	newConfig, err := svc.ds.NewAppConfig(fromPayload(p, kolide.AppConfig{}))
	if err != nil {
		return nil, err
	}
	return newConfig, nil
}

func (svc service) AppConfig(ctx context.Context) (*kolide.AppConfig, error) {
	return svc.ds.AppConfig()
}

func (svc service) ModifyAppConfig(ctx context.Context, r kolide.ModifyAppConfigRequest) (*kolide.ModifyAppConfigPayload, error) {
	// Test SMTP but don't save anything back to the db
	if r.TestSMTP {

	}
	config, err := svc.ds.AppConfig()
	if err != nil {
		return nil, err
	}

	if err := svc.ds.SaveAppConfig(&r.AppConfig); err != nil {
		return nil, err
	}

	response := &kolide.ModifyAppConfigPayload{
		AppConfig: *config,
		SMTPStatus: kolide.SMTPResponse{
			Details: map[string]string{},
			Success: true,
		},
	}
	return response, nil

}

func appendSMTPError(errName, errDescription string, resp *kolide.SMTPConfigResponse) {
	smtpError := kolide.SMTPError{
		Name:        errName,
		Description: errDescription,
	}
	resp.Errors = append(resp.Errors, smtpError)
	resp.Success = false
}

const (
	missingSmtpArg = "missing argument"
	smtpError      = "smtp error"
)

// testsSMTPConfiguration userEmail is the email of the current user
func testSMTPConfiguration(userEmail string, config *kolide.AppConfig) *kolide.SMTPConfigResponse {
	response := &kolide.SMTPConfigResponse{
		Success: false,
	}

	if config.SenderAddress == "" {
		appendSMTPError(missingSmtpArg, "missing smtp sender address", response)
	}
	if config.Server == "" {
		appendSMTPError(missingSmtpArg, "missing smtp server host name", response)
	}
	if config.AuthenticationType == kolide.AuthTypeUserNamePassword {
		if config.UserName == "" {
			appendSMTPError(missingSmtpArg, "missing smtp user name", response)
		}
		if config.Password == "" {
			appendSMTPError(missingSmtpArg, "missing smtp password", response)
		}
	}
	// If we are missing values we need to connect, exit
	if len(response.Errors) > 0 {
		return response
	}

	recipient := []string{
		userEmail,
	}

	smtpHost := fmt.Sprintf("%s:%d", config.Server, config.Port)

	if config.AuthenticationType == kolide.AuthTypeUserNamePassword {

	} else {

		err := smtp.SendMail(smtpHost, nil, config.SenderAddress, recipient, []byte("Hello from Kolide!\r\n"))
		if err != nil {
			appendSMTPError(smtpError, err.Error(), response)
			return response
		}
		response.Success = true
	}

	return response
}

func fromPayload(p kolide.AppConfigPayload, config kolide.AppConfig) *kolide.AppConfig {
	if p.OrgInfo != nil && p.OrgInfo.OrgLogoURL != nil {
		config.OrgLogoURL = *p.OrgInfo.OrgLogoURL
	}
	if p.OrgInfo != nil && p.OrgInfo.OrgName != nil {
		config.OrgName = *p.OrgInfo.OrgName
	}
	if p.ServerSettings != nil && p.ServerSettings.KolideServerURL != nil {
		config.KolideServerURL = *p.ServerSettings.KolideServerURL
	}
	return &config
}
