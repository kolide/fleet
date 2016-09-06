package server

import (
	"encoding/json"
	"fmt"

	"github.com/kolide/kolide-ose/errors"
	"github.com/kolide/kolide-ose/kolide"
	"golang.org/x/net/context"
)

func (svc service) EnrollAgent(ctx context.Context, enrollSecret, hostIdentifier string) (string, error) {
	if enrollSecret != svc.osqueryEnrollSecret {
		return "", errors.New(
			"Node key invalid",
			fmt.Sprintf("Invalid node key provided: %s", enrollSecret),
		)
	}

	host, err := svc.ds.EnrollHost(hostIdentifier, "", "", "", svc.osqueryNodeKeySize)
	if err != nil {
		return "", err
	}

	return host.NodeKey, nil
}

func (svc service) GetClientConfig(ctx context.Context, action string, data json.RawMessage) (kolide.OsqueryConfig, error) {
	var config kolide.OsqueryConfig
	return config, nil
}

func (svc service) SubmitStatusLogs(ctx context.Context, logs []kolide.OsqueryResultLog) error {
	return nil
}

func (svc service) SubmitResultsLogs(ctx context.Context, logs []kolide.OsqueryStatusLog) error {
	return nil
}

func (svc service) GetDistributedQueries(ctx context.Context) (map[string]string, error) {
	var queries map[string]string

	queries["id1"] = "select * from osquery_info"

	return queries, nil
}

func (svc service) SubmitDistributedQueryResults(ctx context.Context, results kolide.OsqueryDistributedQueryResults) error {
	return nil
}
