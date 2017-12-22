package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) ResetOptions(ctx context.Context) ([]kolide.Option, error) {
	return svc.ds.ResetOptions()
}

func (svc service) GetOptions(ctx context.Context) ([]kolide.Option, error) {
	opts, err := svc.ds.ListOptions()
	if err != nil {
		return nil, err
	}
	return opts, nil
}

func (svc service) ModifyOptions(ctx context.Context, req kolide.OptionRequest) ([]kolide.Option, error) {
	if err := svc.ds.SaveOptions(req.Options); err != nil {
		return nil, err
	}
	return req.Options, nil
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeModifyOptionsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req kolide.OptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type optionsResponse struct {
	Options []kolide.Option `json:"options,omitempty"`
	Err     error           `json:"error,omitempty"`
}

func (or optionsResponse) error() error { return or.Err }

func makeGetOptionsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		options, err := svc.GetOptions(ctx)
		if err != nil {
			return optionsResponse{Err: err}, nil
		}
		return optionsResponse{Options: options}, nil
	}
}

func makeModifyOptionsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		payload := request.(kolide.OptionRequest)
		opts, err := svc.ModifyOptions(ctx, payload)
		if err != nil {
			return optionsResponse{Err: err}, nil
		}
		return optionsResponse{Options: opts}, nil
	}
}

func makeResetOptionsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		options, err := svc.ResetOptions(ctx)
		if err != nil {
			return optionsResponse{Err: err}, nil
		}
		return optionsResponse{Options: options}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Metrics
////////////////////////////////////////////////////////////////////////////////

func (mw metricsMiddleware) GetOptions(ctx context.Context) ([]kolide.Option, error) {
	var (
		options []kolide.Option
		err     error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "GetOptions", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	options, err = mw.Service.GetOptions(ctx)
	return options, err
}

func (mw metricsMiddleware) ModifyOptions(ctx context.Context, or kolide.OptionRequest) ([]kolide.Option, error) {
	var (
		options []kolide.Option
		err     error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "ModifyOptions", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	options, err = mw.Service.ModifyOptions(ctx, or)
	return options, err
}

func (mw metricsMiddleware) ResetOptions(ctx context.Context) ([]kolide.Option, error) {
	var (
		options []kolide.Option
		err     error
	)
	defer func(begin time.Time) {
		lvs := []string{"method", "ResetOptions", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	options, err = mw.Service.ResetOptions(ctx)
	return options, err
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) GetOptions(ctx context.Context) ([]kolide.Option, error) {
	var (
		options []kolide.Option
		err     error
	)

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "GetOptions",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	options, err = mw.Service.GetOptions(ctx)
	return options, err
}

func (mw loggingMiddleware) ModifyOptions(ctx context.Context, req kolide.OptionRequest) ([]kolide.Option, error) {
	var (
		options []kolide.Option
		err     error
	)

	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ModifyOptions",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	options, err = mw.Service.ModifyOptions(ctx, req)
	return options, err
}

func (mw loggingMiddleware) ResetOptions(ctx context.Context) ([]kolide.Option, error) {
	var (
		options []kolide.Option
		err     error
	)
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "ResetOptions",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	options, err = mw.Service.ResetOptions(ctx)
	return options, err
}

////////////////////////////////////////////////////////////////////////////////
// Validation
////////////////////////////////////////////////////////////////////////////////

func (mw validationMiddleware) ModifyOptions(ctx context.Context, req kolide.OptionRequest) ([]kolide.Option, error) {
	invalid := &invalidArgumentError{}
	for _, opt := range req.Options {
		if opt.ReadOnly {
			invalid.Append(opt.Name, "readonly option")
			continue
		}
		if err := validateValueMapsToOptionType(opt); err != nil {
			invalid.Append(opt.Name, err.Error())
		}
	}
	if invalid.HasErrors() {
		return nil, invalid
	}
	return mw.Service.ModifyOptions(ctx, req)
}

var (
	errTypeMismatch = errors.New("type mismatch")
	errInvalidType  = errors.New("invalid option type")
)

func validateValueMapsToOptionType(opt kolide.Option) error {
	if !opt.OptionSet() {
		return nil
	}
	switch opt.GetValue().(type) {
	case int, uint, uint64, float64:
		if opt.Type != kolide.OptionTypeInt {
			return errTypeMismatch
		}
	case string:
		if opt.Type != kolide.OptionTypeString {
			return errTypeMismatch
		}
	case bool:
		if opt.Type != kolide.OptionTypeBool {
			return errTypeMismatch
		}
	default:
		return errInvalidType
	}
	return nil
}
