package mock

import "github.com/kolide/kolide-ose/server/kolide"

var _ kolide.InviteStore = (*MockInviteStore)(nil)

type NewInviteFunc func(i *kolide.Invite) (*kolide.Invite, error)

type MockInviteStore struct {
	NewInviteFunc        NewInviteFunc
	NewInviteFuncInvoked bool

	ListInvitesFunc        func(opt kolide.ListOptions) ([]*kolide.Invite, error)
	ListInvitesFuncInvoked bool

	InviteFunc        func(id uint) (*kolide.Invite, error)
	InviteFuncInvoked bool

	InviteByEmailFunc        func(email string) (*kolide.Invite, error)
	InviteByEmailFuncInvoked bool

	SaveInviteFunc        func(i *kolide.Invite) error
	SaveInviteFuncInvoked bool

	DeleteInviteFunc        func(i *kolide.Invite) error
	DeleteInviteFuncInvoked bool
}

func (s *MockInviteStore) NewInvite(i *kolide.Invite) (*kolide.Invite, error) {
	s.NewInviteFuncInvoked = true
	return s.NewInviteFunc(i)
}

func (s *MockInviteStore) ListInvites(opt kolide.ListOptions) ([]*kolide.Invite, error) {
	s.ListInvitesFuncInvoked = true
	return s.ListInvitesFunc(opt)
}

func (s *MockInviteStore) Invite(id uint) (*kolide.Invite, error) {
	s.InviteFuncInvoked = true
	return s.InviteFunc(id)
}

func (s *MockInviteStore) InviteByEmail(email string) (*kolide.Invite, error) {
	s.InviteByEmailFuncInvoked = true
	return s.InviteByEmailFunc(email)
}

func (s *MockInviteStore) SaveInvite(i *kolide.Invite) error {
	s.SaveInviteFuncInvoked = true
	return s.SaveInviteFunc(i)
}

func (s *MockInviteStore) DeleteInvite(i *kolide.Invite) error {
	s.DeleteInviteFuncInvoked = true
	return s.DeleteInviteFunc(i)
}

// helpers

func ReturnNewInivite(fake *kolide.Invite) NewInviteFunc {
	return func(i *kolide.Invite) (*kolide.Invite, error) {
		return fake, nil
	}
}
