package kolide

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/net/context"
)

const (
	LicenseTimeLayout  = "2006-01-02 15:04:05 MST"
	LicenseGracePeriod = time.Hour * 24 * 60 // 60 days
)

type LicenseStore interface {
	// SaveLicense saves jwt formatted customer license information
	SaveLicense(tokenString string) (*License, error)
	// License returns a structure with the jwt customer license if it exists.
	License() (*License, error)
}

type LicenseService interface {
	// LicenseClaims returns details of a customer license that determine authorization
	// to use the Kolide product.
	License(ctx context.Context) (*License, error)

	// SaveLicense writes jwt token string to database after performing
	// validation
	SaveLicense(ctx context.Context, jwtToken string) (*License, error)
}

// Contains information needed to extract customer license particulars.
type License struct {
	UpdateTimestamp
	ID uint
	// Token is a jwt token
	Token *string `db:"token"`
	// PublicKey is used to validate the Token and extract claims
	PublicKey string `db:"key"`
	// Revoked if true overrides a license that might otherwise be valid
	Revoked bool
	// HostCount is the count of enrolled hosts
	HostCount uint `db:"-"`
}

// LicenseClaims contains information about the rights of a customer to
// use the Kolide product
type Claims struct {
	LicenseUUID      string
	OrganizationName string
	OrganizationUUID string
	// HostLimit the maximum number of hosts that a customer can use. 0 is unlimited.
	HostLimit int
	// Evaluation indicates that Kolide can be used for eval only.
	Evaluation bool
	// ExpiresAt time when license expires
	ExpiresAt time.Time
}

// Expired returns true if the license is expired
func (c *Claims) Expired(current time.Time) bool {
	if c.Evaluation {
		if current.Sub(c.ExpiresAt) >= 0 {
			return true
		}
		return false
	}
	if current.Sub(c.ExpiresAt.Add(LicenseGracePeriod)) >= 0 {
		return true
	}
	return false
}

// Claims returns information contained in the jwt license token
func (l *License) Claims() (*Claims, error) {
	if l.Token == nil {
		return nil, errors.New("license missing")
	}
	token, err := jwt.Parse(*l.Token, func(token *jwt.Token) (interface{}, error) {
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
	var result Claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.LicenseUUID = claims["license_uuid"].(string)
		result.OrganizationName = claims["organization_name"].(string)
		result.OrganizationUUID = claims["organization_uuid"].(string)
		result.HostLimit = int(claims["host_limit"].(float64))
		result.Evaluation = claims["evaluation"].(bool)
		expiry, err := time.Parse(LicenseTimeLayout, claims["expires_at"].(string))
		if err != nil {
			return nil, err
		}
		result.ExpiresAt = expiry
	}
	return &result, nil
}
