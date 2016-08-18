package app

import "time"

// Datastore combines all methods for backend interactions
type Datastore interface {
	UserStore
	HostStore
	CampaignStore
	Drop() error
	Migrate() error
}

type HostStore interface {
	EnrollHost(uuid, hostname, ip, platform string, nodeKeySize int) (*Host, error)
}

type CampaignStore interface {
	CreatePassworResetRequest(userID uint, expires time.Time, token string) (*PasswordResetRequest, error)

	DeletePasswordResetRequest(req *PasswordResetRequest) error

	FindPassswordResetByID(id uint) (*PasswordResetRequest, error)

	FindPassswordResetByToken(token string) (*PasswordResetRequest, error)

	FindPassswordResetByTokenAndUserID(token string, id uint) (*PasswordResetRequest, error)
}
