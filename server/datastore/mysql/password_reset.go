package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewPasswordResetRequest(req *kolide.PasswordResetRequest) (res *kolide.PasswordResetRequest, e error) {
	return
}

func (d *Datastore) SavePasswordResetRequest(req *kolide.PasswordResetRequest) (e error) {
	return
}

func (d *Datastore) DeletePasswordResetRequest(req *kolide.PasswordResetRequest) (e error) {
	return
}

func (d *Datastore) DeletePasswordResetRequestsForUser(userID uint) (e error) {
	return
}

func (d *Datastore) FindPassswordResetByID(id uint) (req *kolide.PasswordResetRequest, e error) {
	return
}

func (d *Datastore) FindPassswordResetsByUserID(id uint) (requests []*kolide.PasswordResetRequest, e error) {
	return
}

func (d *Datastore) FindPassswordResetByToken(token string) (req *kolide.PasswordResetRequest, e error) {
	return
}

func (d *Datastore) FindPassswordResetByTokenAndUserID(token string, id uint) (req *kolide.PasswordResetRequest, e error) {
	return
}
