package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewPasswordResetRequest(req *kolide.PasswordResetRequest) (*kolide.PasswordResetRequest, error) {
	panic("not implemented")
}

func (d *Datastore) SavePasswordResetRequest(req *kolide.PasswordResetRequest) error {
	panic("not implemented")
}

func (d *Datastore) DeletePasswordResetRequest(req *kolide.PasswordResetRequest) error {
	panic("not implemented")
}

func (d *Datastore) DeletePasswordResetRequestsForUser(userID uint) error {
	panic("not implemented")
}

func (d *Datastore) FindPassswordResetByID(id uint) (*kolide.PasswordResetRequest, error) {
	panic("not implemented")
}

func (d *Datastore) FindPassswordResetsByUserID(id uint) ([]*kolide.PasswordResetRequest, error) {
	panic("not implemented")
}

func (d *Datastore) FindPassswordResetByToken(token string) (*kolide.PasswordResetRequest, error) {
	panic("not implemented")
}

func (d *Datastore) FindPassswordResetByTokenAndUserID(token string, id uint) (*kolide.PasswordResetRequest, error) {
	panic("not implemented")
}
