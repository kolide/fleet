package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/mail"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

// mailError is set when an error performing mail operations
type mailError struct {
	message string
}

func (e mailError) Error() string {
	return fmt.Sprintf("a mail error occurred: %s", e.message)
}

func (e mailError) MailError() []map[string]string {
	return []map[string]string{
		map[string]string{
			"name":   "base",
			"reason": e.message,
		},
	}
}

func (svc service) NewAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	config, err := svc.ds.AppConfig()
	if err != nil {
		return nil, err
	}
	fromPayload := appConfigFromAppConfigPayload(p, *config)
	if fromPayload.EnrollSecret == "" {
		// generate a random string if the user hasn't set one in the form.
		rand, err := kolide.RandomText(24)
		if err != nil {
			return nil, errors.Wrap(err, "generate enroll secret string")
		}
		fromPayload.EnrollSecret = rand
	}
	newConfig, err := svc.ds.NewAppConfig(fromPayload)
	if err != nil {
		return nil, err
	}
	return newConfig, nil
}

func (svc service) AppConfig(ctx context.Context) (*kolide.AppConfig, error) {
	return svc.ds.AppConfig()
}

func (svc service) SendTestEmail(ctx context.Context, config *kolide.AppConfig) error {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return errNoContext
	}

	testMail := kolide.Email{
		Subject: "Hello from Kolide",
		To:      []string{vc.User.Email},
		Mailer: &kolide.SMTPTestMailer{
			KolideServerURL: config.KolideServerURL,
		},
		Config: config,
	}

	if err := mail.Test(svc.mailService, testMail); err != nil {
		return mailError{message: err.Error()}
	}
	return nil

}

func (svc service) ModifyAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	oldAppConfig, err := svc.AppConfig(ctx)
	if err != nil {
		return nil, err
	}
	config := appConfigFromAppConfigPayload(p, *oldAppConfig)

	if p.SMTPSettings != nil {
		if err = svc.SendTestEmail(ctx, config); err != nil {
			return nil, err
		}
		config.SMTPConfigured = true
	}

	if err := svc.ds.SaveAppConfig(config); err != nil {
		return nil, err
	}
	return config, nil
}

func cleanupURL(url string) string {
	return strings.TrimRight(strings.Trim(url, " \t\n"), "/")
}

