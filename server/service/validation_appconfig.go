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

	if !r.AppConfig.Disabled {
		if r.AppConfig.SenderAddress != "" {
			invalid.Append("smtp_sender_address", "missing required argument")
		}

		if r.AppConfig.AuthenticationType != kolide.AuthTypeUserNamePassword &&
			r.AppConfig.AuthenticationType != kolide.AuthTypeNone {
			invalid.Append("smtp_authentication_type", "invalid value")
		}

		if r.AppConfig.AuthenticationType != kolide.AuthTypeNone {
			if r.AppConfig.AuthenticationMethod != kolide.AuthMethodCramMD5 &&
				r.AppConfig.AuthenticationMethod != kolide.AuthMethodPlain {
				invalid.Append("smtp_authentication_method", "invalid value")
			}
		}

		if r.AppConfig.AuthenticationMethod == kolide.AuthTypeUserNamePassword {
			if r.AppConfig.UserName == "" {
				invalid.Append("smtp_user_name", "missing required argument")
			}
			if r.AppConfig.Password == "" {
				invalid.Append("smtp_password", "missing required argument")
			}
		}

		if r.AppConfig.Server == "" {
			invalid.Append("smtp_server", "missing require argument")
		}

		if !invalid.HasErrors() {
			if !r.AppConfig.Configured || r.TestSMTP {
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
					r.AppConfig.Configured = true
				}

			}
		}

	}

	if invalid.HasErrors() {
		return nil, invalid
	}

	return mw.Service.ModifyAppConfig(ctx, r)

}
