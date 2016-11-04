package mysql

import (
	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewPack(pack *kolide.Pack) (e error) {
	return
}

func (d *Datastore) SavePack(pack *kolide.Pack) (e error) {
	return
}

func (d *Datastore) DeletePack(pid uint) (e error) {
	return
}

func (d *Datastore) Pack(pid uint) (p *kolide.Pack, e error) {
	return
}

func (d *Datastore) ListPacks(opt kolide.ListOptions) (packs []*kolide.Pack, e error) {
	return
}

func (d *Datastore) AddQueryToPack(qid uint, pid uint) (e error) {
	return
}

func (d *Datastore) ListQueriesInPack(pack *kolide.Pack) (queries []*kolide.Query, e error) {
	return
}

func (d *Datastore) RemoveQueryFromPack(query *kolide.Query, pack *kolide.Pack) (e error) {
	return
}

func (d *Datastore) AddLabelToPack(lid uint, pid uint) (e error) {
	return
}

func (d *Datastore) ListLabelsForPack(pack *kolide.Pack) (labels []*kolide.Label, e error) {
	return
}

func (d *Datastore) RemoveLabelFromPack(label *kolide.Label, pack *kolide.Pack) (e error) {
	return
}
