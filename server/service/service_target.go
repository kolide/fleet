package service

import (
	"golang.org/x/net/context"

    "github.com/kolide/kolide-ose/server/kolide"
)

func (svc service) SearchTargets(ctx context.Context, query string, omit []kolide.Target) (kolide.TargetSearchResults, uint, error) {
    return kolide.TargetSearchResults{}, 0, nil
}

func (svc service) CountHostsInTargets(ctx context.Context, targets []kolide.Target) (uint, error) {
    return 0, nil
}
