package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) ListDecorators(ctx context.Context) ([]*kolide.Decorator, error) {
	return svc.ds.ListDecorators()
}

func (svc service) DeleteDecorator(ctx context.Context, uid uint) error {
	return svc.ds.DeleteDecorator(uid)
}

func (svc service) NewDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	var dec kolide.Decorator
	if payload.Name != nil {
		dec.Name = *payload.Name
	}
	dec.Query = *payload.Query
	dec.Type = *payload.DecoratorType
	if payload.Interval != nil {
		dec.Interval = *payload.Interval
	}
	return svc.ds.NewDecorator(&dec)
}

func (svc service) ModifyDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	dec, err := svc.ds.Decorator(payload.ID)
	if err != nil {
		return nil, err
	}
	if payload.Name != nil {
		dec.Name = *payload.Name
	}
	if payload.DecoratorType != nil {
		dec.Type = *payload.DecoratorType
	}
	if payload.Query != nil {
		dec.Query = *payload.Query
	}
	if payload.Interval != nil {
		dec.Interval = *payload.Interval
	}
	err = svc.ds.SaveDecorator(dec)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeNewDecoratorRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var dec newDecoratorRequest
	err := json.NewDecoder(req.Body).Decode(&dec)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

func decodeDeleteDecoratorRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	id, err := idFromRequest(req, "id")
	if err != nil {
		return nil, err
	}
	return deleteDecoratorRequest{ID: id}, nil
}

func decodeModifyDecoratorRequest(ctx context.Context, req *http.Request) (interface{}, error) {
	var request newDecoratorRequest
	id, err := idFromRequest(req, "id")
	if err != nil {
		return nil, err
	}
	err = json.NewDecoder(req.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	request.Payload.ID = id
	return request, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type listDecoratorResponse struct {
	Decorators []*kolide.Decorator `json:"decorators"`
	Err        error               `json:"error,omitempty"`
}

func (r listDecoratorResponse) error() error { return r.Err }

func makeListDecoratorsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		decs, err := svc.ListDecorators(ctx)
		if err != nil {
			return listDecoratorResponse{Err: err}, nil
		}
		return listDecoratorResponse{Decorators: decs}, nil
	}
}

type newDecoratorRequest struct {
	Payload kolide.DecoratorPayload `json:"payload"`
}

type decoratorResponse struct {
	Decorator *kolide.Decorator `json:"decorator,omitempty"`
	Err       error             `json:"error,omitempty"`
}

func (r decoratorResponse) error() error { return r.Err }

func makeNewDecoratorEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(newDecoratorRequest)
		dec, err := svc.NewDecorator(ctx, r.Payload)
		if err != nil {
			return decoratorResponse{Err: err}, nil
		}
		return decoratorResponse{Decorator: dec}, nil
	}
}

type deleteDecoratorRequest struct {
	ID uint
}

type deleteDecoratorResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteDecoratorResponse) error() error { return r.Err }

func makeDeleteDecoratorEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(deleteDecoratorRequest)
		err := svc.DeleteDecorator(ctx, r.ID)

		if err != nil {
			return deleteDecoratorResponse{Err: err}, nil
		}
		return deleteDecoratorResponse{}, nil
	}
}

func makeModifyDecoratorEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r := request.(newDecoratorRequest)
		dec, err := svc.ModifyDecorator(ctx, r.Payload)
		if err != nil {
			return decoratorResponse{Err: err}, nil
		}
		return decoratorResponse{Decorator: dec}, nil

	}
}

////////////////////////////////////////////////////////////////////////////////
// Metrics
////////////////////////////////////////////////////////////////////////////////

