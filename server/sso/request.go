package sso

type RequestSettings struct {
	PublicCert                string
	PrivateKey                string
	IdentityProvderPublicCert string
	IdentityProviderURL       string
	CallbackURL               string
	SignRequest               bool
}
