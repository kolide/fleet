package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) SessionByKey(key string) (sess *kolide.Session, e error) {
	return
}

func (d *Datastore) SessionByID(id uint) (sess *kolide.Session, e error) {
	return
}

func (d *Datastore) ListSessionsForUser(id uint) (sess []*kolide.Session, e error) {
	return
}

func (d *Datastore) NewSession(session *kolide.Session) (sess *kolide.Session, e error) {
	return
}

func (d *Datastore) DestroySession(session *kolide.Session) (e error) {
	return
}

func (d *Datastore) DestroyAllSessionsForUser(id uint) (e error) {
	return
}

func (d *Datastore) MarkSessionAccessed(session *kolide.Session) (e error) {
	return
}
