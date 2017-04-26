package service

import (
	"context"

	"github.com/kolide/kolide/server/kolide"
)

func (vm validationMiddleware) NewIdentityProvider(ctx context.Context, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	var invalid invalidArgumentError
	validateIdentityProvider(&payload, &invalid)
	if invalid.HasErrors() {
		return nil, invalid
	}
	return vm.Service.NewIdentityProvider(ctx, payload)
}

func (vm validationMiddleware) ModifyIdentityProvider(ctx context.Context, id uint, payload kolide.IdentityProviderPayload) (*kolide.IdentityProvider, error) {
	var invalid invalidArgumentError
	validateIdentityProvider(&payload, &invalid)
	if invalid.HasErrors() {
		return nil, invalid
	}
	return vm.Service.ModifyIdentityProvider(ctx, id, payload)
}

func validateIdentityProvider(payload *kolide.IdentityProviderPayload, invalid *invalidArgumentError) {
	if payload.Name == nil {
		invalid.Append("name", "name of identity provider is required")
	}
	if payload.Name != nil && len(*payload.Name) == 0 {
		invalid.Append("name", "can't be empty")
	}
	if payload.Metadata == nil && payload.MetadataURL == nil {
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
	// can only define one or the other
	if payload.Metadata != nil && payload.MetadataURL != nil {
		invalid.Append("metadata", "defining both metadata and metadata url is not allowed")
	}
}
