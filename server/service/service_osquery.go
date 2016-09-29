package service

import (
	"encoding/json"
	"net/http"

	hostctx "github.com/kolide/kolide-ose/server/contexts/host"
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

func (e osqueryError) NodeInvalid() bool {
	return e.nodeInvalid
}

func (svc service) AuthenticateHost(ctx context.Context, nodeKey string) (*kolide.Host, error) {
	if nodeKey == "" {
		return nil, osqueryError{
			message:     "authentication error: missing node key",
			nodeInvalid: true,
		}
	}
	host, err := svc.ds.AuthenticateHost(nodeKey)
	if err != nil {
		return nil, osqueryError{
			message:     "authentication error: " + err.Error(),
			nodeInvalid: true,
		}
	}
	return host, nil
}

func (svc service) EnrollAgent(ctx context.Context, enrollSecret, hostIdentifier string) (string, error) {
	if enrollSecret != svc.config.Osquery.EnrollSecret {
		return "", osqueryError{message: "invalid enroll secret", nodeInvalid: true}
	}

	host, err := svc.ds.EnrollHost(hostIdentifier, "", "", "", svc.config.Osquery.NodeKeySize)
	if err != nil {
		return "", osqueryError{message: "enrollment failed: " + err.Error(), nodeInvalid: true}
	}

	return host.NodeKey, nil
}

func (svc service) GetClientConfig(ctx context.Context) (*kolide.OsqueryConfig, error) {
	host, ok := hostctx.FromContext(ctx)
	if !ok {
		return nil, osqueryError{message: "internal error: missing host from request context"}
	}

	config := new(kolide.OsqueryConfig)

	// we will need to give some subset of packs to this host based on the
	// labels which this host is known to belond to
	packs, err := svc.ds.Packs()
	if err != nil {
		return nil, osqueryError{message: "database error: " + err.Error()}
	}

	// pull the labels that this host belongs to
	labels, err := svc.ds.LabelsForHost(&host)
	if err != nil {
		return nil, osqueryError{message: "database error: " + err.Error()}
	}

	// in order to use o(1) array indexing in an o(n) loop vs a o(n^2) double
	// for loop iteration, we must create the array which may be indexed below
	var labelIDs map[uint]bool
	for _, label := range labels {
		labelIDs[label.ID] = true
	}

	for _, pack := range packs {
		// for each pack, we must know what labels have been assigned to that
		// pack
		labelsForPack, err := svc.ds.GetLabelsForPack(pack)
		if err != nil {
			return nil, osqueryError{message: "database error: " + err.Error()}
		}

		// o(n) iteration to determine whether or not a pack is enabled
		// in this case, n is len(labelsForPack)
		enabled := false
		for _, label := range labelsForPack {
			if labelIDs[label.ID] {
				enabled = true
				break
			}
		}

		// if the pack is enabled, we must add the content of the pack  the
		// osquery config struct which we will return
		if enabled {
			// first, we must figure out what queries are in this pack
			queries, err := svc.ds.GetQueriesInPack(pack)
			if err != nil {
				return nil, osqueryError{message: "database error: " + err.Error()}
			}

			// the serializable osquery config struct expects content in a
			// particular format, so we do the conversion here
			var configQueries kolide.Queries
			for _, query := range queries {
				configQueries[query.Name] = kolide.QueryContent{
					Query:    query.Query,
					Interval: query.Interval,
					Platform: query.Platform,
					Version:  query.Version,
					Snapshot: query.Snapshot,
				}
			}

			// finally, we add the pack to the client config struct with all of
			// the packs queries
			config.Packs[pack.Name] = kolide.PackContent{
				Platform: pack.Platform,
				Queries:  configQueries,
			}
		}
	}

	// add global options to the config struct to be returned
	config.Options = kolide.OsqueryOptions{
		PackDelimiter:      "/",
		DisableDistributed: false,
	}

	return config, nil
}

func (svc service) SubmitStatusLogs(ctx context.Context, logs []kolide.OsqueryStatusLog) error {
	host, ok := hostctx.FromContext(ctx)
	if !ok {
		return osqueryError{message: "internal error: missing host from request context"}
	}

	for _, log := range logs {
		err := json.NewEncoder(svc.osqueryStatusLogWriter).Encode(log)
		if err != nil {
			return errors.NewFromError(err, http.StatusInternalServerError, "error writing status log")
		}
	}

	err := svc.ds.MarkHostSeen(&host, svc.clock.Now())
	if err != nil {
		return osqueryError{message: "failed to update host seen: " + err.Error()}
	}

	return nil
}

func (svc service) SubmitResultLogs(ctx context.Context, logs []kolide.OsqueryResultLog) error {
	host, ok := hostctx.FromContext(ctx)
	if !ok {
		return osqueryError{message: "internal error: missing host from request context"}
	}

	for _, log := range logs {
		err := json.NewEncoder(svc.osqueryResultLogWriter).Encode(log)
		if err != nil {
			return errors.NewFromError(err, http.StatusInternalServerError, "error writing result log")
		}
	}

	err := svc.ds.MarkHostSeen(&host, svc.clock.Now())
	if err != nil {
		return osqueryError{message: "failed to update host seen: " + err.Error()}
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

func (svc service) GetDistributedQueries(ctx context.Context) (map[string]string, error) {
	queries := make(map[string]string)

	host, ok := hostctx.FromContext(ctx)
	if !ok {
		return nil, osqueryError{message: "internal error: missing host from request context"}
	}

	queries = hostDetailQueries(host)
	if len(queries) > 0 {
		// If the host details need to be updated, we should do so
		// before checking for any other queries
		return queries, nil
	}

	// Retrieve the label queries that should be updated
	cutoff := svc.clock.Now().Add(-svc.config.Osquery.LabelUpdateInterval)
	labelQueries, err := svc.ds.LabelQueriesForHost(&host, cutoff)
	if err != nil {
		return nil, err
	}

	for name, query := range labelQueries {
		queries[hostLabelQueryPrefix+name] = query
	}

	// TODO: retrieve the active distributed queries for this host

	return queries, nil
}

func (svc service) SubmitDistributedQueryResults(ctx context.Context, results kolide.OsqueryDistributedQueryResults) error {
	return nil
}
