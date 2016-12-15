package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/contexts/viewer"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

type modifyAppConfigResponse struct {
	Response *kolide.AppConfig `json:"app_config,omitempty"`
	Err      error             `json:"error,omitempty"`
}

func (m modifyAppConfigResponse) error() error { return m.Err }

type getAppConfigResponse struct {
	OrgInfo        *kolide.OrgInfo        `json:"org_info,omitemtpy"`
	ServerSettings *kolide.ServerSettings `json:"server_settings,omitempty"`
	AppConfig      *kolide.AppConfig      `json:"app_config,omitempty"`
	Err            error                  `json:"error,omitempty"`
}

func (r getAppConfigResponse) error() error { return r.Err }

func makeGetAppConfigEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		config, err := svc.AppConfig(ctx)
		if err != nil {
			return getAppConfigResponse{Err: err}, nil
		}
		v, _ := viewer.FromContext(ctx)

		if !v.IsAdmin() {
			// make a copy of config so we don't munge inmem
			copyConfig := *config
			copyConfig.SMTPConfig = nil
			return getAppConfigResponse{AppConfig: &copyConfig}, nil
		}

		response := getAppConfigResponse{
			// TODO: OrgInfo and ServerSettings should be removed once front end is updated to
			// get OrgName and OrgLogoURL from AppConfig see Issue #649
			OrgInfo: &kolide.OrgInfo{
				OrgName:    &config.OrgName,
				OrgLogoURL: &config.OrgLogoURL,
			},
			ServerSettings: &kolide.ServerSettings{
				KolideServerURL: &config.KolideServerURL,
			},
			AppConfig: config,
		}

		return response, nil
	}
}

func makeModifyAppConfigRequest(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(kolide.ModifyAppConfigRequest)
		payload, err := svc.ModifyAppConfig(ctx, req)
		if err != nil {
			return modifyAppConfigResponse{Err: err}, nil
		}
		return modifyAppConfigResponse{Response: payload}, nil
	}
}
