package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) SessionByKey(key string) (*kolide.Session, error) {
	panic("not implemented")
}

func (d *Datastore) SessionByID(id uint) (*kolide.Session, error) {
	panic("not implemented")
}

func (d *Datastore) ListSessionsForUser(id uint) ([]*kolide.Session, error) {
	panic("not implemented")
}

func (d *Datastore) NewSession(session *kolide.Session) (*kolide.Session, error) {
	panic("not implemented")
}

func (d *Datastore) DestroySession(session *kolide.Session) error {
	panic("not implemented")
}

func (d *Datastore) DestroyAllSessionsForUser(id uint) error {
	panic("not implemented")
}

func (d *Datastore) MarkSessionAccessed(session *kolide.Session) error {
	panic("not implemented")
}
