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

	Delete(ctx context.Context, entity Entity) (err error)
}
