package kolide

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

const LicenseTimeLayout = "2006-01-02 15:04:05 MST"

type LicenseStore interface {
	// SaveLicense saves jwt formatted customer license information
	SaveLicense(tokenString string) error
	// License returns a structure with the jwt customer license if it exists.
	License() (*License, error)
}

type LicenseService interface {
	// LicenseClaims returns details of a customer license that determine authorization
	// to use the Kolide product.
	LicenseClaims(ctx context.Context) (*Claims, error)
}

// Contains information needed to extract customer license particulars.
type License struct {
	UpdateTimestamp
	ID          uint
	TokenString *string `db:"license"`
	PublicKey   string  `db:"public_key"`
	Revoked     bool
}

// LicenseClaims contains information about the rights of a customer to
// use the Kolide product
type Claims struct {
	// Licensed indicates whether the application license has been installed. It does
	// not indicate whether or not the license is valid.
	Licensed         bool
	LicenseUUID      string
	OrganizationName string
	OrganizationUUID string
	// HostLimit the maximum number of hosts that a customer can use. 0 is unlimited.
	HostLimit int
	// Evaluation indicates that Kolide can be used for eval only.
	Evaluation bool
	// ExpiresAt time when license expires
	ExpiresAt time.Time
	// Revoked if true overrides a license that might otherwise be valid
	Revoked bool
}

// Claims returns information contained in the jwt license token
func (l License) Claims() (*Claims, error) {
	if l.TokenString == nil {
		return &Claims{Licensed: false}, nil
	}
	token, err := jwt.Parse(*l.TokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(l.PublicKey))
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	result := &Claims{Licensed: true}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.LicenseUUID = claims["license_uuid"].(string)
		result.OrganizationName = claims["organization_name"].(string)
		result.OrganizationUUID = claims["organization_uuid"].(string)
		result.HostLimit = int(claims["host_limit"].(float64))
		result.Evaluation = claims["evaluation"].(bool)
		result.Revoked = l.Revoked
		expiry, err := time.Parse(LicenseTimeLayout, claims["expires_at"].(string))
		if err != nil {
			return nil, err
		}
		result.ExpiresAt = expiry
	}
	return result, nil
}
