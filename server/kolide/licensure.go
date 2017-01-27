package kolide

import "time"

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
	LicenseClaims() (*LicenseClaims, error)
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
type LicenseClaims struct {
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
