package service

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

////////////////////////////////////////////////////////////////////////////////
// Enroll Agent
////////////////////////////////////////////////////////////////////////////////

type enrollAgentRequest struct {
	EnrollSecret   string `json:"enroll_secret"`
	HostIdentifier string `json:"host_identifier"`
}

type enrollAgentResponse struct {
	NodeKey string `json:"node_key,omitempty"`
	Err     error  `json:"error,omitempty"`
}

func (r enrollAgentResponse) error() error { return r.Err }

func makeEnrollAgentEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(enrollAgentRequest)
		nodeKey, err := svc.EnrollAgent(ctx, req.EnrollSecret, req.HostIdentifier)
		if err != nil {
			return enrollAgentResponse{Err: err}, nil
		}
		return enrollAgentResponse{NodeKey: nodeKey}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Client Config
////////////////////////////////////////////////////////////////////////////////

type getClientConfigRequest struct {
	NodeKey string `json:"node_key"`
}

type getClientConfigResponse struct {
	Config kolide.OsqueryConfig `json:"config,omitempty"`
	Err    error                `json:"error,omitempty"`
}

func (r getClientConfigResponse) error() error { return r.Err }

func makeGetClientConfigEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		config, err := svc.GetClientConfig(ctx)
		if err != nil {
			return getClientConfigResponse{Err: err}, nil
		}
		return getClientConfigResponse{Config: *config}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Distributed Queries
////////////////////////////////////////////////////////////////////////////////

type getDistributedQueriesRequest struct {
	NodeKey string `json:"node_key"`
}

type getDistributedQueriesResponse struct {
	Queries map[string]string `json:"queries"`
	Err     error             `json:"error,omitempty"`
}

func (r getDistributedQueriesResponse) error() error { return r.Err }

func makeGetDistributedQueriesEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		queries, err := svc.GetDistributedQueries(ctx)
		if err != nil {
			return getDistributedQueriesResponse{Err: err}, nil
		}
		return getDistributedQueriesResponse{Queries: queries}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Write Distributed Query Results
////////////////////////////////////////////////////////////////////////////////

type submitDistributedQueryResultsRequest struct {
	NodeKey string                                `json:"node_key"`
	Results kolide.OsqueryDistributedQueryResults `json:"queries"`
}

type submitDistributedQueryResultsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r submitDistributedQueryResultsResponse) error() error { return r.Err }

func makeSubmitDistributedQueryResultsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitDistributedQueryResultsRequest)
		err := svc.SubmitDistributedQueryResults(ctx, req.Results)
		if err != nil {
			return submitDistributedQueryResultsResponse{Err: err}, nil
		}
		return submitDistributedQueryResultsResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Submit Status Logs
////////////////////////////////////////////////////////////////////////////////

type submitStatusLogsRequest struct {
	NodeKey string `json:"node_key"`
	Logs    []kolide.OsqueryStatusLog
}

type submitStatusLogsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r submitStatusLogsResponse) error() error { return r.Err }

func makeSubmitStatusLogsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitStatusLogsRequest)
		err := svc.SubmitStatusLogs(ctx, req.Logs)
		if err != nil {
			return submitStatusLogsResponse{Err: err}, nil
		}
		return submitStatusLogsResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Submit Result Logs
////////////////////////////////////////////////////////////////////////////////

type submitResultLogsRequest struct {
	NodeKey string `json:"node_key"`
	Logs    []kolide.OsqueryResultLog
}

type submitResultLogsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r submitResultLogsResponse) error() error { return r.Err }

func makeSubmitResultLogsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitResultLogsRequest)
		err := svc.SubmitResultLogs(ctx, req.Logs)
		if err != nil {
			return submitResultLogsResponse{Err: err}, nil
		}
		return submitResultLogsResponse{}, nil
	}
}
