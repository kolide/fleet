package inmem

import (
	"sort"

	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
)

func (orm *Datastore) NewScheduledQuery(sq *kolide.PackQuery) (*kolide.PackQuery, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	newScheduledQuery := *sq

	newScheduledQuery.ID = orm.nextID(newScheduledQuery)
	orm.packQueries[newScheduledQuery.ID] = &newScheduledQuery

	return &newScheduledQuery, nil
}

func (orm *Datastore) SaveScheduledQuery(sq *kolide.PackQuery) (*kolide.PackQuery, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if _, ok := orm.packQueries[sq.ID]; !ok {
		return nil, errors.ErrNotFound
	}

	orm.packQueries[sq.ID] = sq
	return sq, nil
}

func (orm *Datastore) DeleteScheduledQuery(id uint) error {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	if _, ok := orm.packQueries[id]; !ok {
		return errors.ErrNotFound
	}

	delete(orm.packQueries, id)
	return nil
}

func (orm *Datastore) ScheduledQuery(id uint) (*kolide.PackQuery, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	sq, ok := orm.packQueries[id]
	if !ok {
		return nil, errors.ErrNotFound
	}

	return sq, nil
}

func (orm *Datastore) ListScheduledQueriesInPack(id uint, opt kolide.ListOptions) ([]*kolide.PackQuery, error) {
	orm.mtx.Lock()
	defer orm.mtx.Unlock()

	// We need to sort by keys to provide reliable ordering
	keys := []int{}
	for k, sq := range orm.packQueries {
		if sq.PackID == id {
			keys = append(keys, int(k))
		}
	}

	if len(keys) == 0 {
		return []*kolide.PackQuery{}, nil
	}

	sort.Ints(keys)

	packQueries := []*kolide.PackQuery{}
	for _, k := range keys {
		packQueries = append(packQueries, orm.packQueries[uint(k)])
	}

	// Apply ordering
	if opt.OrderKey != "" {
		var fields = map[string]string{
			"id":           "ID",
			"created_at":   "CreatedAt",
			"updated_at":   "UpdatedAt",
			"name":         "Name",
			"query":        "Query",
			"interval":     "Interval",
			"snapshot":     "Snapshot",
			"differential": "Differential",
			"platform":     "Platform",
			"version":      "Version",
		}
		if err := sortResults(packQueries, opt, fields); err != nil {
			return nil, err
		}
	}

	// Apply limit/offset
	low, high := orm.getLimitOffsetSliceBounds(opt, len(packQueries))
	packQueries = packQueries[low:high]

	return packQueries, nil
}
