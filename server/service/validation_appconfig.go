package service

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (mw validationMiddleware) ModifyAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	invalid := &invalidArgumentError{}

	if p.ServerSettings.KolideServerURL == nil || *p.ServerSettings.KolideServerURL == "" {
		invalid.Append("kolide_server_url", "missing")
	}
	if p.ServerSettings.KolideServerURL != nil && *p.ServerSettings.KolideServerURL != "" {
		if err := validateKolideServerURL(*p.ServerSettings.KolideServerURL); err != nil {
			invalid.Append("kolide_server_url", err.Error())
		}
	}
	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.ModifyAppConfig(ctx, p)
}
