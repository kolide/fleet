package service

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide/server/kolide"
)

type getIdentityProviderRequest struct {
	ID uint `json:"id"`
}

type getIdentityProviderResponse struct {
	IdentityProvider *kolide.IdentityProvider `json:"identity_provider,omitempty"`
	Err              error                    `json:"error,omitempty"`
}

func (r getIdentityProviderResponse) error() error { return r.Err }

func makeGetIdentityProviderEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getIdentityProviderRequest)
		idp, err := svc.GetIdentityProvider(ctx, req.ID)
		if err != nil {
			return getIdentityProviderResponse{Err: err}, nil
		}
		return getIdentityProviderResponse{IdentityProvider: idp}, nil
	}
}

func makeListIdentityProvidersEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

type newIdentityProviderRequest struct {
	Payload kolide.IdentityProviderPayload `json:"payload"`
}

type newIdentityProviderResponse struct {
	IdentityProvider *kolide.IdentityProvider `json:"identity_provider"`
	Err              error                    `json:"error,omitempty"`
}

func (r newIdentityProviderResponse) error() error { return r.Err }

func makeNewIdentityProviderEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(newIdentityProviderRequest)
		idp, err := svc.NewIdentityProvider(ctx, req.Payload)
		if err != nil {
			return newIdentityProviderResponse{Err: err}, nil
		}
		return newIdentityProviderResponse{IdentityProvider: idp}, nil
	}
}

func makeModifyIdentityProviderEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func makeDeleteIdentityProviderEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}
