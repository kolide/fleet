package app

import "github.com/kolide/kolide-ose/kolide"

// Datastore combines all methods for backend interactions
type Datastore interface {
	kolide.UserStore
	kolide.HostStore
	kolide.CampaignStore
	Drop() error
	Migrate() error
}
