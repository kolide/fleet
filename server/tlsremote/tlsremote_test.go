package tlsremote

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/kolide/fleet/server/config"
	hostctx "github.com/kolide/fleet/server/contexts/host"
	"github.com/kolide/fleet/server/datastore/inmem"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnrollAgent(t *testing.T) {
	ds, svc, _ := setupOsqueryTests(t)
	ctx := context.Background()

	hosts, err := ds.ListHosts(kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, hosts, 0)

	nodeKey, err := svc.EnrollAgent(ctx, "", "host123")
	require.Nil(t, err)
	assert.NotEmpty(t, nodeKey)

	hosts, err = ds.ListHosts(kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, hosts, 1)
}

func TestEnrollAgentIncorrectEnrollSecret(t *testing.T) {
	ds, svc, _ := setupOsqueryTests(t)
	ctx := context.Background()

	hosts, err := ds.ListHosts(kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, hosts, 0)

	nodeKey, err := svc.EnrollAgent(ctx, "not_correct", "host123")
	assert.NotNil(t, err)
	assert.Empty(t, nodeKey)

	hosts, err = ds.ListHosts(kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, hosts, 0)
}

func TestAuthenticateHost(t *testing.T) {
	ds, svc, mockClock := setupOsqueryTests(t)
	ctx := context.Background()

	nodeKey, err := svc.EnrollAgent(ctx, "", "host123")
	require.Nil(t, err)

	mockClock.AddTime(1 * time.Minute)

	host, err := svc.AuthenticateHost(ctx, nodeKey)
	require.Nil(t, err)

	// Verify that the update time is set appropriately
	checkHost, err := ds.Host(host.ID)
	require.Nil(t, err)
	assert.Equal(t, mockClock.Now(), checkHost.UpdatedAt)

	// Advance clock time and check that seen time is updated
	mockClock.AddTime(1*time.Minute + 27*time.Second)

	_, err = svc.AuthenticateHost(ctx, nodeKey)
	require.Nil(t, err)

	checkHost, err = ds.Host(host.ID)
	require.Nil(t, err)
	assert.Equal(t, mockClock.Now().UTC(), checkHost.UpdatedAt.UTC())
}

type nopCloserWriter struct {
	io.Writer
}

func (n *nopCloserWriter) Close() error { return nil }

func TestSubmitStatusLogs(t *testing.T) {
	ds, svc, _ := setupOsqueryTests(t)
	ctx := context.Background()

	_, err := svc.EnrollAgent(ctx, "", "host123")
	require.Nil(t, err)

	hosts, err := ds.ListHosts(kolide.ListOptions{})
	require.Nil(t, err)
	require.Len(t, hosts, 1)
	host := hosts[0]
	ctx = hostctx.NewContext(ctx, *host)

	var statusBuf bytes.Buffer
	svc.osqueryStatusLogWriter = &nopCloserWriter{&statusBuf}

	logs := []string{
		`{"severity":"0","filename":"tls.cpp","line":"216","message":"some message","version":"1.8.2","decorations":{"host_uuid":"uuid_foobar","username":"zwass"}}`,
		`{"severity":"1","filename":"buffered.cpp","line":"122","message":"warning!","version":"1.8.2","decorations":{"host_uuid":"uuid_foobar","username":"zwass"}}`,
	}
	logJSON := fmt.Sprintf("[%s]", strings.Join(logs, ","))

	var status []kolide.OsqueryStatusLog
	err = json.Unmarshal([]byte(logJSON), &status)
	require.Nil(t, err)

	err = svc.SubmitStatusLogs(ctx, status)
	require.Nil(t, err)

	statusJSON := statusBuf.String()
	statusJSON = strings.TrimRight(statusJSON, "\n")
	statusLines := strings.Split(statusJSON, "\n")

	if assert.Equal(t, len(logs), len(statusLines)) {
		for i, line := range statusLines {
			assert.JSONEq(t, logs[i], line)
		}
	}
}

