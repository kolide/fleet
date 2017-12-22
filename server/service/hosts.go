package service

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func (svc service) ListHosts(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Host, error) {
	return svc.ds.ListHosts(opt)
}

func (svc service) GetHost(ctx context.Context, id uint) (*kolide.Host, error) {
	return svc.ds.Host(id)
}

func (svc service) GetHostSummary(ctx context.Context) (*kolide.HostSummary, error) {
	online, offline, mia, new, err := svc.ds.GenerateHostStatusStatistics(svc.clock.Now())
	if err != nil {
		return nil, err
	}
	return &kolide.HostSummary{
		OnlineCount:  online,
		OfflineCount: offline,
		MIACount:     mia,
		NewCount:     new,
	}, nil
}

func (svc service) DeleteHost(ctx context.Context, id uint) error {
	return svc.ds.DeleteHost(id)
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func decodeGetHostRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getHostRequest{ID: id}, nil
}

func decodeDeleteHostRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteHostRequest{ID: id}, nil
}

func decodeListHostsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	opt, err := listOptionsFromRequest(r)
	if err != nil {
		return nil, err
	}
	return listHostsRequest{ListOptions: opt}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Endpoints
////////////////////////////////////////////////////////////////////////////////

type hostResponse struct {
	kolide.Host
	Status      string `json:"status"`
	DisplayText string `json:"display_text"`
}

func hostResponseForHost(ctx context.Context, svc kolide.Service, host *kolide.Host) (*hostResponse, error) {
	return &hostResponse{
		Host:        *host,
		Status:      host.Status(time.Now()),
		DisplayText: host.HostName,
	}, nil
}

type getHostRequest struct {
	ID uint `json:"id"`
}

type getHostResponse struct {
	Host *hostResponse `json:"host"`
	Err  error         `json:"error,omitempty"`
}

func (r getHostResponse) error() error { return r.Err }

func makeGetHostEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getHostRequest)
		host, err := svc.GetHost(ctx, req.ID)
		if err != nil {
			return getHostResponse{Err: err}, nil
		}

		resp, err := hostResponseForHost(ctx, svc, host)
		if err != nil {
			return getHostResponse{Err: err}, nil
		}

		return getHostResponse{
			Host: resp,
		}, nil
	}
}

type listHostsRequest struct {
	ListOptions kolide.ListOptions
}

type listHostsResponse struct {
	Hosts []hostResponse `json:"hosts"`
	Err   error          `json:"error,omitempty"`
}

func (r listHostsResponse) error() error { return r.Err }

func makeListHostsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listHostsRequest)
		hosts, err := svc.ListHosts(ctx, req.ListOptions)
		if err != nil {
			return listHostsResponse{Err: err}, nil
		}

		hostResponses := make([]hostResponse, len(hosts))
		for i, host := range hosts {
			h, err := hostResponseForHost(ctx, svc, host)
			if err != nil {
				return listHostsResponse{Err: err}, nil
			}

			hostResponses[i] = *h
		}
		return listHostsResponse{Hosts: hostResponses}, nil
	}
}

type getHostSummaryResponse struct {
	kolide.HostSummary
	Err error `json:"error,omitempty"`
}

func (r getHostSummaryResponse) error() error { return r.Err }

func makeGetHostSummaryEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		summary, err := svc.GetHostSummary(ctx)
		if err != nil {
			return getHostSummaryResponse{Err: err}, nil
		}

		resp := getHostSummaryResponse{
			HostSummary: *summary,
		}
		return resp, nil
	}
}

type deleteHostRequest struct {
	ID uint `json:"id"`
}

type deleteHostResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteHostResponse) error() error { return r.Err }

func makeDeleteHostEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteHostRequest)
		err := svc.DeleteHost(ctx, req.ID)
		if err != nil {
			return deleteHostResponse{Err: err}, nil
		}
		return deleteHostResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) ListHosts(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Host, error) {
	var (
		hosts []*kolide.Host
		err   error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "ListHosts",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	hosts, err = mw.Service.ListHosts(ctx, opt)
	return hosts, err
}

func (mw loggingMiddleware) GetHost(ctx context.Context, id uint) (*kolide.Host, error) {
	var (
		host *kolide.Host
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetHost",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	host, err = mw.Service.GetHost(ctx, id)
	return host, err
}

func (mw loggingMiddleware) GetHostSummary(ctx context.Context) (*kolide.HostSummary, error) {
	var (
		summary *kolide.HostSummary
		err     error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetHostSummary",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	summary, err = mw.Service.GetHostSummary(ctx)
	return summary, err
}

func (mw loggingMiddleware) DeleteHost(ctx context.Context, id uint) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "DeleteHost",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.DeleteHost(ctx, id)
	return err
}
