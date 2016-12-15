package service

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) NewAppConfig(ctx context.Context, p kolide.AppConfigPayload) (*kolide.AppConfig, error) {
	config, err := svc.ds.AppConfig()
	if err != nil {
		return nil, err
	}
	newConfig, err := svc.ds.NewAppConfig(fromPayload(p, *config))
	if err != nil {
		return nil, err
	}
	return newConfig, nil
}

func (svc service) AppConfig(ctx context.Context) (*kolide.AppConfig, error) {
	return svc.ds.AppConfig()
}

func (svc service) ModifyAppConfig(ctx context.Context, r kolide.ModifyAppConfigRequest) (*kolide.AppConfig, error) {
	// If TestSMTP is true and we got here, we already successfully tested smtp
	// so return without writing to db
	if r.TestSMTP {
		return &r.AppConfig, nil
	}

	if err := svc.ds.SaveAppConfig(&r.AppConfig); err != nil {
		return nil, err
	}

	return &r.AppConfig, nil

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
