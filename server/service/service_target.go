package service

import (
    "errors"

    "github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) SearchTargets(ctx context.Context, query string, omit []kolide.Target) (*kolide.TargetSearchResults, uint, error) {
    results := &kolide.TargetSearchResults{}

    var omitHosts []uint
    var omitLabels []uint
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

    return results, 0, nil
}

func (svc service) CountHostsInTargets(ctx context.Context, targets []kolide.Target) (uint, error) {
    return 0, errors.New("not implemented")
}
