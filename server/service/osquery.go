package service

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	hostctx "github.com/kolide/fleet/server/contexts/host"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/pubsub"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

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

// Sometimes osquery gives us empty string where we expect an integer.
// We change the to "0" so it can be handled by the appropriate string to
// integer conversion function, as these will err on ""
func emptyToZero(val string) string {
	if val == "" {
		return "0"
	}
	return val
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
		switch err.(type) {
		case kolide.NotFoundError:
			return nil, osqueryError{
				message:     "authentication error: invalid node key: " + nodeKey,
				nodeInvalid: true,
			}
		default:
			return nil, osqueryError{
				message: "authentication error: " + err.Error(),
			}
		}
	}

	// Update the "seen" time used to calculate online status
	err = svc.ds.MarkHostSeen(host, svc.clock.Now())
	if err != nil {
		return nil, osqueryError{message: "failed to mark host seen: " + err.Error()}
	}

	return host, nil
}

func (svc service) EnrollAgent(ctx context.Context, enrollSecret, hostIdentifier string) (string, error) {
	config, err := svc.ds.AppConfig()
	if err != nil {
		return "", osqueryError{message: "getting enroll secret: " + err.Error(), nodeInvalid: true}
	}

	if enrollSecret != config.EnrollSecret {
		return "", osqueryError{message: "invalid enroll secret", nodeInvalid: true}
	}

	host, err := svc.ds.EnrollHost(hostIdentifier, svc.config.Osquery.NodeKeySize)
	if err != nil {
		return "", osqueryError{message: "enrollment failed: " + err.Error(), nodeInvalid: true}
	}

	return host.NodeKey, nil
}

func (svc service) getFIMConfig(ctx context.Context, cfg *kolide.OsqueryConfig) (*kolide.OsqueryConfig, error) {
	fimConfig, err := svc.GetFIM(ctx)
	if err != nil {
		return nil, osqueryError{message: "internal error: unable to fetch FIM configuration"}
	}
	if cfg.Schedule == nil {
		cfg.Schedule = make(map[string]kolide.QueryContent)
	}
	removed := false
	// file events scheduled query is required to run file integrity monitors
	cfg.Schedule["file_events"] = kolide.QueryContent{
		Query:    "SELECT * FROM file_events;",
		Interval: fimConfig.Interval,
		Removed:  &removed,
	}
	cfg.FilePaths = fimConfig.FilePaths
	return cfg, nil
}

