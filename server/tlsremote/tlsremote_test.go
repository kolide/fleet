package tlsremote

import (
	"testing"

	"github.com/WatchBeam/clock"
	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/kolide"
)

func TestEnrollAgent(t *testing.T) {
}

func newTestService(t *testing.T, ds kolide.Datastore, rs kolide.QueryResultStore) *OsqueryService {
	cfg := config.TestConfig()
	svc := &OsqueryService{
		ds:          ds,
		resultStore: rs,
		clock:       clock.C,
		nodeKeySize: cfg.Osquery.NodeKeySize,
	}
	return svc
}
