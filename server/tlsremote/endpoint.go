// Package tlsremote implements the osquery TLS Remote service.
package tlsremote

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/go-kit/kit/endpoint"

	hostctx "github.com/kolide/fleet/server/contexts/host"
	"github.com/kolide/fleet/server/kolide"
)

// Endpoints collects all of the endpoints that compose an TLS service. It's meant to
// be used as a helper struct, to collect all of the endpoints into a single
// parameter.
type Endpoints struct {
	EnrollAgentEndpoint                   endpoint.Endpoint
	GetClientConfigEndpoint               endpoint.Endpoint
	GetDistributedQueriesEndpoint         endpoint.Endpoint
	SubmitDistributedQueryResultsEndpoint endpoint.Endpoint
	SubmitLogsEndpoint                    endpoint.Endpoint
}

func MakeServerEndpoints(svc kolide.OsqueryService) Endpoints {
	return Endpoints{
		EnrollAgentEndpoint:                   makeEnrollAgentEndpoint(svc),
		GetClientConfigEndpoint:               authenticatedHost(svc, makeGetClientConfigEndpoint(svc)),
		GetDistributedQueriesEndpoint:         authenticatedHost(svc, makeGetDistributedQueriesEndpoint(svc)),
		SubmitDistributedQueryResultsEndpoint: authenticatedHost(svc, makeSubmitDistributedQueryResultsEndpoint(svc)),
		SubmitLogsEndpoint:                    authenticatedHost(svc, makeSubmitLogsEndpoint(svc)),
	}
}

// authenticatedHost wraps an endpoint, checks the validity of the node_key
// provided in the request, and attaches the corresponding osquery host to the
// context for the request
func authenticatedHost(svc kolide.OsqueryService, next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		nodeKey, err := getNodeKey(request)
		if err != nil {
			return nil, err
		}

		host, err := svc.AuthenticateHost(ctx, nodeKey)
		if err != nil {
			return nil, err
		}

		ctx = hostctx.NewContext(ctx, *host)
		return next(ctx, request)
	}
}

func getNodeKey(r interface{}) (string, error) {
	// Retrieve node key by reflection (note that our options here
	// are limited by the fact that request is an interface{})
	v := reflect.ValueOf(r)
	if v.Kind() != reflect.Struct {
		return "", osqueryError{
			message: "request type is not struct. This is likely a Kolide programmer error.",
		}
	}
	nodeKeyField := v.FieldByName("NodeKey")
	if !nodeKeyField.IsValid() {
		return "", osqueryError{
			message: "request struct missing NodeKey. This is likely a Kolide programmer error.",
		}
	}
	if nodeKeyField.Kind() != reflect.String {
		return "", osqueryError{
			message: "NodeKey is not a string. This is likely a Kolide programmer error.",
		}
	}
	return nodeKeyField.String(), nil
}

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

func makeEnrollAgentEndpoint(svc kolide.OsqueryService) endpoint.Endpoint {
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
	kolide.OsqueryConfig
	Err error `json:"error,omitempty"`
}

func (r getClientConfigResponse) error() error { return r.Err }

func makeGetClientConfigEndpoint(svc kolide.OsqueryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		config, err := svc.GetClientConfig(ctx)
		if err != nil {
			return getClientConfigResponse{Err: err}, nil
		}
		return getClientConfigResponse{OsqueryConfig: *config}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Distributed Queries
////////////////////////////////////////////////////////////////////////////////

type getDistributedQueriesRequest struct {
	NodeKey string `json:"node_key"`
}

type getDistributedQueriesResponse struct {
	Queries    map[string]string `json:"queries"`
	Accelerate uint              `json:"accelerate,omitempty"`
	Err        error             `json:"error,omitempty"`
}

func (r getDistributedQueriesResponse) error() error { return r.Err }

func makeGetDistributedQueriesEndpoint(svc kolide.OsqueryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		queries, accelerate, err := svc.GetDistributedQueries(ctx)
		if err != nil {
			return getDistributedQueriesResponse{Err: err}, nil
		}
		return getDistributedQueriesResponse{Queries: queries, Accelerate: accelerate}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Write Distributed Query Results
////////////////////////////////////////////////////////////////////////////////

type submitDistributedQueryResultsRequest struct {
	NodeKey  string                                `json:"node_key"`
	Results  kolide.OsqueryDistributedQueryResults `json:"queries"`
	Statuses map[string]string                     `json:"statuses"`
}

type submitDistributedQueryResultsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r submitDistributedQueryResultsResponse) error() error { return r.Err }

func makeSubmitDistributedQueryResultsEndpoint(svc kolide.OsqueryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitDistributedQueryResultsRequest)
		err := svc.SubmitDistributedQueryResults(ctx, req.Results, req.Statuses)
		if err != nil {
			return submitDistributedQueryResultsResponse{Err: err}, nil
		}
		return submitDistributedQueryResultsResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Submit Logs
////////////////////////////////////////////////////////////////////////////////

type submitLogsRequest struct {
	NodeKey string          `json:"node_key"`
	LogType string          `json:"log_type"`
	Data    json.RawMessage `json:"data"`
}

type submitLogsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r submitLogsResponse) error() error { return r.Err }

func makeSubmitLogsEndpoint(svc kolide.OsqueryService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitLogsRequest)

		var err error
		switch req.LogType {
		case "status":
			var statuses []kolide.OsqueryStatusLog
			if err := json.Unmarshal(req.Data, &statuses); err != nil {
				err = osqueryError{message: "unmarshalling status logs: " + err.Error()}
				break
			}

			err = svc.SubmitStatusLogs(ctx, statuses)
			if err != nil {
				break
			}

		case "result":
			var results []json.RawMessage
			if err := json.Unmarshal(req.Data, &results); err != nil {
				err = osqueryError{message: "unmarshalling result logs: " + err.Error()}
				break
			}
			err = svc.SubmitResultLogs(ctx, results)
			if err != nil {
				break
			}

		default:
			err = osqueryError{message: "unknown log type: " + req.LogType}
		}

		return submitLogsResponse{Err: err}, nil
	}
}
