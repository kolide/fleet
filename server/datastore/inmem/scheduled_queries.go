package inmem

import (
	"sort"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewScheduledQuery(sq *kolide.ScheduledQuery) (*kolide.ScheduledQuery, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	newScheduledQuery := *sq

	newScheduledQuery.ID = d.nextID(newScheduledQuery)
	d.ScheduledQueries[newScheduledQuery.ID] = &newScheduledQuery

	newScheduledQuery.Query = d.Queries[newScheduledQuery.QueryID].Query
	newScheduledQuery.Name = d.Queries[newScheduledQuery.QueryID].Name

	return &newScheduledQuery, nil
}

func (d *Datastore) SaveScheduledQuery(sq *kolide.ScheduledQuery) (*kolide.ScheduledQuery, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if _, ok := d.ScheduledQueries[sq.ID]; !ok {
		return nil, notFound("ScheduledQuery").WithID(sq.ID)
	}

	d.ScheduledQueries[sq.ID] = sq
	return sq, nil
}

func (d *Datastore) DeleteScheduledQuery(id uint) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	if _, ok := d.ScheduledQueries[id]; !ok {
		return notFound("ScheduledQuery").WithID(id)
	}

	delete(d.ScheduledQueries, id)
	return nil
}

func (d *Datastore) ScheduledQuery(id uint) (*kolide.ScheduledQuery, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	sq, ok := d.ScheduledQueries[id]
	if !ok {
		return nil, notFound("ScheduledQuery").WithID(id)
	}

	sq.Name = d.Queries[sq.QueryID].Name
	sq.Query = d.Queries[sq.QueryID].Query

	return sq, nil
}

func (d *Datastore) ListScheduledQueriesInPack(id uint, opt kolide.ListOptions) ([]*kolide.ScheduledQuery, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	// We need to sort by keys to provide reliable ordering
	keys := []int{}
	for k, sq := range d.ScheduledQueries {
		if sq.PackID == id {
			keys = append(keys, int(k))
		}
	}

	if len(keys) == 0 {
		return []*kolide.ScheduledQuery{}, nil
	}

	sort.Ints(keys)

	scheduledQueries := []*kolide.ScheduledQuery{}
	for _, k := range keys {
		q := d.ScheduledQueries[uint(k)]
		q.Name = d.Queries[q.QueryID].Name
		q.Query = d.Queries[q.QueryID].Query
		scheduledQueries = append(scheduledQueries, q)
	}

	// Apply ordering
	if opt.OrderKey != "" {
		var fields = map[string]string{
			"id":         "ID",
			"created_at": "CreatedAt",
			"updated_at": "UpdatedAt",
			"name":       "Name",
			"query":      "Query",
			"interval":   "Interval",
			"snapshot":   "Snapshot",
			"removed":    "Removed",
			"platform":   "Platform",
			"version":    "Version",
		}
		if err := sortResults(scheduledQueries, opt, fields); err != nil {
			return nil, err
		}
	}

	// Apply limit/offset
	low, high := d.getLimitOffsetSliceBounds(opt, len(scheduledQueries))
	scheduledQueries = scheduledQueries[low:high]

	return scheduledQueries, nil
}
