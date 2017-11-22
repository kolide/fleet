package kolide

// service a interface stub
type Service interface {
	UserService
	SessionService
	PackService
	LabelService
	QueryService
	CampaignService
	HostService
	AppConfigService
	InviteService
	TargetService
	ScheduledQueryService
	OptionService
	ImportConfigService
	DecoratorService
	FileIntegrityMonitoringService
}
