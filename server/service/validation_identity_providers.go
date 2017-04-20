package service

import (
	"context"

	"github.com/kolide/kolide/server/kolide"
)

func (vm validationMiddleware) NewIdentityProvider(ctx context.Context, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	var invalid invalidArgumentError
	if payload.Name == nil {
		invalid.Append("name", "name of identity provider is required")
	}
	if payload.Metadata == nil {
		if payload.Certificate == nil {
			invalid.Append("cert", "must be defined if metadata is undefined")
		}
		if payload.DestinationURL == nil {
			invalid.Append("destination_url", "must be defined if metadata is undefined")
		}
		if payload.IssuerURI == nil {
			invalid.Append("issuer_url", "must be defined if metadata is undefined")
		}
	}
	if invalid.HasErrors() {
		return nil, invalid
	}
	return vm.Service.NewIdentityProvider(ctx, payload)
}
