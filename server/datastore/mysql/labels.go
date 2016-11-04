package mysql

import (
	"time"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewLabel(label *kolide.Label) (l *kolide.Label, e error) {
	return
}

func (d *Datastore) DeleteLabel(lid uint) (e error) {
	return
}

func (d *Datastore) Label(lid uint) (l *kolide.Label, e error) {
	return
}

func (d *Datastore) ListLabels(opt kolide.ListOptions) (labels []*kolide.Label, e error) {
	return
}

func (d *Datastore) LabelQueriesForHost(host *kolide.Host, cutoff time.Time) (q map[string]string, e error) {
	return
}

func (d *Datastore) RecordLabelQueryExecutions(host *kolide.Host, results map[string]bool, t time.Time) (e error) {
	return
}

func (d *Datastore) ListLabelsForHost(hid uint) (labels []kolide.Label, e error) {
	return
}

func (d *Datastore) ListHostsInLabel(lid uint) (hosts []kolide.Host, e error) {
	return
}

func (d *Datastore) ListUniqueHostsInLabels(labels []uint) (hosts []kolide.Host, e error) {
	return
}

func (d *Datastore) SearchLabels(query string, omit []uint) (labels []kolide.Label, e error) {
	return
}
