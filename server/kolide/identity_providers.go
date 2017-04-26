package kolide

import "context"

// IdentityProviderStore exposes methods to persist IdentityProviders.
// IdentityProvider is an entity used for single sign on.
type IdentityProviderStore interface {
	// NewIdentityProvider creates a new IdentityProvider.
	NewIdentityProvider(idp IdentityProvider) (*IdentityProvider, error)
	// SaveIdentityProvider saves changes to an IdentityProvider.
	SaveIdentityProvider(idb IdentityProvider) error
	// IdentityProvider retrieves an IdentityProvider identified by id.
	IdentityProvider(id uint) (*IdentityProvider, error)
	// DeleteIdentityProvider soft deletes an IdentityProvider
	DeleteIdentityProvider(id uint) error
	// ListIdentityProviders returns all IdentityProvider entities
	ListIdentityProviders() ([]IdentityProvider, error)
	// ListIdentityProvidersNoAuth returns IDP information with sensitive information
	// omitted.
	ListIdentityProvidersNoAuth() ([]IdentityProviderNoAuth, error)
}

// IdentityProvider represents a SAML identity provider.
type IdentityProvider struct {
	UpdateCreateTimestamps
	DeleteFields
	ID uint `json:"id"`
	// DestinationURL is the URL for the identity provider and is required
	// if Metadata is not present.
	DestinationURL string `json:"destination_url" db:"destination_url"`
	// IssuerURI is an optional identifier for this service provider. If not
	// supplied it will default to the host name.
	IssuerURI string `json:"issuer_uri" db:"issuer_uri"`
	// Certificate is the identity provider's public certificate
	// and is required if Metadata is not present.
	Certificate string `json:"cert" db:"cert"`
	// Name is the descriptive name for the identity provider that will
	// be displayed in the UI.
	Name string `json:"name"`
	// ImageURL is a link to an icon that will be displayed on the SSO
	// button for a particular identity provider.
	ImageURL string `json:"image_url" db:"image_url"`
	// Metadata if the definition of the SAML contract for a particular
	// IDP.
	Metadata string `json:"metadata" db:"metadata"`
	// MetadataURL is the discovery endpoint for the IDP, if a provider supplies
	// this it is preferable to the metadata field as it gets the latest state
	// for information on interacting with the IDP as opposed to the Metadata field which
	// is static.
	MetadataURL string `json:"metadata_url" db:"metadata_url"`
}

// IdentityProviderPayload user to update one or more fields of an IdentityProvider
// by supplying values that correspond to fields that will be changed.
type IdentityProviderPayload struct {
	DestinationURL *string `json:"destination_url"`
	IssuerURI      *string `json:"issuer_uri"`
	Certificate    *string `json:"cert"`
	Name           *string `json:"name"`
	ImageURL       *string `json:"image_url"`
	Metadata       *string `json:"metadata"`
	MetadataURL    *string `json:"metadata_url"`
}

// IdentityProviderNoAuth used to display SSO information for unauthorized users.
type IdentityProviderNoAuth struct {
	ID uint `json:"id"`
	// Name is the descriptive name for the identity provider that will
	// be displayed in the UI.
	Name string `json:"name"`
	// ImageURL is a link to an icon that will be displayed on the SSO
	// button for a particular identity provider.
	ImageURL string `json:"image_url" db:"image_url"`
}

// IdentityProviderService exposes methods to manage IdentityProvider entities
type IdentityProviderService interface {
	// NewIdentityProvider creates a IdentityProvider
	NewIdentityProvider(ctx context.Context, payload IdentityProviderPayload) (*IdentityProvider, error)
	// SaveIdentityProvider is used to modify an existing IdentityProvider.  Nonnil
	// fields in the payload argument will be changed for an existing IdentityProvider
	ModifyIdentityProvider(ctx context.Context, id uint, payload IdentityProviderPayload) (*IdentityProvider, error)
	// GetIdentityProvider retrieves an IdentityProvider given it's ID.
	GetIdentityProvider(ctx context.Context, id uint) (*IdentityProvider, error)
	// DeleteIdentityProvider removes an IdentityProvider
	DeleteIdentityProvider(ctx context.Context, id uint) error
	// ListIdentityProviders returns a list of all IdentityProvider entities
	ListIdentityProviders(ctx context.Context) ([]IdentityProvider, error)
	// ListIdentityProvidersNoAuth returns a list of identity providers with sensitive
	// security info removed so it can be used with unauthorized users.
	ListIdentityProvidersNoAuth(ctx context.Context) ([]IdentityProviderNoAuth, error)
}
