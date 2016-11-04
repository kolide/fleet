package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewQuery(query *kolide.Query) (q *kolide.Query, e error) {
	return
}

func (d *Datastore) SaveQuery(query *kolide.Query) (e error) {
	return
}

func (d *Datastore) DeleteQuery(query *kolide.Query) (e error) {
	return
}

func (d *Datastore) Query(id uint) (q *kolide.Query, e error) {
	return
}

func (d *Datastore) ListQueries(opt kolide.ListOptions) (queries []*kolide.Query, e error) {
	return
}
