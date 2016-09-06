package kolide

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"
)

type OsqueryStore interface {
	LabelQueriesForHost(host *Host, cutoff time.Time) (map[string]string, error)
	RecordLabelQueryExecutions(host *Host, results map[string]bool, t time.Time) error
	NewLabel(label *Label) error
}

type OsqueryService interface {
	EnrollAgent(ctx context.Context, enrollSecret, hostIdentifier string) (string, error)
	GetClientConfig(ctx context.Context, action string, data *json.RawMessage) (*OsqueryConfig, error)
	Log(ctx context.Context, logType string, data *json.RawMessage) error
	GetDistributedQueries(ctx context.Context) (map[string]string, error)
	LogDistributedQueryResults(ctx context.Context, queries map[string][]map[string]string) error
}

type OsqueryConfig struct {
	Packs    []Pack
	Schedule []Query
	Options  map[string]interface{}
}

// TODO: move this to just use LabelQueriesForHot
// LabelQueriesForHost calculates the appropriate update cutoff (given
// interval) and uses the datastore to retrieve the label queries for the
// provided host.
func LabelQueriesForHost(store OsqueryStore, host *Host, interval time.Duration) (map[string]string, error) {
	cutoff := time.Now().Add(-interval)
	return store.LabelQueriesForHost(host, cutoff)
}
