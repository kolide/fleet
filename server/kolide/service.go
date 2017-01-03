package kolide

import (
	"github.com/WatchBeam/clock"
)

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

	Clock() clock.Clock
}
