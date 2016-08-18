package app

import (
	"time"

	"github.com/kolide/kolide-ose/kolide"
)

// Datastore combines all methods for backend interactions
type Datastore interface {
	kolide.UserStore
	HostStore
	CampaignStore
	Drop() error
	Migrate() error
}

// HostStore enrolls hosts in the datastore
type HostStore interface {
	EnrollHost(uuid, hostname, ip, platform string, nodeKeySize int) (*Host, error)
}

// CampaignStore manages email campaigns in the database
type CampaignStore interface {
	CreatePassworResetRequest(userID uint, expires time.Time, token string) (*PasswordResetRequest, error)

	DeletePasswordResetRequest(req *PasswordResetRequest) error

	FindPassswordResetByID(id uint) (*PasswordResetRequest, error)

	FindPassswordResetByToken(token string) (*PasswordResetRequest, error)

	FindPassswordResetByTokenAndUserID(token string, id uint) (*PasswordResetRequest, error)
}
