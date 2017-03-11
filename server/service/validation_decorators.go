package service

import (
	"github.com/kolide/kolide/server/kolide"
	"golang.org/x/net/context"
)

// NewDecorator validator checks to make sure that a valid decorator type exists and
// if the decorator is of an interval type, an interval value is present and is
// divisable by 60
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
func (mw validationMiddleware) NewDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	invalid := &invalidArgumentError{}
	decoratorType, err := kolide.DecoratorTypeFromName(payload.DecoratorType)
	if err != nil {
		invalid.Append("type", err.Error())
	}
	// If the type is interval, the interval value must be present and
	// must be divisible by 60
	if decoratorType == kolide.DecoratorInterval {
		if payload.Interval != nil && *payload.Interval != 0 {
			if *payload.Interval%60 != 0 {
				invalid.Append("interal", "interval value must be divisible by 60")
			}
		} else {
			invalid.Append("interval", "missing required argument")
		}
	}
	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.NewDecorator(ctx, payload)
}

func (mw validationMiddleware) DeleteDecorator(ctx context.Context, id uint) error {
	invalid := &invalidArgumentError{}
	dec, err := mw.ds.Decorator(id)
	if err != nil {
		return err
	}
	if dec.BuiltIn {
		invalid.Append("decorator", "read only")
		return invalid
	}
	return mw.Service.DeleteDecorator(ctx, id)
}
