package mock

//go:generate mockimpl -o datastore_users.go "s *UserStore" "kolide.UserStore"
//go:generate mockimpl -o datastore_invites.go "s *InviteStore" "kolide.InviteStore"
//go:generate mockimpl -o datastore_appconfig.go "s *AppConfigStore" "kolide.AppConfigStore"
//go:generate mockimpl -o datastore_licenses.go "s *LicenseStore" "kolide.LicenseStore"
//go:generate mockimpl -o datastore_labels.go "s *LabelStore" "kolide.LabelStore"
//go:generate mockimpl -o datastore_decorators.go "s *DecoratorStore" "kolide.DecoratorStore"
//go:generate mockimpl -o datastore_options.go "s *OptionStore" "kolide.OptionStore"
//go:generate mockimpl -o datastore_packs.go "s *PackStore" "kolide.PackStore"
//go:generate mockimpl -o datastore_hosts.go "s *HostStore" "kolide.HostStore"
//go:generate mockimpl -o datastore_identity_providers.go "s *IdentityProviderStore" "kolide.IdentityProviderStore"

import "github.com/kolide/kolide/server/kolide"

var _ kolide.Datastore = (*Store)(nil)

type Store struct {
	kolide.CampaignStore
	kolide.SessionStore
	kolide.PasswordResetStore
	kolide.QueryStore
	kolide.ScheduledQueryStore
	kolide.FileIntegrityMonitoringStore
	kolide.YARAStore
	kolide.TargetStore
	AppConfigStore
	DecoratorStore
	HostStore
	InviteStore
	LabelStore
	LicenseStore
	OptionStore
	PackStore
	UserStore
	IdentityProviderStore
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
func (m *Store) MigrationStatus() (kolide.MigrationStatus, error) {
	return 0, nil
}
func (m *Store) Name() string {
	return "mock"
}
