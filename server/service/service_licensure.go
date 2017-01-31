package service

import (
	"github.com/kolide/kolide-ose/server/kolide"
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
	license, err := svc.ds.License()
	if err != nil {
		return nil, err
	}
	// check license validity
	license.Token = &jwtToken
	_, err = license.Claims()
	if err != nil {
		return nil, err
	}
	updated, err := svc.ds.SaveLicense(jwtToken)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