func (svc service) GetClientConfig(ctx context.Context) (*kolide.OsqueryConfig, error) {
	host, ok := hostctx.FromContext(ctx)
	if !ok {
		return nil, osqueryError{message: "internal error: missing host from request context"}
	}

	options, err := svc.ds.GetOsqueryConfigOptions()
	if err != nil {
		return nil, osqueryError{message: "internal error: unable to fetch configuration options"}
	}

	decorators, err := svc.ds.ListDecorators()
	if err != nil {
		return nil, osqueryError{message: "internal error: unable to fetch decorators"}
	}
	decConfig := kolide.Decorators{
		Interval: make(map[string][]string),
	}
	for _, dec := range decorators {
		switch dec.Type {
		case kolide.DecoratorLoad:
			decConfig.Load = append(decConfig.Load, dec.Query)
		case kolide.DecoratorAlways:
			decConfig.Always = append(decConfig.Always, dec.Query)
		case kolide.DecoratorInterval:
			key := strconv.Itoa(int(dec.Interval))
			decConfig.Interval[key] = append(decConfig.Interval[key], dec.Query)
		default:
			svc.logger.Log("component", "service", "method", "GetClientConfig", "err",
				"unknown decorator type")
		}
	}

	config := &kolide.OsqueryConfig{
		Options:    options,
		Decorators: decConfig,
		Packs:      kolide.Packs{},
	}

	packs, err := svc.ListPacksForHost(ctx, host.ID)
	if err != nil {
		return nil, osqueryError{message: "database error: " + err.Error()}
	}

	for _, pack := range packs {
		// first, we must figure out what queries are in this pack
		queries, err := svc.ds.ListScheduledQueriesInPack(pack.ID, kolide.ListOptions{})
		if err != nil {
			return nil, osqueryError{message: "database error: " + err.Error()}
		}

		// the serializable osquery config struct expects content in a
		// particular format, so we do the conversion here
		configQueries := kolide.Queries{}
		for _, query := range queries {
			queryContent := kolide.QueryContent{
				Query:    query.Query,
				Interval: query.Interval,
				Platform: query.Platform,
				Version:  query.Version,
			}

			if query.Snapshot != nil && *query.Snapshot {
				queryContent.Snapshot = query.Snapshot
			}

			configQueries[query.Name] = queryContent
		}

		// finally, we add the pack to the client config struct with all of
		// the packs queries
		config.Packs[pack.Name] = kolide.PackContent{
			Platform: pack.Platform,
			Queries:  configQueries,
		}
	}

	// Save interval values if they have been updated. Note
	// config_tls_refresh can only be set in the osquery flags so is
	// ignored here.
	saveHost := false

	distributedIntervalVal, ok := config.Options["distributed_interval"]
	distributedInterval, err := cast.ToUintE(distributedIntervalVal)
	if ok && err == nil && host.DistributedInterval != distributedInterval {
		host.DistributedInterval = distributedInterval
		saveHost = true
	}

	loggerTLSPeriodVal, ok := config.Options["logger_tls_period"]
	loggerTLSPeriod, err := cast.ToUintE(loggerTLSPeriodVal)
	if ok && err == nil && host.LoggerTLSPeriod != loggerTLSPeriod {
		host.LoggerTLSPeriod = loggerTLSPeriod
		saveHost = true
	}

	if saveHost {
		err := svc.ds.SaveHost(&host)
		if err != nil {
			return nil, err
		}
	}

	return svc.getFIMConfig(ctx, config)
}

// If osqueryWriters are based on bufio we want to flush after a batch of
// writes so log entry gets completely written to the logfile.
type flusher interface {
	Flush() error
}

func (svc service) SubmitStatusLogs(ctx context.Context, logs []json.RawMessage) error {
	for _, log := range logs {
		if _, err := svc.osqueryStatusLogWriter.Write(append(log, '\n')); err != nil {
			return osqueryError{message: "error writing status log: " + err.Error()}
		}
	}
	if writer, ok := svc.osqueryStatusLogWriter.(flusher); ok {
		err := writer.Flush()
		if err != nil {
			return osqueryError{message: "error flushing status log: " + err.Error()}
		}
	}
	return nil
}

