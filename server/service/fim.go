package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
	"github.com/pkg/errors"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) GetFIM(ctx context.Context) (*kolide.FIMConfig, error) {
	config, err := svc.ds.AppConfig()
	if err != nil {
		return nil, errors.Wrap(err, "getting fim config")
	}
	paths, err := svc.ds.FIMSections()
	if err != nil {
		return nil, errors.Wrap(err, "getting fim paths")
	}
	result := &kolide.FIMConfig{
		Interval:  uint(config.FIMInterval),
		FilePaths: paths,
	}
	return result, nil
}

// ModifyFIM will remove existing FIM settings and replace it
func (svc service) ModifyFIM(ctx context.Context, fim kolide.FIMConfig) error {
	if err := svc.ds.ClearFIMSections(); err != nil {
		return errors.Wrap(err, "updating fim")
	}
	config, err := svc.ds.AppConfig()
	if err != nil {
		return errors.Wrap(err, "updating fim")
	}
	config.FIMInterval = int(fim.Interval)
	for sectionName, paths := range fim.FilePaths {
		section := kolide.FIMSection{
			SectionName: sectionName,
			Paths:       paths,
		}
		if _, err := svc.ds.NewFIMSection(&section); err != nil {
			return errors.Wrap(err, "creating fim section")
		}
	}
	return svc.ds.SaveAppConfig(config)
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeModifyFIMRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var fimConfig kolide.FIMConfig
	if err := json.NewDecoder(r.Body).Decode(&fimConfig); err != nil {
		return nil, err
	}
	return fimConfig, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type modifyFIMResponse struct {
	Err error `json:"error,omitempty"`
}

func (m modifyFIMResponse) error() error { return m.Err }

func makeModifyFIMEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		fimConfig := req.(kolide.FIMConfig)
		var resp modifyFIMResponse
		if err := svc.ModifyFIM(ctx, fimConfig); err != nil {
			resp.Err = err
		}
		return resp, nil
	}
}

type getFIMResponse struct {
	Err     error             `json:"error,omitempty"`
	Payload *kolide.FIMConfig `json:"payload,omitempty"`
}

func makeGetFIMEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		fimConfig, err := svc.GetFIM(ctx)
		if err != nil {
			return getFIMResponse{Err: err}, nil
		}
		return getFIMResponse{Payload: fimConfig}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Metrics
////////////////////////////////////////////////////////////////////////////////

func (mw metricsMiddleware) GetFIM(ctx context.Context) (cfg *kolide.FIMConfig, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetFIM", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	cfg, err = mw.Service.GetFIM(ctx)
	return cfg, err
}

func (mw metricsMiddleware) ModifyFIM(ctx context.Context, fim kolide.FIMConfig) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "ModifyFIM", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Service.ModifyFIM(ctx, fim)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (lm loggingMiddleware) GetFIM(ctx context.Context) (cfg *kolide.FIMConfig, err error) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "GetFIM",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	cfg, err = lm.Service.GetFIM(ctx)
	return cfg, err
}

func (lm loggingMiddleware) ModifyFIM(ctx context.Context, fim kolide.FIMConfig) (err error) {
	defer func(begin time.Time) {
		lm.logger.Log(
			"method", "ModifyFIM",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	err = lm.Service.ModifyFIM(ctx, fim)
	return err
}