func appConfigFromAppConfigPayload(p kolide.AppConfigPayload, config kolide.AppConfig) *kolide.AppConfig {
	if p.OrgInfo != nil && p.OrgInfo.OrgLogoURL != nil {
		config.OrgLogoURL = *p.OrgInfo.OrgLogoURL
	}
	if p.OrgInfo != nil && p.OrgInfo.OrgName != nil {
		config.OrgName = *p.OrgInfo.OrgName
	}
	if p.ServerSettings != nil && p.ServerSettings.KolideServerURL != nil {
		config.KolideServerURL = cleanupURL(*p.ServerSettings.KolideServerURL)
	}
	if p.ServerSettings != nil && p.ServerSettings.EnrollSecret != nil {
		config.EnrollSecret = *p.ServerSettings.EnrollSecret
	}

	if p.SSOSettings != nil {
		if p.SSOSettings.EnableSSO != nil {
			config.EnableSSO = *p.SSOSettings.EnableSSO
		}
		if p.SSOSettings.EntityID != nil {
			config.EntityID = *p.SSOSettings.EntityID
		}
		if p.SSOSettings.IDPImageURL != nil {
			config.IDPImageURL = *p.SSOSettings.IDPImageURL
		}
		if p.SSOSettings.IDPName != nil {
			config.IDPName = *p.SSOSettings.IDPName
		}
		if p.SSOSettings.IssuerURI != nil {
			config.IssuerURI = *p.SSOSettings.IssuerURI
		}
		if p.SSOSettings.Metadata != nil {
			config.Metadata = *p.SSOSettings.Metadata
		}
		if p.SSOSettings.MetadataURL != nil {
			config.MetadataURL = *p.SSOSettings.MetadataURL
		}
	}

	populateSMTP := func(p *kolide.SMTPSettingsPayload) {
		if p.SMTPAuthenticationMethod != nil {
			switch *p.SMTPAuthenticationMethod {
			case kolide.AuthMethodNameCramMD5:
				config.SMTPAuthenticationMethod = kolide.AuthMethodCramMD5
			case kolide.AuthMethodNamePlain:
				config.SMTPAuthenticationMethod = kolide.AuthMethodPlain
			default:
				panic("unknown SMTP AuthMethod: " + *p.SMTPAuthenticationMethod)
			}
		}
		if p.SMTPAuthenticationType != nil {
			switch *p.SMTPAuthenticationType {
			case kolide.AuthTypeNameUserNamePassword:
				config.SMTPAuthenticationType = kolide.AuthTypeUserNamePassword
			case kolide.AuthTypeNameNone:
				config.SMTPAuthenticationType = kolide.AuthTypeNone
			default:
				panic("unknown SMTP AuthType: " + *p.SMTPAuthenticationType)
			}
		}

		if p.SMTPDomain != nil {
			config.SMTPDomain = *p.SMTPDomain
		}

		if p.SMTPEnableStartTLS != nil {
			config.SMTPEnableStartTLS = *p.SMTPEnableStartTLS
		}

		if p.SMTPEnableTLS != nil {
			config.SMTPEnableTLS = *p.SMTPEnableTLS
		}

		if p.SMTPPassword != nil {
			config.SMTPPassword = *p.SMTPPassword
		}

		if p.SMTPPort != nil {
			config.SMTPPort = *p.SMTPPort
		}

		if p.SMTPSenderAddress != nil {
			config.SMTPSenderAddress = *p.SMTPSenderAddress
		}

		if p.SMTPServer != nil {
			config.SMTPServer = *p.SMTPServer
		}

		if p.SMTPUserName != nil {
			config.SMTPUserName = *p.SMTPUserName
		}

		if p.SMTPVerifySSLCerts != nil {
			config.SMTPVerifySSLCerts = *p.SMTPVerifySSLCerts
		}
	}

	if p.SMTPSettings != nil {
		populateSMTP(p.SMTPSettings)
	}
	return &config
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeModifyAppConfigRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var payload kolide.AppConfigPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return appConfigRequest{Payload: payload}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type appConfigRequest struct {
	Payload kolide.AppConfigPayload
}

type appConfigResponse struct {
	OrgInfo        *kolide.OrgInfo             `json:"org_info,omitemtpy"`
	ServerSettings *kolide.ServerSettings      `json:"server_settings,omitempty"`
	SMTPSettings   *kolide.SMTPSettingsPayload `json:"smtp_settings,omitempty"`
	SSOSettings    *kolide.SSOSettingsPayload  `json:"sso_settings,omitempty"`
	Err            error                       `json:"error,omitempty"`
}

func (r appConfigResponse) error() error { return r.Err }

func makeGetAppConfigEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vc, ok := viewer.FromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("could not fetch user")
		}
		config, err := svc.AppConfig(ctx)
		if err != nil {
			return nil, err
		}
		var smtpSettings *kolide.SMTPSettingsPayload
		var ssoSettings *kolide.SSOSettingsPayload
		// only admin can see smtp settings
		if vc.IsAdmin() {
			smtpSettings = smtpSettingsFromAppConfig(config)
			if smtpSettings.SMTPPassword != nil {
				*smtpSettings.SMTPPassword = "********"
			}
			ssoSettings = &kolide.SSOSettingsPayload{
				EntityID:    &config.EntityID,
				IssuerURI:   &config.IssuerURI,
				IDPImageURL: &config.IDPImageURL,
				Metadata:    &config.Metadata,
				MetadataURL: &config.MetadataURL,
				IDPName:     &config.IDPName,
				EnableSSO:   &config.EnableSSO,
			}
		}
		response := appConfigResponse{
			OrgInfo: &kolide.OrgInfo{
				OrgName:    &config.OrgName,
				OrgLogoURL: &config.OrgLogoURL,
			},
			ServerSettings: &kolide.ServerSettings{
				KolideServerURL: &config.KolideServerURL,
				EnrollSecret:    &config.EnrollSecret,
			},
			SMTPSettings: smtpSettings,
			SSOSettings:  ssoSettings,
		}
		return response, nil
	}
}

func makeModifyAppConfigEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(appConfigRequest)
		config, err := svc.ModifyAppConfig(ctx, req.Payload)
		if err != nil {
			return appConfigResponse{Err: err}, nil
		}
		response := appConfigResponse{
			OrgInfo: &kolide.OrgInfo{
				OrgName:    &config.OrgName,
				OrgLogoURL: &config.OrgLogoURL,
			},
			ServerSettings: &kolide.ServerSettings{
				KolideServerURL: &config.KolideServerURL,
				EnrollSecret:    &config.EnrollSecret,
			},
			SMTPSettings: smtpSettingsFromAppConfig(config),
			SSOSettings: &kolide.SSOSettingsPayload{
				EntityID:    &config.EntityID,
				IssuerURI:   &config.IssuerURI,
				IDPImageURL: &config.IDPImageURL,
				Metadata:    &config.Metadata,
				MetadataURL: &config.MetadataURL,
				IDPName:     &config.IDPName,
				EnableSSO:   &config.EnableSSO,
			},
		}
		if response.SMTPSettings.SMTPPassword != nil {
			*response.SMTPSettings.SMTPPassword = "********"
		}
		return response, nil
	}
}

