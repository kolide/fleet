package mock

import "github.com/kolide/kolide-ose/server/kolide"

type Store struct {
	kolide.HostStore
	kolide.LabelStore
	kolide.PackStore
	kolide.CampaignStore
	kolide.SessionStore
	kolide.AppConfigStore
	kolide.PasswordResetStore
	kolide.QueryStore

	MockInviteStore
	UserStore
}

func (m *Store) Drop() error {
	return nil
}
func (m *Store) Migrate() error {
	return nil
}
func (m *Store) Name() string {
	return "mock"
}
