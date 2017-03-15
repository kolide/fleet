package mock

//go:generate mockimpl -o datastore_users.go "s *UserStore" "kolide.UserStore"
//go:generate mockimpl -o datastore_invites.go "s *InviteStore" "kolide.InviteStore"
//go:generate mockimpl -o datastore_appconfig.go "s *AppConfigStore" "kolide.AppConfigStore"
//go:generate mockimpl -o datastore_licenses.go "s *LicenseStore" "kolide.LicenseStore"
//go:generate mockimpl -o datastore_labels.go "s *LabelStore" "kolide.LabelStore"
//go:generate mockimpl -o datastore_decorators.go "s *DecoratorStore" "kolide.DecoratorStore"
//go:generate mockimpl -o datastore_sessions.go "s *SessionStore" "kolide.SessionStore"
//go:generate mockimpl -o datastore_options.go "s *OptionStore" "kolide.OptionStore"
//go:generate mockimpl -o datastore_packs.go "s *PackStore" "kolide.PackStore"
//go:generate mockimpl -o datastore_queries.go "s *QueryStore" "kolide.QueryStore"
//go:generate mockimpl -o datastore_scheduled_queries.go "s *ScheduledQueryStore" "kolide.ScheduledQueryStore"
//go:generate mockimpl -o "datastore_file_integrity_monitoring.go" "s *FileIntegrityMonitoringStore" "kolide.FileIntegrityMonitoringStore"
//go:generate mockimpl -o "datastore_yara.go" "s *YARAStore" "kolide.YARAStore"

import "github.com/kolide/kolide/server/kolide"

var _ kolide.Datastore = (*Store)(nil)

type Store struct {
	kolide.HostStore
	kolide.CampaignStore
	kolide.PasswordResetStore
	LicenseStore
	InviteStore
	UserStore
	AppConfigStore
	LabelStore
	DecoratorStore
	SessionStore
	OptionStore
	PackStore
	QueryStore
	ScheduledQueryStore
	FileIntegrityMonitoringStore
	YARAStore
}

func (m *Store) Drop() error {
	return nil
}
func (m *Store) MigrateTables() error {
	return nil
}
func (m *Store) MigrateData() error {
	return nil
}
func (m *Store) Name() string {
	return "mock"
}

func (m *Store) MigrationStatus() (kolide.MigrationStatus, error) {
	return 0, nil
}
