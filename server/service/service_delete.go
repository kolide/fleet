package service

import (
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) Delete(ctx context.Context, entity kolide.Entity) error {
	return svc.ds.Delete(ctx, entity)
}
