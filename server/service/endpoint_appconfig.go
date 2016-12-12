package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

type modifyAppConfigResponse struct {
	Response *kolide.AppConfig `json:"app_config,omitempty"`
	Err      error             `json:"error,omitempty"`
}

func (m modifyAppConfigResponse) error() error { return m.Err }

type getAppConfigResponse struct {
	AppConfig *kolide.AppConfig `json:"app_config,omitempty"`
	Err       error             `json:"error,omitempty"`
}

func (r getAppConfigResponse) error() error { return r.Err }

func makeGetAppConfigEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		config, err := svc.AppConfig(ctx)

		if err != nil {
			return getAppConfigResponse{Err: err}, nil
		}

		return getAppConfigResponse{AppConfig: config}, nil
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