func (svc service) SubmitResultLogs(ctx context.Context, logs []json.RawMessage) error {
	for _, log := range logs {
		if _, err := svc.osqueryResultLogWriter.Write(append(log, '\n')); err != nil {
			return osqueryError{message: "error writing result log: " + err.Error()}
		}
	}
	if writer, ok := svc.osqueryResultLogWriter.(flusher); ok {
		err := writer.Flush()
		if err != nil {
			return osqueryError{message: "error flushing status log: " + err.Error()}
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

// hostDistributedQueryPrefix is appended before the query name when a query is
// run from a distributed query campaign
const hostDistributedQueryPrefix = "kolide_distributed_query_"

// detailQueries defines the detail queries that should be run on the host, as
// well as how the results of those queries should be ingested into the
// kolide.Host data model. This map should not be modified at runtime.
var detailQueries = map[string]struct {
	Query      string
	IngestFunc func(logger log.Logger, host *kolide.Host, rows []map[string]string) error
}{
	"network_interface": {
		Query: `select * from interface_details id join interface_addresses ia
                        on ia.interface = id.interface where broadcast != ""
                        order by (ibytes + obytes) desc`,
		IngestFunc: func(logger log.Logger, host *kolide.Host, rows []map[string]string) (err error) {
			if len(rows) == 0 {
				logger.Log("component", "service", "method", "IngestFunc", "err",
					"detail_query_network_interface expected 1 or more results")
				return nil
			}
			networkInterfaces := []*kolide.NetworkInterface{}

			for _, row := range rows {
				nic := kolide.NetworkInterface{}

				nic.MAC = row["mac"]
				nic.IPAddress = row["address"]
				nic.Broadcast = row["broadcast"]
				if nic.IBytes, err = strconv.ParseInt(emptyToZero(row["ibytes"]), 10, 64); err != nil {
					return err
				}
				if nic.IErrors, err = strconv.ParseInt(emptyToZero(row["ierrors"]), 10, 64); err != nil {
					return err
				}
				nic.Interface = row["interface"]
				if nic.IPackets, err = strconv.ParseInt(emptyToZero(row["ipackets"]), 10, 64); err != nil {
					return err
				}
				// Optional last_change
				if lastChange, ok := row["last_change"]; ok {
					if nic.LastChange, err = strconv.ParseInt(emptyToZero(lastChange), 10, 64); err != nil {
						return err
					}
				}
				nic.Mask = row["mask"]
				if nic.Metric, err = strconv.Atoi(emptyToZero(row["metric"])); err != nil {
					return err
				}
				if nic.MTU, err = strconv.Atoi(emptyToZero(row["mtu"])); err != nil {
					return err
				}
				if nic.OBytes, err = strconv.ParseInt(emptyToZero(row["obytes"]), 10, 64); err != nil {
					return err
				}
				if nic.OErrors, err = strconv.ParseInt(emptyToZero(row["oerrors"]), 10, 64); err != nil {
					return err
				}
				if nic.OPackets, err = strconv.ParseInt(emptyToZero(row["opackets"]), 10, 64); err != nil {
					return err
				}
				nic.PointToPoint = row["point_to_point"]
				if nic.Type, err = strconv.Atoi(emptyToZero(row["type"])); err != nil {
					return err
				}
				networkInterfaces = append(networkInterfaces, &nic)
			}

			host.NetworkInterfaces = networkInterfaces

			return nil
		},
	},
	"os_version": {
		Query: "select * from os_version limit 1",
		IngestFunc: func(logger log.Logger, host *kolide.Host, rows []map[string]string) error {
			if len(rows) != 1 {
				logger.Log("component", "service", "method", "IngestFunc", "err",
					fmt.Sprintf("detail_query_os_version expected single result got %d", len(rows)))
				return nil
			}

			host.OSVersion = fmt.Sprintf(
				"%s %s.%s.%s",
				rows[0]["name"],
				rows[0]["major"],
				rows[0]["minor"],
				rows[0]["patch"],
			)
			host.OSVersion = strings.Trim(host.OSVersion, ".")

			if build, ok := rows[0]["build"]; ok {
				host.Build = build
			}

			host.Platform = rows[0]["platform"]
			host.PlatformLike = rows[0]["platform_like"]
			host.CodeName = rows[0]["code_name"]

			// On centos6 there is an osquery bug that leaves
			// platform empty. Here we workaround.
			if host.Platform == "" &&
				strings.Contains(strings.ToLower(rows[0]["name"]), "centos") {
				host.Platform = "centos"
			}

			return nil
		},
	},
	"osquery_flags": {
		// Collect the interval info (used for online status
		// calculation) from the osquery flags. We typically control
		// distributed_interval (but it's not required), and typically
		// do not control config_tls_refresh.
		Query: `select name, value from osquery_flags where name in ("distributed_interval", "config_tls_refresh", "config_refresh", "logger_tls_period")`,
		IngestFunc: func(logger log.Logger, host *kolide.Host, rows []map[string]string) error {
			var configTLSRefresh, configRefresh uint
			for _, row := range rows {
				switch row["name"] {

				case "distributed_interval":
					interval, err := strconv.Atoi(emptyToZero(row["value"]))
					if err != nil {
						return errors.Wrap(err, "parsing distributed_interval")
					}
					host.DistributedInterval = uint(interval)

				case "config_tls_refresh":
					// Prior to osquery 2.4.6, the flag was
					// called `config_tls_refresh`.
					interval, err := strconv.Atoi(emptyToZero(row["value"]))
					if err != nil {
						return errors.Wrap(err, "parsing config_tls_refresh")
					}
					configTLSRefresh = uint(interval)

				case "config_refresh":
					// After 2.4.6 `config_tls_refresh` was
					// aliased to `config_refresh`.
					interval, err := strconv.Atoi(emptyToZero(row["value"]))
					if err != nil {
						return errors.Wrap(err, "parsing config_refresh")
					}
					configRefresh = uint(interval)

				case "logger_tls_period":
					interval, err := strconv.Atoi(emptyToZero(row["value"]))
					if err != nil {
						return errors.Wrap(err, "parsing logger_tls_period")
					}
					host.LoggerTLSPeriod = uint(interval)
				}
			}

			// Since the `config_refresh` flag existed prior to
			// 2.4.6 and had a different meaning, we prefer
			// `config_tls_refresh` if it was set, and use
			// `config_refresh` as a fallback.
			if configTLSRefresh != 0 {
				host.ConfigTLSRefresh = configTLSRefresh
			} else {
				host.ConfigTLSRefresh = configRefresh
			}

			return nil
		},
	},
	"osquery_info": {
		Query: "select * from osquery_info limit 1",
		IngestFunc: func(logger log.Logger, host *kolide.Host, rows []map[string]string) error {
			if len(rows) != 1 {
				logger.Log("component", "service", "method", "IngestFunc", "err",
					fmt.Sprintf("detail_query_osquery_info expected single result got %d", len(rows)))
				return nil
			}

			host.OsqueryVersion = rows[0]["version"]

			return nil
		},
	},
	"system_info": {
		Query: "select * from system_info limit 1",
		IngestFunc: func(logger log.Logger, host *kolide.Host, rows []map[string]string) error {
			if len(rows) != 1 {
				logger.Log("component", "service", "method", "IngestFunc", "err",
					fmt.Sprintf("detail_query_system_info expected single result got %d", len(rows)))
				return nil
			}

			var err error
			host.PhysicalMemory, err = strconv.Atoi(emptyToZero(rows[0]["physical_memory"]))
			if err != nil {
				return err
			}
			host.HostName = rows[0]["hostname"]
			host.UUID = rows[0]["uuid"]
			host.CPUType = rows[0]["cpu_type"]
			host.CPUSubtype = rows[0]["cpu_subtype"]
			host.CPUBrand = rows[0]["cpu_brand"]
			host.CPUPhysicalCores, err = strconv.Atoi(emptyToZero(rows[0]["cpu_physical_cores"]))
			if err != nil {
				return err
			}
			host.CPULogicalCores, err = strconv.Atoi(emptyToZero(rows[0]["cpu_logical_cores"]))
			if err != nil {
				return err
			}
			host.HardwareVendor = rows[0]["hardware_vendor"]
			host.HardwareModel = rows[0]["hardware_model"]
			host.HardwareVersion = rows[0]["hardware_version"]
			host.HardwareSerial = rows[0]["hardware_serial"]
			host.ComputerName = rows[0]["computer_name"]
			return nil
		},
	},
	"uptime": {
		Query: "select * from uptime limit 1",
		IngestFunc: func(logger log.Logger, host *kolide.Host, rows []map[string]string) error {
			if len(rows) != 1 {
				logger.Log("component", "service", "method", "IngestFunc", "err",
					fmt.Sprintf("detail_query_uptime expected single result got %d", len(rows)))
				return nil
			}

			uptimeSeconds, err := strconv.Atoi(emptyToZero(rows[0]["total_seconds"]))
			if err != nil {
				return err
			}
			host.Uptime = time.Duration(uptimeSeconds) * time.Second

			return nil
		},
	},
}

// detailUpdateInterval determines how often the detail queries should be
// updated
const detailUpdateInterval = 1 * time.Hour

// hostDetailQueries returns the map of queries that should be executed by
// osqueryd to fill in the host details
func (svc service) hostDetailQueries(host kolide.Host) map[string]string {
	queries := make(map[string]string)
	if host.DetailUpdateTime.After(svc.clock.Now().Add(-detailUpdateInterval)) {
		// No need to update already fresh details
		return queries
	}

	for name, query := range detailQueries {
		queries[hostDetailQueryPrefix+name] = query.Query
	}
	return queries
}

func (svc service) GetDistributedQueries(ctx context.Context) (map[string]string, uint, error) {
	host, ok := hostctx.FromContext(ctx)
	if !ok {
		return nil, 0, osqueryError{message: "internal error: missing host from request context"}
	}

	queries := svc.hostDetailQueries(host)

	// Retrieve the label queries that should be updated
	cutoff := svc.clock.Now().Add(-svc.config.Osquery.LabelUpdateInterval)
	labelQueries, err := svc.ds.LabelQueriesForHost(&host, cutoff)
	if err != nil {
		return nil, 0, osqueryError{message: "retrieving label queries: " + err.Error()}
	}

	for name, query := range labelQueries {
		queries[hostLabelQueryPrefix+name] = query
	}

	distributedQueries, err := svc.ds.DistributedQueriesForHost(&host)
	if err != nil {
		return nil, 0, osqueryError{message: "retrieving query campaigns: " + err.Error()}
	}

	for id, query := range distributedQueries {
		queries[hostDistributedQueryPrefix+strconv.Itoa(int(id))] = query
	}

	accelerate := uint(0)
	if host.HostName == "" && host.Platform == "" {
		// Assume this host is just enrolling, and accelerate checkins
		// (to allow for platform restricted labels to run quickly
		// after platform is retrieved from details)
		accelerate = 10
	}

	return queries, accelerate, nil
}

// ingestDetailQuery takes the results of a detail query and modifies the
// provided kolide.Host appropriately.
func (svc service) ingestDetailQuery(host *kolide.Host, name string, rows []map[string]string) error {
	trimmedQuery := strings.TrimPrefix(name, hostDetailQueryPrefix)
	query, ok := detailQueries[trimmedQuery]
	if !ok {
		return osqueryError{message: "unknown detail query " + trimmedQuery}
	}

	err := query.IngestFunc(svc.logger, host, rows)
	if err != nil {
		return osqueryError{
			message: fmt.Sprintf("ingesting query %s: %s", name, err.Error()),
		}
	}

	return nil
}

// ingestLabelQuery records the results of label queries run by a host
func (svc service) ingestLabelQuery(host kolide.Host, query string, rows []map[string]string, results map[uint]bool) error {
	trimmedQuery := strings.TrimPrefix(query, hostLabelQueryPrefix)
	trimmedQueryNum, err := strconv.Atoi(emptyToZero(trimmedQuery))
	if err != nil {
		return errors.Wrap(err, "converting query from string to int")
	}
	// A label query matches if there is at least one result for that
	// query. We must also store negative results.
	results[uint(trimmedQueryNum)] = len(rows) > 0
	return nil
}

// ingestDistributedQuery takes the results of a distributed query and modifies the
// provided kolide.Host appropriately.
func (svc service) ingestDistributedQuery(host kolide.Host, name string, rows []map[string]string, failed bool) error {
	trimmedQuery := strings.TrimPrefix(name, hostDistributedQueryPrefix)

	campaignID, err := strconv.Atoi(emptyToZero(trimmedQuery))
	if err != nil {
		return osqueryError{message: "unable to parse campaign ID: " + trimmedQuery}
	}

	// Write the results to the pubsub store
	res := kolide.DistributedQueryResult{
		DistributedQueryCampaignID: uint(campaignID),
		Host: host,
		Rows: rows,
	}
	if failed {
		// osquery errors are not currently helpful, but we should fix
		// them to be better in the future
		errString := "failed"
		res.Error = &errString
	}

	err = svc.resultStore.WriteResult(res)
	if err != nil {
		nErr, ok := err.(pubsub.Error)
		if !ok || !nErr.NoSubscriber() {
			return osqueryError{message: "writing results: " + err.Error()}
		}

		// If there are no subscribers, the campaign is "orphaned"
		// and should be closed so that we don't continue trying to
		// execute that query when we can't write to any subscriber
		campaign, err := svc.ds.DistributedQueryCampaign(uint(campaignID))
		if err != nil {
			return osqueryError{message: "loading orphaned campaign: " + err.Error()}
		}

		campaign.Status = kolide.QueryComplete
		if err := svc.ds.SaveDistributedQueryCampaign(campaign); err != nil {
			return osqueryError{message: "closing orphaned campaign: " + err.Error()}
		}
	}

	// Record execution of the query
	status := kolide.ExecutionSucceeded
	if failed {
		status = kolide.ExecutionFailed
	}
	exec := &kolide.DistributedQueryExecution{
		HostID: host.ID,
		DistributedQueryCampaignID: uint(campaignID),
		Status: status,
	}

	_, err = svc.ds.NewDistributedQueryExecution(exec)
	if err != nil {
		return osqueryError{message: "recording execution: " + err.Error()}
	}

	return nil
}

func (svc service) SubmitDistributedQueryResults(ctx context.Context, results kolide.OsqueryDistributedQueryResults, statuses map[string]string) error {
	host, ok := hostctx.FromContext(ctx)

	if !ok {
		return osqueryError{message: "internal error: missing host from request context"}
	}

	var err error
	detailUpdated := false
	labelResults := map[uint]bool{}
	for query, rows := range results {
		switch {
		case strings.HasPrefix(query, hostDetailQueryPrefix):
			err = svc.ingestDetailQuery(&host, query, rows)
			detailUpdated = true
		case strings.HasPrefix(query, hostLabelQueryPrefix):
			err = svc.ingestLabelQuery(host, query, rows, labelResults)
		case strings.HasPrefix(query, hostDistributedQueryPrefix):
			// osquery docs say any nonzero (string) value for
			// status indicates a query error
			status, ok := statuses[query]
			failed := ok && status != "0"
			err = svc.ingestDistributedQuery(host, query, rows, failed)
		default:
			err = osqueryError{message: "unknown query prefix: " + query}
		}

		if err != nil {
			return osqueryError{message: "failed to ingest result: " + err.Error()}
		}

	}

	if len(labelResults) > 0 {
		err = svc.ds.RecordLabelQueryExecutions(&host, labelResults, svc.clock.Now())
		if err != nil {
			return osqueryError{message: "failed to save labels: " + err.Error()}
		}
	}

	if detailUpdated {
		host.DetailUpdateTime = svc.clock.Now()
	}

	if len(labelResults) > 0 || detailUpdated {
		err = svc.ds.SaveHost(&host)
		if err != nil {
			return osqueryError{message: "failed to update host details: " + err.Error()}
		}
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////
// Endpoints
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

type getClientConfigRequest struct {
	NodeKey string `json:"node_key"`
}

type getClientConfigResponse struct {
	kolide.OsqueryConfig
	Err error `json:"error,omitempty"`
}

func (r getClientConfigResponse) error() error { return r.Err }

func makeGetClientConfigEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		config, err := svc.GetClientConfig(ctx)
		if err != nil {
			return getClientConfigResponse{Err: err}, nil
		}
		return getClientConfigResponse{OsqueryConfig: *config}, nil
	}
}

type getDistributedQueriesRequest struct {
	NodeKey string `json:"node_key"`
}

type getDistributedQueriesResponse struct {
	Queries    map[string]string `json:"queries"`
	Accelerate uint              `json:"accelerate,omitempty"`
	Err        error             `json:"error,omitempty"`
}

func (r getDistributedQueriesResponse) error() error { return r.Err }

func makeGetDistributedQueriesEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		queries, accelerate, err := svc.GetDistributedQueries(ctx)
		if err != nil {
			return getDistributedQueriesResponse{Err: err}, nil
		}
		return getDistributedQueriesResponse{Queries: queries, Accelerate: accelerate}, nil
	}
}

type submitDistributedQueryResultsRequest struct {
	NodeKey  string                                `json:"node_key"`
	Results  kolide.OsqueryDistributedQueryResults `json:"queries"`
	Statuses map[string]string                     `json:"statuses"`
}

type submitDistributedQueryResultsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r submitDistributedQueryResultsResponse) error() error { return r.Err }

func makeSubmitDistributedQueryResultsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitDistributedQueryResultsRequest)
		err := svc.SubmitDistributedQueryResults(ctx, req.Results, req.Statuses)
		if err != nil {
			return submitDistributedQueryResultsResponse{Err: err}, nil
		}
		return submitDistributedQueryResultsResponse{}, nil
	}
}

type submitLogsRequest struct {
	NodeKey string          `json:"node_key"`
	LogType string          `json:"log_type"`
	Data    json.RawMessage `json:"data"`
}

type submitLogsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r submitLogsResponse) error() error { return r.Err }

func makeSubmitLogsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(submitLogsRequest)

		var err error
		switch req.LogType {
		case "status":
			var statuses []json.RawMessage
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

////////////////////////////////////////////////////////////////////////////////
// Logging
////////////////////////////////////////////////////////////////////////////////

func (mw loggingMiddleware) EnrollAgent(ctx context.Context, enrollSecret string, hostIdentifier string) (string, error) {
	var (
		nodeKey string
		err     error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "EnrollAgent",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	nodeKey, err = mw.Service.EnrollAgent(ctx, enrollSecret, hostIdentifier)
	return nodeKey, err
}

func (mw loggingMiddleware) AuthenticateHost(ctx context.Context, nodeKey string) (*kolide.Host, error) {
	var (
		host *kolide.Host
		err  error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "AuthenticateHost",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	host, err = mw.Service.AuthenticateHost(ctx, nodeKey)
	return host, err
}

func (mw loggingMiddleware) GetClientConfig(ctx context.Context) (*kolide.OsqueryConfig, error) {
	var (
		config *kolide.OsqueryConfig
		err    error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetClientConfig",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	config, err = mw.Service.GetClientConfig(ctx)
	return config, err
}

func (mw loggingMiddleware) GetDistributedQueries(ctx context.Context) (map[string]string, uint, error) {
	var (
		queries    map[string]string
		err        error
		accelerate uint
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetDistributedQueries",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	queries, accelerate, err = mw.Service.GetDistributedQueries(ctx)
	return queries, accelerate, err
}

func (mw loggingMiddleware) SubmitDistributedQueryResults(ctx context.Context, results kolide.OsqueryDistributedQueryResults, statuses map[string]string) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SubmitDistributedQueryResults",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.SubmitDistributedQueryResults(ctx, results, statuses)
	return err
}

func (mw loggingMiddleware) SubmitStatusLogs(ctx context.Context, logs []json.RawMessage) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SubmitStatusLogs",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.SubmitStatusLogs(ctx, logs)
	return err
}

func (mw loggingMiddleware) SubmitResultLogs(ctx context.Context, logs []json.RawMessage) error {
	var (
		err error
	)

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "SubmitResultLogs",
			"ip_addr", ctx.Value(kithttp.ContextKeyRequestRemoteAddr).(string),
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.SubmitResultLogs(ctx, logs)
	return err
}
