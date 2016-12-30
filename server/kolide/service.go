package kolide

import "golang.org/x/net/context"

// service a interface stub
type Service interface {
	UserService
	SessionService
	PackService
	LabelService
	QueryService
	CampaignService
	OsqueryService
	HostService
	AppConfigService
	InviteService
	TargetService
	ScheduledQueryService
	OptionService
	Deleter
}

// Deleter removes an Entity from the application datastore.
// Deleter should be implemented using a SoftDelete strategy, marking objects
// as deleted, but not necessarily removing the Entity permanently.
type Deleter interface {
	Delete(ctx context.Context, entity Entity) (err error)
}
