package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/fleet/server/kolide"
)

////////////////////////////////////////////////////////////////////////////////
// Host response format
////////////////////////////////////////////////////////////////////////////////

type hostResponse struct {
	kolide.Host
	Status      string `json:"status"`
	DisplayText string `json:"display_text"`
}

func hostResponseForHost(ctx context.Context, svc kolide.API, host *kolide.Host) (*hostResponse, error) {
	return &hostResponse{
		Host:        *host,
		Status:      host.Status(time.Now()),
		DisplayText: host.HostName,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Get Host
////////////////////////////////////////////////////////////////////////////////

func (svc service) GetHost(ctx context.Context, id uint) (*kolide.Host, error) {
	return svc.ds.Host(id)
}

type getHostRequest struct {
	ID uint `json:"id"`
}

type getHostResponse struct {
	Host *hostResponse `json:"host"`
	Err  error         `json:"error,omitempty"`
}

func (r getHostResponse) error() error { return r.Err }

func makeGetHostEndpoint(svc kolide.API) endpoint.Endpoint {
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

func decodeGetHostRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return getHostRequest{ID: id}, nil
}

////////////////////////////////////////////////////////////////////////////////
// List Hosts
////////////////////////////////////////////////////////////////////////////////

func (svc service) ListHosts(ctx context.Context, opt kolide.ListOptions) ([]*kolide.Host, error) {
	return svc.ds.ListHosts(opt)
}

type listHostsRequest struct {
	ListOptions kolide.ListOptions
}

type listHostsResponse struct {
	Hosts []hostResponse `json:"hosts"`
	Err   error          `json:"error,omitempty"`
}

func (r listHostsResponse) error() error { return r.Err }

func makeListHostsEndpoint(svc kolide.API) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listHostsRequest)
		hosts, err := svc.ListHosts(ctx, req.ListOptions)
		if err != nil {
			return listHostsResponse{Err: err}, nil
		}

		hostResponses := make([]hostResponse, len(hosts), len(hosts))
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

func decodeListHostsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	opt, err := listOptionsFromRequest(r)
	if err != nil {
		return nil, err
	}
	return listHostsRequest{ListOptions: opt}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Delete Host
////////////////////////////////////////////////////////////////////////////////

func (svc service) DeleteHost(ctx context.Context, id uint) error {
	return svc.ds.DeleteHost(id)
}

type deleteHostRequest struct {
	ID uint `json:"id"`
}

type deleteHostResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteHostResponse) error() error { return r.Err }

func makeDeleteHostEndpoint(svc kolide.API) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteHostRequest)
		err := svc.DeleteHost(ctx, req.ID)
		if err != nil {
			return deleteHostResponse{Err: err}, nil
		}
		return deleteHostResponse{}, nil
	}
}

func decodeDeleteHostRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, err := idFromRequest(r, "id")
	if err != nil {
		return nil, err
	}
	return deleteHostRequest{ID: id}, nil
}

////////////////////////////////////////////////////////////////////////////////
// Get Host Summary
////////////////////////////////////////////////////////////////////////////////

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