func TestSubmitResultLogs(t *testing.T) {
	ds, svc, _ := setupOsqueryTests(t)
	ctx := context.Background()

	_, err := svc.EnrollAgent(ctx, "", "host123")
	require.Nil(t, err)

	hosts, err := ds.ListHosts(kolide.ListOptions{})
	require.Nil(t, err)
	require.Len(t, hosts, 1)
	host := hosts[0]
	ctx = hostctx.NewContext(ctx, *host)

	var resultBuf bytes.Buffer
	svc.osqueryResultLogWriter = &nopCloserWriter{&resultBuf}

	logs := []string{
		`{"name":"system_info","hostIdentifier":"some_uuid","calendarTime":"Fri Sep 30 17:55:15 2016 UTC","unixTime":"1475258115","decorations":{"host_uuid":"some_uuid","username":"zwass"},"columns":{"cpu_brand":"Intel(R) Core(TM) i7-4770HQ CPU @ 2.20GHz","hostname":"hostimus","physical_memory":"17179869184"},"action":"added"}`,
		`{"name":"encrypted","hostIdentifier":"some_uuid","calendarTime":"Fri Sep 30 21:19:15 2016 UTC","unixTime":"1475270355","decorations":{"host_uuid":"4740D59F-699E-5B29-960B-979AAF9BBEEB","username":"zwass"},"columns":{"encrypted":"1","name":"\/dev\/disk1","type":"AES-XTS","uid":"","user_uuid":"","uuid":"some_uuid"},"action":"added"}`,
		`{"snapshot":[{"hour":"20","minutes":"8"}],"action":"snapshot","name":"time","hostIdentifier":"1379f59d98f4","calendarTime":"Tue Jan 10 20:08:51 2017 UTC","unixTime":"1484078931","decorations":{"host_uuid":"EB714C9D-C1F8-A436-B6DA-3F853C5502EA"}}`,
		`{"diffResults":{"removed":[{"address":"127.0.0.1","hostnames":"kl.groob.io"}],"added":""},"name":"pack\/test\/hosts","hostIdentifier":"FA01680E-98CA-5557-8F59-7716ECFEE964","calendarTime":"Sun Nov 19 00:02:08 2017 UTC","unixTime":"1511049728","epoch":"0","counter":"10","decorations":{"host_uuid":"FA01680E-98CA-5557-8F59-7716ECFEE964","hostname":"kl.groob.io"}}`,
		// fleet will accept anything in the "data" field of an log request.
		`{"unknown":{"foo": [] }}`,
	}
	logJSON := fmt.Sprintf("[%s]", strings.Join(logs, ","))

	var results []json.RawMessage
	err = json.Unmarshal([]byte(logJSON), &results)
	require.Nil(t, err)

	err = svc.SubmitResultLogs(ctx, results)
	require.Nil(t, err)

	resultJSON := resultBuf.String()
	resultJSON = strings.TrimRight(resultJSON, "\n")
	resultLines := strings.Split(resultJSON, "\n")

	if assert.Equal(t, len(logs), len(resultLines)) {
		for i, line := range resultLines {
			assert.JSONEq(t, logs[i], line)
		}
	}
}

type packService struct {
}

func (svc *packService) ListPacksForHost(ctx context.Context, hid uint) (packs []*kolide.Pack, err error) {
	return
}

type fimService struct{}

func (svc *fimService) GetFIM(ctx context.Context) (*kolide.FIMConfig, error) {
	return new(kolide.FIMConfig), nil
}

func TestGetClientConfig(t *testing.T) {
	ds, svc, mockClock := setupOsqueryTests(t)
	svc.packs = new(packService)
	svc.fim = new(fimService)

	ctx := context.Background()

	hosts, err := ds.ListHosts(kolide.ListOptions{})
	require.Nil(t, err)
	require.Len(t, hosts, 0)

	_, err = svc.EnrollAgent(ctx, "", "user.local")
	assert.Nil(t, err)

	hosts, err = ds.ListHosts(kolide.ListOptions{})
	require.Nil(t, err)
	require.Len(t, hosts, 1)
	host := hosts[0]

	ctx = hostctx.NewContext(ctx, *host)

	// with no queries, packs, labels, etc. verify the state of a fresh host
	// asking for a config
	config, err := svc.GetClientConfig(ctx)
	require.Nil(t, err)
	assert.NotNil(t, config)
	val, ok := config.Options["disable_distributed"]
	require.True(t, ok)
	disabled, ok := val.(bool)
	require.True(t, ok)
	assert.False(t, disabled)
	val, ok = config.Options["pack_delimiter"]
	require.True(t, ok)
	delim, ok := val.(string)
	require.True(t, ok)
	assert.Equal(t, "/", delim)

	// this will be greater than 0 if we ever start inserting an administration
	// pack
	assert.Len(t, config.Packs, 0)

	// let's populate the database with some info

	infoQuery := &kolide.Query{
		Name:  "Info",
		Query: "select * from osquery_info;",
	}
	infoQueryInterval := uint(60)
	infoQuery, err = ds.NewQuery(infoQuery)
	assert.Nil(t, err)

	monitoringPack := &kolide.Pack{
		Name: "monitoring",
	}
	_, err = ds.NewPack(monitoringPack)
	assert.Nil(t, err)

	test.NewScheduledQuery(t, ds, monitoringPack.ID, infoQuery.ID, infoQueryInterval, false, false)

	mysqlLabel := &kolide.Label{
		Name:  "MySQL Monitoring",
		Query: "select pid from processes where name = 'mysqld';",
	}
	mysqlLabel, err = ds.NewLabel(mysqlLabel)
	assert.Nil(t, err)

	err = ds.AddLabelToPack(mysqlLabel.ID, monitoringPack.ID)
	assert.Nil(t, err)

	err = ds.RecordLabelQueryExecutions(
		host,
		map[uint]bool{mysqlLabel.ID: true},
		mockClock.Now(),
	)
	assert.Nil(t, err)

	// with a minimal setup of packs, labels, and queries, will our host get the
	// pack
	config, err = svc.GetClientConfig(ctx)
	require.Nil(t, err)
	assert.Len(t, config.Packs, 1)
	assert.Len(t, config.Packs["monitoring"].Queries, 1)
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

func setupOsqueryTests(t *testing.T) (kolide.Datastore, *OsqueryService, *clock.MockClock) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)

	_, err = ds.NewAppConfig(&kolide.AppConfig{EnrollSecret: ""})
	require.Nil(t, err)

	mockClock := clock.NewMockClock()
	svc := newTestService(t, ds, nil)
	svc.clock = mockClock

	return ds, svc, mockClock
}
