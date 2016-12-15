package service

import (
	"fmt"

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
		v, ok := viewer.FromContext(ctx)
		if !ok {
			return nil, fmt.Errorf("could not fetch user")
		}
		var (
			config *kolide.AppConfig
			err    error
		)
		if v.IsAdmin() {
			config, err = svc.AppConfig(ctx)
			if err != nil {
				return getAppConfigResponse{Err: err}, nil
			}
		}

		response := getAppConfigResponse{
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
