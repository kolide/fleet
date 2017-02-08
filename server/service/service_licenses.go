package service

import (
	"github.com/kolide/kolide/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) License(ctx context.Context) (*kolide.License, error) {
	license, err := svc.ds.License()
	if err != nil {
		return nil, err
	}
	return license, nil
}

func (svc service) SaveLicense(ctx context.Context, jwtToken string) (*kolide.License, error) {
	publicKey, err := svc.ds.LicensePublicKey(jwtToken)
	if err != nil {
		return nil, err
	}
	updated, err := svc.ds.SaveLicense(jwtToken, publicKey)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (svc service) RequireLicense() (bool, error) {
	lic, err := svc.ds.License()
	if err != nil {
		return false, err
	}
	return lic.Token == nil, nil
}
