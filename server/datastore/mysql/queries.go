package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewQuery(query *kolide.Query) (*kolide.Query, error) {
	panic("not implemented")
}

func (d *Datastore) SaveQuery(query *kolide.Query) error {
	panic("not implemented")
}

func (d *Datastore) DeleteQuery(query *kolide.Query) error {
	panic("not implemented")
}

func (d *Datastore) Query(id uint) (*kolide.Query, error) {
	panic("not implemented")
}

func (d *Datastore) ListQueries(opt kolide.ListOptions) ([]*kolide.Query, error) {
	panic("not implemented")
}
