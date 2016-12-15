package service

import (
	"github.com/kolide/kolide-ose/server/contexts/viewer"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (mw validationMiddleware) ModifyAppConfig(ctx context.Context, r kolide.ModifyAppConfigRequest) (*kolide.AppConfig, error) {
	invalid := &invalidArgumentError{}
	if err := validateKolideServerURL(r.AppConfig.KolideServerURL); err != nil {
		invalid.Append("kolide_server_url", err.Error())
	}

	if !r.AppConfig.SMTPDisabled {
		if r.AppConfig.SMTPSenderAddress != "" {
			invalid.Append("smtp_sender_address", "missing required argument")
		}

		if r.AppConfig.SMTPAuthenticationType != kolide.AuthTypeUserNamePassword &&
			r.AppConfig.SMTPAuthenticationType != kolide.AuthTypeNone {
			invalid.Append("smtp_authentication_type", "invalid value")
		}

		if r.AppConfig.SMTPAuthenticationType != kolide.AuthTypeNone {
			if r.AppConfig.SMTPAuthenticationMethod != kolide.AuthMethodCramMD5 &&
				r.AppConfig.SMTPAuthenticationMethod != kolide.AuthMethodPlain {
				invalid.Append("smtp_authentication_method", "invalid value")
			}
		}

		if r.AppConfig.SMTPAuthenticationMethod == kolide.AuthTypeUserNamePassword {
			if r.AppConfig.SMTPUserName == "" {
				invalid.Append("smtp_user_name", "missing required argument")
			}
			if r.AppConfig.SMTPPassword == "" {
				invalid.Append("smtp_password", "missing required argument")
			}
		}

		if r.AppConfig.SMTPServer == "" {
			invalid.Append("smtp_server", "missing require argument")
		}

		if !invalid.HasErrors() {
			if !r.AppConfig.SMTPConfigured || r.TestSMTP {
				v, ok := viewer.FromContext(ctx)
				if !ok {
					invalid.Append("user", "missing user")
					return nil, invalid
				}

				mail := kolide.Email{
					Subject: "Kolide",
					To:      []string{v.User.Email},
					Mailer: &kolide.SMTPTestMailer{
						KolideServerURL: r.AppConfig.KolideServerURL,
					},
				}

				if err := mw.Service.SendEmail(mail); err != nil {
					invalid.Append("smtp connection", err.Error())
				} else {
					r.AppConfig.SMTPConfigured = true
				}

			}
		}

	}

	if invalid.HasErrors() {
		return nil, invalid
	}

	return mw.Service.ModifyAppConfig(ctx, r)

}
