package kolide

import "errors"

var (
	// ErrNotFound is returned when the datastore resource cannot be found
	ErrNotFound = errors.New("resource not found")

	// ErrExists is returned when creating a datastore resource that already exists
	ErrExists = errors.New("resource already created")
)

// Datastore combines all the interfaces in the Kolide DAL
type Datastore interface {
	UserStore
	QueryStore
	CampaignStore
	PackStore
	LabelStore
	HostStore
	PasswordResetStore
	SessionStore
	AppConfigStore
	InviteStore
	Name() string
	Drop() error
	Migrate() error
}
