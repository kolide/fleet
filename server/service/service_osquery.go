package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kolide/kolide-ose/server/errors"
	"github.com/kolide/kolide-ose/server/kolide"
	"golang.org/x/net/context"
)

type osqueryError struct {
	message     string
	nodeInvalid bool
}

func (e osqueryError) Error() string {
	return e.message
}

func (svc service) EnrollAgent(ctx context.Context, enrollSecret, hostIdentifier string) (string, error) {
	if enrollSecret != svc.config.Osquery.EnrollSecret {
		return "", errors.New(
			"Node key invalid",
			fmt.Sprintf("Invalid node key provided: %s", enrollSecret),
		)
	}

	host, err := svc.ds.EnrollHost(hostIdentifier, "", "", "", svc.config.Osquery.NodeKeySize)
	if err != nil {
		return "", err
	}

	return host.NodeKey, nil
}

func (svc service) GetClientConfig(ctx context.Context, nodeKey string) (*kolide.OsqueryConfig, error) {
	_, err := svc.ds.AuthenticateHost(nodeKey)
	if err != nil {
		return nil, osqueryError{message: "authentication error", nodeInvalid: true}
	}

	var config kolide.OsqueryConfig
	return &config, nil
}

func (svc service) SubmitStatusLogs(ctx context.Context, nodeKey string, logs []kolide.OsqueryResultLog) error {
	_, err := svc.ds.AuthenticateHost(nodeKey)
	if err != nil {
		return osqueryError{message: "authentication error", nodeInvalid: true}
	}

	for _, log := range logs {
		err := json.NewEncoder(svc.osqueryStatusLogWriter).Encode(log)
		if err != nil {
			return errors.NewFromError(err, http.StatusInternalServerError, "error writing status log")
		}
	}
	return nil
}

func (svc service) SubmitResultsLogs(ctx context.Context, nodeKey string, logs []kolide.OsqueryStatusLog) error {
	_, err := svc.ds.AuthenticateHost(nodeKey)
	if err != nil {
		return osqueryError{message: "authentication error", nodeInvalid: true}
	}

	for _, log := range logs {
		err := json.NewEncoder(svc.osqueryResultsLogWriter).Encode(log)
		if err != nil {
			return errors.NewFromError(err, http.StatusInternalServerError, "error writing result log")
		}
	}
	return nil
}

// hostLabelQueryPrefix is appended before the query name when a query is
// provided as a label query. This allows the results to be retrieved when
// osqueryd writes the distributed query results.
const hostLabelQueryPrefix = "kolide_label_query_"

// hostDetailQueryPrefix is appended before the query name when a query is
// provided as a detail query.
const hostDetailQueryPrefix = "kolide_detail_query_"

// hostDetailQueries returns the map of queries that should be executed by
// osqueryd to fill in the host details
func hostDetailQueries(host kolide.Host) map[string]string {
	queries := make(map[string]string)
	if host.Platform == "" {
		queries[hostDetailQueryPrefix+"platform"] = "select build_platform from osquery_info;"
	}
	return queries
}

func (svc service) GetDistributedQueries(ctx context.Context, nodeKey string) (map[string]string, error) {
	queries := make(map[string]string)

	host, err := svc.ds.AuthenticateHost(nodeKey)
	if err != nil {
		return nil, osqueryError{message: "authentication error", nodeInvalid: true}
	}

	queries = hostDetailQueries(*host)
	if len(queries) > 0 {
		// If the host details need to be updated, we should do so
		// before checking for any other queries
		return queries, nil
	}

	// Retrieve the label queries that should be updated
	cutoff := svc.clock.Now().Add(-svc.config.Osquery.LabelUpdateInterval)
	labelQueries, err := svc.ds.LabelQueriesForHost(host, cutoff)
	if err != nil {
		return nil, err
	}

	for name, query := range labelQueries {
		queries[hostLabelQueryPrefix+name] = query
	}

	// TODO: retrieve the active distributed queries for this host

	return queries, nil
}

func (svc service) SubmitDistributedQueryResults(ctx context.Context, nodeKey string, results kolide.OsqueryDistributedQueryResults) error {
	_, err := svc.ds.AuthenticateHost(nodeKey)
	if err != nil {
		return osqueryError{message: "authentication error", nodeInvalid: true}
	}

	return nil
}
