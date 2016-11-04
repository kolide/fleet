package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewUser(user *kolide.User) (*kolide.User, error) {
	panic("not implemented")
}
func (d *Datastore) User(username string) (*kolide.User, error) {
	panic("not implemented")
}

func (d *Datastore) ListUsers(opt kolide.ListOptions) ([]*kolide.User, error) {
	panic("not implemented")
}

func (d *Datastore) UserByEmail(email string) (*kolide.User, error) {
	panic("not implemented")
}

func (d *Datastore) UserByID(id uint) (*kolide.User, error) {
	panic("not implemented")
}

func (d *Datastore) SaveUser(user *kolide.User) error {
	panic("not implemented")
}
