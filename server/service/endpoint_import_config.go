package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

type importResponse struct {
	Response *kolide.ImportConfigResponse `json:"response,omitempty"`
	Err      error                        `json:"error,omitempty"`
}

func (ir importResponse) error() error { return ir.Err }

func makeImportConfig(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		config := request.(kolide.ImportConfig)
		resp, err := svc.ImportConfig(ctx, &config)
		if err != nil {
			return importResponse{Err: err}, nil
		}
		return importResponse{Response: resp}, nil
	}
}
