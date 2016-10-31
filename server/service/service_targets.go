package service

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) SearchTargets(ctx context.Context, query string, selectedHostIDs []uint, selectedLabelIDs []uint) (*kolide.TargetSearchResults, uint, error) {
	results := &kolide.TargetSearchResults{}

	hosts, err := svc.ds.SearchHosts(query, selectedHostIDs)
	if err != nil {
		return nil, 0, err
	}
	results.Hosts = hosts

	labels, err := svc.ds.SearchLabels(query, selectedLabelIDs)
	if err != nil {
		return nil, 0, err
	}
	results.Labels = labels

	count, err := svc.CountHostsInTargets(ctx, hosts, labels)
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

func (svc service) CountHostsInTargets(ctx context.Context, hosts []kolide.Host, labels []kolide.Label) (uint, error) {
	// make a lookup map for constant time deduplication
	hostLookup := map[uint]bool{}
	for _, host := range hosts {
		hostLookup[host.ID] = true
	}

	count := uint(len(hostLookup))

	for _, label := range labels {
		hostsInLabel, err := svc.ds.ListHostsInLabel(label.ID)
		if err != nil {
			return 0, err
		}
		for _, host := range hostsInLabel {
			if !hostLookup[host.ID] {
				hostLookup[host.ID] = true
				count++
			}
		}
	}

	return count, nil
}
