package mysql

import (
	"time"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewLabel(label *kolide.Label) (*kolide.Label, error) {
	panic("not implemented")
}

func (d *Datastore) DeleteLabel(lid uint) error {
	panic("not implemented")
}

func (d *Datastore) Label(lid uint) (*kolide.Label, error) {
	panic("not implemented")
}

func (d *Datastore) ListLabels(opt kolide.ListOptions) ([]*kolide.Label, error) {
	panic("not implemented")
}

func (d *Datastore) LabelQueriesForHost(host *kolide.Host, cutoff time.Time) (map[string]string, error) {
	panic("not implemented")
}

func (d *Datastore) RecordLabelQueryExecutions(host *kolide.Host, results map[string]bool, t time.Time) error {
	panic("not implemented")
}

func (d *Datastore) ListLabelsForHost(hid uint) ([]kolide.Label, error) {
	panic("not implemented")
}

func (d *Datastore) ListHostsInLabel(lid uint) ([]kolide.Host, error) {
	panic("not implemented")
}

func (d *Datastore) ListUniqueHostsInLabels(labels []uint) ([]kolide.Host, error) {
	panic("not implemented")
}

func (d *Datastore) SearchLabels(query string, omit []uint) ([]kolide.Label, error) {
	panic("not implemented")
}
