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

// UserStore contains methods for managing users in a datastore
type UserStore interface {
	NewUser(user *User) (*User, error)
	User(username string) (*User, error)
	UserByID(id uint) (*User, error)
	SaveUser(user *User) error
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
