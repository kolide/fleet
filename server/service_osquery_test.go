package server

import (
	"context"
	"testing"

	"github.com/WatchBeam/clock"
	"github.com/kolide/kolide-ose/datastore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnrollAgent(t *testing.T) {
	ds, err := datastore.New("gorm-sqlite3", ":memory:")
	assert.Nil(t, err)

	svc, err := newTestService(ds)
	assert.Nil(t, err)

	ctx := context.Background()

	hosts, err := ds.Hosts()
	assert.Nil(t, err)
	assert.Len(t, hosts, 0)

	nodeKey, err := svc.EnrollAgent(ctx, "", "host123")
	assert.Nil(t, err)
	assert.NotEmpty(t, nodeKey)

	hosts, err = ds.Hosts()
	assert.Nil(t, err)
	assert.Len(t, hosts, 1)
}

func TestEnrollAgentIncorrectEnrollSecret(t *testing.T) {
	ds, err := datastore.New("gorm-sqlite3", ":memory:")
	assert.Nil(t, err)

	svc, err := newTestService(ds)
	assert.Nil(t, err)

	ctx := context.Background()

	hosts, err := ds.Hosts()
	assert.Nil(t, err)
	assert.Len(t, hosts, 0)

	nodeKey, err := svc.EnrollAgent(ctx, "not_correct", "host123")
	assert.NotNil(t, err)
	assert.Empty(t, nodeKey)

	hosts, err = ds.Hosts()
	assert.Nil(t, err)
	assert.Len(t, hosts, 0)
}

func TestGetDistributedQueries(t *testing.T) {
	ds, err := datastore.New("gorm-sqlite3", ":memory:")
	assert.Nil(t, err)

	mockClock := clock.NewMockClock()

	svc, err := newTestServiceWithClock(ds, mockClock)
	assert.Nil(t, err)

	ctx := context.Background()

	_, err = svc.EnrollAgent(ctx, "", "host123")
	assert.Nil(t, err)

	hosts, err := ds.Hosts()
	require.Nil(t, err)
	require.Len(t, hosts, 1)
	host := hosts[0]

	ctx = context.WithValue(ctx, "osqueryHost", host)

	// With no platform set, we should get the details query
	queries, err := svc.GetDistributedQueries(ctx)
	assert.Nil(t, err)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries, "kolide_detail_query_platform")
	assert.Equal(t,
		"select build_platform from osquery_info;",
		queries["kolide_detail_query_platform"],
	)

	host.Platform = "darwin"

	ctx = context.WithValue(ctx, "osqueryHost", host)

	// With the platform set, we should get the label queries (but none
	// exist yet)
	queries, err = svc.GetDistributedQueries(ctx)
	assert.Nil(t, err)
	assert.Len(t, queries, 0)
}
