package service

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) LicenseClaims(ctx context.Context) (*kolide.Claims, error) {
	license, err := svc.ds.License()
	if err != nil {
		return nil, err
	}
	claims, err := license.Claims()
	if err != nil {
		return nil, err
	}
	return claims, nil
}
