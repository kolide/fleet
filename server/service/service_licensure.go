package service

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (svc service) LicenseClaims() (*kolide.LicenseClaims, error) {
	license, err := svc.ds.License()
	if err != nil {
		return nil, err
	}
	if license.TokenString == nil {
		return &kolide.LicenseClaims{Licensed: false}, nil
	}
	token, err := jwt.Parse(*license.TokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(license.PublicKey))
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	result := &kolide.LicenseClaims{Licensed: true}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.LicenseUUID = claims["license_uuid"].(string)
		result.OrganizationName = claims["organization_name"].(string)
		result.OrganizationUUID = claims["organization_uuid"].(string)
		result.HostLimit = int(claims["host_limit"].(float64))
		result.Evaluation = claims["evaluation"].(bool)
		result.Revoked = license.Revoked
		expiry, err := time.Parse(kolide.LicenseTimeLayout, claims["expires_at"].(string))
		if err != nil {
			return nil, err
		}
		result.ExpiresAt = expiry
	}

	return result, nil
}
