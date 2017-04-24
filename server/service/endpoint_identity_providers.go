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

type listIdentityProviderResponse struct {
	IdentityProviders []kolide.IdentityProvider `json:"identity_providers,omitempty"`
	Err               error                     `json:"error,omitempty"`
}

func (r listIdentityProviderResponse) error() error { return r.Err }

func makeListIdentityProvidersEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		idps, err := svc.ListIdentityProviders(ctx)
		if err != nil {
			return listIdentityProviderResponse{Err: err}, nil
		}
		return listIdentityProviderResponse{IdentityProviders: idps}, nil
	}
}

type listIdentityProvidersNoAuthResponse struct {
	IdentityProvidersNoAuth []kolide.IdentityProviderNoAuth `json:"identity_providers,omitempty"`
	Err                     error                           `json:"error,omitempty"`
}

func (r listIdentityProvidersNoAuthResponse) error() error { return r.Err }

func makeListIdentityProvidersNoAuthEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		idps, err := svc.ListIdentityProvidersNoAuth(ctx)
		if err != nil {
			return listIdentityProvidersNoAuthResponse{Err: err}, nil
		}
		return listIdentityProvidersNoAuthResponse{IdentityProvidersNoAuth: idps}, nil
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

type modifyIdentityProviderRequest struct {
	id      uint
	Payload *kolide.IdentityProviderPayload `json:"payload"`
}

type modifyIdentityProviderResponse struct {
	IdentityProvider *kolide.IdentityProvider `json:"identity_provider,omitempty"`
	Err              error                    `json:"error,omitempty"`
}

func (r modifyIdentityProviderResponse) error() error { return r.Err }

func makeModifyIdentityProviderEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyIdentityProviderRequest)
		idp, err := svc.ModifyIdentityProvider(ctx, req.id, *req.Payload)
		if err != nil {
			return modifyIdentityProviderResponse{Err: err}, nil
		}
		return modifyIdentityProviderResponse{IdentityProvider: idp}, nil
	}
}

type deleteIdentityProviderRequest struct {
	id uint
}

type deleteIdentityProviderResponse struct {
	Err error `json:"error"`
}

func (r deleteIdentityProviderResponse) error() error { return r.Err }

func makeDeleteIdentityProviderEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteIdentityProviderRequest)
		err := svc.DeleteIdentityProvider(ctx, req.id)
		return deleteIdentityProviderResponse{Err: err}, nil
	}
}
