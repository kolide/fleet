package inmem

import (
	"sort"

	"github.com/kolide/kolide-ose/server/kolide"
)

func (d *Datastore) NewPack(pack *kolide.Pack) (*kolide.Pack, error) {
	newPack := *pack

	for _, q := range d.Packs {
		if pack.Name == q.Name {
			return nil, alreadyExists("Pack", q.ID)
		}
	}

	d.mtx.Lock()
	newPack.ID = d.nextID(pack)
	d.Packs[newPack.ID] = &newPack
	d.mtx.Unlock()

	pack.ID = newPack.ID

	return pack, nil
}

func (d *Datastore) SavePack(pack *kolide.Pack) error {
	if _, ok := d.Packs[pack.ID]; !ok {
		return notFound("Pack").WithID(pack.ID)
	}

	d.mtx.Lock()
	d.Packs[pack.ID] = pack
	d.mtx.Unlock()

	return nil
}

func (d *Datastore) Pack(id uint) (*kolide.Pack, error) {
	p, err := d.byID(&kolide.Pack{ID: id})
	if err != nil {
		return nil, err
	}
	return p.(*kolide.Pack), nil
}

func (d *Datastore) ListPacks(opt kolide.ListOptions) ([]*kolide.Pack, error) {
	// We need to sort by keys to provide reliable ordering
	keys := []int{}
	d.mtx.Lock()
	for k, _ := range d.Packs {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	packs := []*kolide.Pack{}
	for _, k := range keys {
		packs = append(packs, d.Packs[uint(k)])
	}
	d.mtx.Unlock()

	// Apply ordering
	if opt.OrderKey != "" {
		var fields = map[string]string{
			"id":         "ID",
			"created_at": "CreatedAt",
			"updated_at": "UpdatedAt",
			"name":       "Name",
			"platform":   "Platform",
		}
		if err := sortResults(packs, opt, fields); err != nil {
			return nil, err
		}
	}

	// Apply limit/offset
	low, high := d.getLimitOffsetSliceBounds(opt, len(packs))
	packs = packs[low:high]

	return packs, nil
}

func (d *Datastore) AddLabelToPack(lid uint, pid uint) error {
	pt := &kolide.PackTarget{
		PackID: pid,
		Target: kolide.Target{
			Type:     kolide.TargetLabel,
			TargetID: lid,
		},
	}

	d.mtx.Lock()
	pt.ID = d.nextID(pt)
	d.PackTargets[pt.ID] = pt
	d.mtx.Unlock()

	return nil
}

func (d *Datastore) ListLabelsForPack(pid uint) ([]*kolide.Label, error) {
	var labels []*kolide.Label

	d.mtx.Lock()
	for _, pt := range d.PackTargets {
		if pt.Type == kolide.TargetLabel && pt.PackID == pid {
			labels = append(labels, d.Labels[pt.TargetID])
		}
	}
	d.mtx.Unlock()

	return labels, nil
}

func (d *Datastore) RemoveLabelFromPack(label *kolide.Label, pack *kolide.Pack) error {
	var labelsToDelete []uint

	d.mtx.Lock()
	for _, pt := range d.PackTargets {
		if pt.Type == kolide.TargetLabel && pt.TargetID == label.ID && pt.PackID == pack.ID {
			labelsToDelete = append(labelsToDelete, pt.ID)
		}
	}

	for _, id := range labelsToDelete {
		delete(d.PackTargets, id)
	}
	d.mtx.Unlock()

	return nil
}

func (d *Datastore) ListHostsInPack(pid uint, opt kolide.ListOptions) ([]*kolide.Host, error) {
	hosts := []*kolide.Host{}
	hostLookup := map[uint]bool{}

	d.mtx.Lock()
	for _, pt := range d.PackTargets {
		if pt.PackID != pid {
			continue
		}

		switch pt.Type {
		case kolide.TargetHost:
			if !hostLookup[pt.TargetID] {
				hostLookup[pt.TargetID] = true
				hosts = append(hosts, d.Hosts[pt.TargetID])
			}
		case kolide.TargetLabel:
			for _, lqe := range d.LabelQueryExecutions {
				if lqe.LabelID == pt.TargetID && lqe.Matches && !hostLookup[lqe.HostID] {
					hostLookup[lqe.HostID] = true
					hosts = append(hosts, d.Hosts[lqe.HostID])
				}
			}
		}
	}
	d.mtx.Unlock()

	// Apply ordering
	if opt.OrderKey != "" {
		var fields = map[string]string{
			"id":                 "ID",
			"created_at":         "CreatedAt",
			"updated_at":         "UpdatedAt",
			"detail_update_time": "DetailUpdateTime",
			"hostname":           "HostName",
			"uuid":               "UUID",
			"platform":           "Platform",
			"osquery_version":    "OsqueryVersion",
			"os_version":         "OSVersion",
			"uptime":             "Uptime",
			"memory":             "PhysicalMemory",
			"mac":                "PrimaryMAC",
			"ip":                 "PrimaryIP",
		}
		if err := sortResults(hosts, opt, fields); err != nil {
			return nil, err
		}
	}

	// Apply limit/offset
	low, high := d.getLimitOffsetSliceBounds(opt, len(hosts))
	hosts = hosts[low:high]

	return hosts, nil
}