func smtpSettingsFromAppConfig(config *kolide.AppConfig) *kolide.SMTPSettingsPayload {
	authType := config.SMTPAuthenticationType.String()
	authMethod := config.SMTPAuthenticationMethod.String()
	return &kolide.SMTPSettingsPayload{
		SMTPConfigured:           &config.SMTPConfigured,
		SMTPSenderAddress:        &config.SMTPSenderAddress,
		SMTPServer:               &config.SMTPServer,
		SMTPPort:                 &config.SMTPPort,
		SMTPAuthenticationType:   &authType,
		SMTPUserName:             &config.SMTPUserName,
		SMTPPassword:             &config.SMTPPassword,
		SMTPEnableTLS:            &config.SMTPEnableTLS,
		SMTPAuthenticationMethod: &authMethod,
		SMTPDomain:               &config.SMTPDomain,
		SMTPVerifySSLCerts:       &config.SMTPVerifySSLCerts,
		SMTPEnableStartTLS:       &config.SMTPEnableStartTLS,
	}
}

////////////////////////////////////////////////////////////////////////////////
// Metrics
////////////////////////////////////////////////////////////////////////////////

func (mw metricsMiddleware) NewAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	var (
		info *kolide.AppConfig
		err  error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "NewOrgInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	info, err = mw.Service.NewAppConfig(ctx, p)
	return info, err
}

func (mw metricsMiddleware) AppConfig(ctx context.Context) (*kolide.AppConfig, error) {
	var (
		info *kolide.AppConfig
		err  error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "OrgInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	info, err = mw.Service.AppConfig(ctx)
	return info, err
}

func (mw metricsMiddleware) ModifyAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	var (
		info *kolide.AppConfig
		err  error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "ModifyOrgInfo", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	info, err = mw.Service.ModifyAppConfig(ctx, p)
	return info, err
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) NewAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	var (
		info *kolide.AppConfig
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "NewAppConfig",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	info, err = mw.Service.NewAppConfig(ctx, p)
	return info, err
}

func (mw loggingMiddleware) AppConfig(ctx context.Context) (*kolide.AppConfig, error) {
	var (
		info *kolide.AppConfig
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "AppConfig",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	info, err = mw.Service.AppConfig(ctx)
	return info, err
}

func (mw loggingMiddleware) ModifyAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	var (
		info *kolide.AppConfig
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ModifyAppConfig",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	info, err = mw.Service.ModifyAppConfig(ctx, p)
	return info, err
}

////////////////////////////////////////////////////////////////////////////////
// Validation
////////////////////////////////////////////////////////////////////////////////

func (mw validationMiddleware) ModifyAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	existing, err := mw.ds.AppConfig()
	if err != nil {
		return nil, errors.Wrap(err, "fetching existing app config in validation")
	}
	invalid := &invalidArgumentError{}
	validateSSOSettings(p, existing, invalid)
	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.ModifyAppConfig(ctx, p)
}

func isSet(val *string) bool {
	if val != nil {
		return len(*val) > 0
	}
	return false
}

func validateSSOSettings(p kolide.AppConfigPayload, existing *kolide.AppConfig, invalid *invalidArgumentError) {
	if p.SSOSettings != nil && p.SSOSettings.EnableSSO != nil {
		if *p.SSOSettings.EnableSSO {
			if !isSet(p.SSOSettings.Metadata) && !isSet(p.SSOSettings.MetadataURL) {
				if existing.Metadata == "" && existing.MetadataURL == "" {
					invalid.Append("metadata", "either metadata or metadata_url must be defined")
				}
			}
			if isSet(p.SSOSettings.Metadata) && isSet(p.SSOSettings.MetadataURL) {
				invalid.Append("metadata", "both metadata and metadata_url are defined, only one is allowed")
			}
			if !isSet(p.SSOSettings.EntityID) {
				if existing.EntityID == "" {
					invalid.Append("entity_id", "required")
				}
			} else {
				if len(*p.SSOSettings.EntityID) < 5 {
					invalid.Append("entity_id", "must be 5 or more characters")
				}
			}
			if !isSet(p.SSOSettings.IDPName) {
				if existing.IDPName == "" {
					invalid.Append("idp_name", "required")
				}
			} else {
				if len(*p.SSOSettings.IDPName) < 5 {
					invalid.Append("idp_name", "must be 5 or more characters")
				}
			}
		}
	}
}
