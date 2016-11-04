package mysql

import "github.com/kolide/kolide-ose/server/kolide"

func (d *Datastore) NewInvite(i *kolide.Invite) (*kolide.Invite, error) {
	panic("not implemented")
}

func (d *Datastore) ListInvites(opt kolide.ListOptions) ([]*kolide.Invite, error) {
	panic("not implemented")
}

func (d *Datastore) Invite(id uint) (*kolide.Invite, error) {
	panic("not implemented")
}

func (d *Datastore) InviteByEmail(email string) (*kolide.Invite, error) {
	panic("not implemented")
}

func (d *Datastore) SaveInvite(i *kolide.Invite) error {
	panic("not implemented")
}

func (d *Datastore) DeleteInvite(i *kolide.Invite) error {
	panic("not implemented")
}
