package service

import (
	"context"

	"github.com/kolide/kolide/server/kolide"
)

func (svc service) GetIdentityProvider(ctx context.Context, id uint) (*kolide.IdentityProvider, error) {
	return svc.ds.IdentityProvider(id)
}

func (svc service) ModifyIdentityProvider(ctx context.Context, id uint, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	return nil, nil
}

func (svc service) DeleteIdentityProvider(ctx context.Context, id uint) error {
	return nil
}

func (svc service) ListIdentityProviders(ctx context.Context, id uint) ([]kolide.IdentityProvider, error) {
	return nil, nil
}

func (svc service) NewIdentityProvider(ctx context.Context, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	var idp kolide.IdentityProvider
	if payload.Name != nil {
		idp.Name = *payload.Name
	}
	if payload.Certificate != nil {
		idp.Certificate = *payload.Certificate
	}
	if payload.DestinationURL != nil {
		idp.DestinationURL = *payload.DestinationURL
	}
	if payload.ImageURL != nil {
		idp.ImageURL = *payload.ImageURL
	}
	if payload.IssuerURI != nil {
		idp.IssuerURI = *payload.IssuerURI
	}
	if payload.Metadata != nil {
		idp.Metadata = *payload.Metadata
	}
	return svc.ds.NewIdentityProvider(idp)
}
