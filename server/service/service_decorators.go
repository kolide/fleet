package service

import (
	"github.com/kolide/kolide/server/kolide"
	"golang.org/x/net/context"
)

func (svc service) ListDecorators(ctx context.Context) ([]*kolide.Decorator, error) {
	return svc.ds.ListDecorators()
}

func (svc service) DeleteDecorator(ctx context.Context, uid uint) error {
	return svc.ds.DeleteDecorator(uid)
}

func (svc service) NewDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	var dec kolide.Decorator
	dec.Query = payload.Query
	if payload.Interval != nil {
		dec.Interval = *payload.Interval
	}
	dec.Type, _ = kolide.DecoratorTypeFromName(payload.DecoratorType)
	return svc.ds.NewDecorator(&dec)
}
