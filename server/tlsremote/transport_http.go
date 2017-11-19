package tlsremote

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/kolide/fleet/server/kolide"
)

// NewHTTPHandler returns an osquery TLS Remote HTTP Handler.
func NewHTTPHandler(e Endpoints, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
		kithttp.ServerAfter(kithttp.SetContentType("application/json; charset=utf-8")),
	}
	newServer := func(e endpoint.Endpoint, decodeFn kithttp.DecodeRequestFunc) http.Handler {
		return kithttp.NewServer(e, decodeFn, encodeResponse, opts...)
	}

	enrollAgentHander := newServer(
		e.EnrollAgentEndpoint,
		decodeEnrollAgentRequest,
	)

	getClientConfighandler := newServer(
		e.GetClientConfigEndpoint,
		decodeGetClientConfigRequest,
	)

	getDistributedQueriesHandler := newServer(
		e.GetDistributedQueriesEndpoint,
		decodeGetDistributedQueriesRequest,
	)

	submitDistributedQueryResultsHandler := newServer(
		e.SubmitDistributedQueryResultsEndpoint,
		decodeSubmitDistributedQueryResultsRequest,
	)

	submitLogsHandler := newServer(
		e.SubmitLogsEndpoint,
		decodeSubmitLogsRequest,
	)

	r := mux.NewRouter()
	r.Handle("/api/v1/osquery/enroll", enrollAgentHander).
		Methods("POST").
		Name("enroll_agent")
	r.Handle("/api/v1/osquery/config", getClientConfighandler).
		Methods("POST").
		Name("get_client_config")
	r.Handle("/api/v1/osquery/distributed/read", getDistributedQueriesHandler).
		Methods("POST").
		Name("get_distributed_queries")
	r.Handle("/api/v1/osquery/distributed/write", submitDistributedQueryResultsHandler).
		Methods("POST").
		Name("submit_distributed_query_results")
	r.Handle("/api/v1/osquery/log", submitLogsHandler).
		Methods("POST").
		Name("submit_logs")

	return r
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(response)
}

// erroer interface is implemented by response structs to encode business logic errors
type errorer interface {
	error() error
}

// encode error and status header to the client
func encodeError(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")

	type osqueryError interface {
		error
		NodeInvalid() bool
	}
	if e, ok := err.(osqueryError); ok {
		// osquery expects to receive the node_invalid key when a TLS
		// request provides an invalid node_key for authentication. It
		// doesn't use the error message provided, but we provide this
		// for debugging purposes (and perhaps osquery will use this
		// error message in the future).

		errMap := map[string]interface{}{"error": e.Error()}
		if e.NodeInvalid() {
			w.WriteHeader(http.StatusUnauthorized)
			errMap["node_invalid"] = true
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		enc.Encode(errMap)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	enc.Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func decodeEnrollAgentRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req enrollAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return req, nil
}

func decodeGetClientConfigRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getClientConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return req, nil
}

func decodeGetDistributedQueriesRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req getDistributedQueriesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return req, nil
}

func decodeSubmitDistributedQueryResultsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	// When a distributed query has no results, the JSON schema is
	// inconsistent, so we use this shim and massage into a consistent
	// schema. For example (simplified from actual osqueryd 1.8.2 output):
	// {
	// "queries": {
	//   "query_with_no_results": "", // <- Note string instead of array
	//   "query_with_results": [{"foo":"bar","baz":"bang"}]
	//  },
	// "node_key":"IGXCXknWQ1baTa8TZ6rF3kAPZ4\/aTsui"
	// }
	type distributedQueryResultsShim struct {
		NodeKey  string                     `json:"node_key"`
		Results  map[string]json.RawMessage `json:"queries"`
		Statuses map[string]string          `json:"statuses"`
	}

	var shim distributedQueryResultsShim
	if err := json.NewDecoder(r.Body).Decode(&shim); err != nil {
		return nil, err
	}
	defer r.Body.Close()

	results := kolide.OsqueryDistributedQueryResults{}
	for query, raw := range shim.Results {
		queryResults := []map[string]string{}
		// No need to handle error because the empty array is what we
		// want if there was an error parsing the JSON (the error
		// indicates that osquery sent us incosistently schemaed JSON)
		_ = json.Unmarshal(raw, &queryResults)
		results[query] = queryResults
	}

	req := submitDistributedQueryResultsRequest{
		NodeKey:  shim.NodeKey,
		Results:  results,
		Statuses: shim.Statuses,
	}

	return req, nil
}

func decodeSubmitLogsRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var err error
	body := r.Body
	if r.Header.Get("content-encoding") == "gzip" {
		body, err = gzip.NewReader(body)
		if err != nil {
			return nil, errors.Wrap(err, "decoding gzip")
		}
		defer body.Close()
	}

	var req submitLogsRequest
	if err = json.NewDecoder(body).Decode(&req); err != nil {
		return nil, errors.Wrap(err, "decoding JSON")
	}
	defer r.Body.Close()

	return req, nil
}
