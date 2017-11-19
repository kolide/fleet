package tlsremote

import (
	"context"
	"testing"

	"github.com/WatchBeam/clock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/datastore/inmem"
	"github.com/kolide/fleet/server/kolide"
)

// TestGetNodeKey tests the reflection logic for pulling the node key from
// various (fake) request types
func TestGetNodeKey(t *testing.T) {
	type Foo struct {
		Foo     string
		NodeKey string
	}

	type Bar struct {
		Bar     string
		NodeKey string
	}

	type Nope struct {
		Nope string
	}

	type Almost struct {
		NodeKey int
	}

	var getNodeKeyTests = []struct {
		i         interface{}
		expectKey string
		shouldErr bool
	}{
		{
			i:         Foo{Foo: "foo", NodeKey: "fookey"},
			expectKey: "fookey",
			shouldErr: false,
		},
		{
			i:         Bar{Bar: "bar", NodeKey: "barkey"},
			expectKey: "barkey",
			shouldErr: false,
		},
		{
			i:         Nope{Nope: "nope"},
			expectKey: "",
			shouldErr: true,
		},
		{
			i:         Almost{NodeKey: 10},
			expectKey: "",
			shouldErr: true,
		},
	}

	for _, tt := range getNodeKeyTests {
		t.Run("", func(t *testing.T) {
			key, err := getNodeKey(tt.i)
			assert.Equal(t, tt.expectKey, key)
			if tt.shouldErr {
				assert.IsType(t, osqueryError{}, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestAuthenticatedHost(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)

	_, err = ds.NewAppConfig(&kolide.AppConfig{EnrollSecret: "foobarbaz"})
	require.Nil(t, err)

	svc := newTestService(t, ds, nil)
	require.Nil(t, err)

	endpoint := authenticatedHost(
		svc,
		func(ctx context.Context, request interface{}) (interface{}, error) {
			return nil, nil
		},
	)

	ctx := context.Background()
	goodNodeKey, err := svc.EnrollAgent(ctx, "foobarbaz", "host123")
	require.Nil(t, err)
	require.NotEmpty(t, goodNodeKey)

	var authenticatedHostTests = []struct {
		nodeKey   string
		shouldErr bool
	}{
		{
			nodeKey:   "invalid",
			shouldErr: true,
		},
		{
			nodeKey:   "",
			shouldErr: true,
		},
		{
			nodeKey:   goodNodeKey,
			shouldErr: false,
		},
	}

	for _, tt := range authenticatedHostTests {
		t.Run("", func(t *testing.T) {
			var r = struct{ NodeKey string }{NodeKey: tt.nodeKey}
			_, err = endpoint(context.Background(), r)
			if tt.shouldErr {
				assert.IsType(t, osqueryError{}, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
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
