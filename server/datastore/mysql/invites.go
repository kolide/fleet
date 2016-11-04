package mysql

import "github.com/kolide/kolide-ose/server/kolide"

func (d *Datastore) NewInvite(i *kolide.Invite) (invite *kolide.Invite, e error) {
	return
}

func (d *Datastore) ListInvites(opt kolide.ListOptions) (invites []*kolide.Invite, e error) {
	return
}

func (d *Datastore) Invite(id uint) (invite *kolide.Invite, e error) {
	return
}

func (d *Datastore) InviteByEmail(email string) (invite *kolide.Invite, e error) {
	return
}

func (d *Datastore) SaveInvite(i *kolide.Invite) (e error) {
	return
}

func (d *Datastore) DeleteInvite(i *kolide.Invite) (e error) {
	return
}