func (mw metricsMiddleware) ListDecorators(ctx context.Context) ([]*kolide.Decorator, error) {
	var (
		decs []*kolide.Decorator
		err  error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "ListDecorators", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	decs, err = mw.Service.ListDecorators(ctx)
	return decs, err
}

func (mw metricsMiddleware) NewDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	var (
		dec *kolide.Decorator
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "NewDecorator", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	dec, err = mw.Service.NewDecorator(ctx, payload)
	return dec, err
}

func (mw metricsMiddleware) ModifyDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	var (
		dec *kolide.Decorator
		err error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "ModifyDecorator", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	dec, err = mw.Service.ModifyDecorator(ctx, payload)
	return dec, err
}

func (mw metricsMiddleware) DeleteDecorator(ctx context.Context, id uint) error {
	var err error
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteDecorator", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.DeleteDecorator(ctx, id)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) ListDecorators(ctx context.Context) ([]*kolide.Decorator, error) {
	var (
		decs []*kolide.Decorator
		err  error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ListDecorators",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	decs, err = mw.Service.ListDecorators(ctx)
	return decs, err
}

func (mw loggingMiddleware) NewDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	var (
		dec *kolide.Decorator
		err error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "NewDecorator",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	dec, err = mw.Service.NewDecorator(ctx, payload)
	return dec, err
}

func (mw loggingMiddleware) ModifyDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	var (
		dec *kolide.Decorator
		err error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ModifyDecorator",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	dec, err = mw.Service.ModifyDecorator(ctx, payload)
	return dec, err
}

func (mw loggingMiddleware) DeleteDecorator(ctx context.Context, id uint) error {
	var err error
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "DeleteDecorator",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = mw.Service.DeleteDecorator(ctx, id)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// Validation
////////////////////////////////////////////////////////////////////////////////

func validateNewDecoratorType(payload kolide.DecoratorPayload, invalid *invalidArgumentError) {
	if payload.DecoratorType == nil {
		invalid.Append("type", "missing required argument")
		return
	}
	if *payload.DecoratorType == kolide.DecoratorUndefined {
		invalid.Append("type", "invalid value, must be load, always, or interval")
		return
	}
	if *payload.DecoratorType == kolide.DecoratorInterval {
		if payload.Interval == nil {
			invalid.Append("interval", "missing required argument")
			return
		}
		if *payload.Interval%60 != 0 {
			invalid.Append("interval", "must be divisible by 60")
			return
		}
	}
}

// NewDecorator validator checks to make sure that a valid decorator type exists and
// if the decorator is of an interval type, an interval value is present and is
// divisable by 60
// See https://osquery.readthedocs.io/en/stable/deployment/configuration/
func (mw validationMiddleware) NewDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	invalid := &invalidArgumentError{}
	validateNewDecoratorType(payload, invalid)

	if payload.Query == nil {
		invalid.Append("query", "missing required argument")
	}

	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.NewDecorator(ctx, payload)
}

func (mw validationMiddleware) validateModifyDecoratorType(payload kolide.DecoratorPayload, invalid *invalidArgumentError) error {
	if payload.DecoratorType != nil {

		if *payload.DecoratorType == kolide.DecoratorUndefined {
			invalid.Append("type", "invalid value, must be load, always, or interval")
			return nil
		}
		if *payload.DecoratorType == kolide.DecoratorInterval {
			// special processing for interval type
			existingDec, err := mw.ds.Decorator(payload.ID)
			if err != nil {
				// if decorator is not present we want to return a 404 to the client
				return err
			}
			// if the type has changed from always or load to interval we need to
			// check suitability of interval value
			if existingDec.Type != kolide.DecoratorInterval {
				if payload.Interval == nil {
					invalid.Append("interval", "missing required argument")
					return nil
				}
			}
		}

		if payload.Interval != nil {
			if *payload.Interval%60 != 0 {
				invalid.Append("interval", "value must be divisible by 60")
			}
		}
	}
	return nil
}

func (mw validationMiddleware) ModifyDecorator(ctx context.Context, payload kolide.DecoratorPayload) (*kolide.Decorator, error) {
	invalid := &invalidArgumentError{}
	err := mw.validateModifyDecoratorType(payload, invalid)
	if err != nil {
		return nil, err
	}
	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.ModifyDecorator(ctx, payload)
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
