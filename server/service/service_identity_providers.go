package service

import (
	"context"

	"github.com/kolide/kolide/server/kolide"
	"github.com/pkg/errors"
)

func (svc service) GetIdentityProvider(ctx context.Context, id uint) (*kolide.IdentityProvider, error) {
	return svc.ds.IdentityProvider(id)
}

func (svc service) ModifyIdentityProvider(ctx context.Context, id uint, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	idp, err := svc.ds.IdentityProvider(id)
	if err != nil {
		return nil, errors.Wrap(err, "modifying identity provider")
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
	if payload.Name != nil {
		idp.Name = *payload.Name
	}
	err = svc.ds.SaveIdentityProvider(*idp)
	if err != nil {
		return nil, errors.Wrap(err, "modifying identity provider")
	}
	return idp, nil
}

func (svc service) DeleteIdentityProvider(ctx context.Context, id uint) error {
	return svc.ds.DeleteIdentityProvider(id)
}

func (svc service) ListIdentityProviders(ctx context.Context) ([]kolide.IdentityProvider, error) {
	return svc.ds.ListIdentityProviders()
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
