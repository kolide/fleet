package service

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) SearchTargets(ctx context.Context, query string, omit []kolide.Target) (*kolide.TargetSearchResults, uint, error) {
	results := &kolide.TargetSearchResults{}

	// assemble the omit sets for the calls to the individual datastore methods
	omitHosts := []uint{}
	omitLabels := []uint{}

	for _, omitTarget := range omit {
		switch omitTarget.Type {
		case kolide.TargetHost:
			omitHosts = append(omitHosts, omitTarget.TargetID)
		case kolide.TargetLabel:
			omitLabels = append(omitLabels, omitTarget.TargetID)
		}
	}

	hosts, err := svc.ds.SearchHosts(query, omitHosts)
	if err != nil {
		return nil, 0, err
	}
	results.Hosts = hosts

	labels, err := svc.ds.SearchLabels(query, omitLabels)
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
