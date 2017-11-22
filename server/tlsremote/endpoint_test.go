package tlsremote

import (
	"context"
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/go-kit/kit/endpoint"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

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

// helpers for the TestAuthenticatedHostMiddleware test.
type (
	authenticationTestService struct {
		ds *authenticationTestDatastore
		kolide.OsqueryService
	}

	authenticationTestDatastore struct {
		AuthenticateHostInvoked bool
		clock                   clock.Clock
		nodeKey                 string
		Datastore
	}
)

func (ds *authenticationTestDatastore) AuthenticateHost(nodeKey string) (*kolide.Host, error) {
	ds.AuthenticateHostInvoked = true
	if nodeKey == ds.nodeKey {
		return &kolide.Host{NodeKey: ds.nodeKey}, nil
	}
	return nil, errors.New("test: bad node key")
}

func (ds *authenticationTestDatastore) MarkHostSeen(host *kolide.Host, t time.Time) error {
	host.SeenTime = ds.clock.Now()
	return nil
}

func setupAuthenticatedHostTest(t *testing.T) (*authenticationTestService, endpoint.Endpoint) {
	goodNodeKey := "good-node-key"
	clock := clock.NewMockClock()

	// create db with with all the necessary dependencies
	ds := &authenticationTestDatastore{
		nodeKey: goodNodeKey,
		clock:   clock,
	}

	// create service using teh datastore
	svc := &authenticationTestService{
		ds: ds,
		OsqueryService: &OsqueryService{
			ds:    ds,
			clock: clock,
		},
	}

	// setup endpoint with authenticateHost middleware
	endpoint := authenticatedHost(
		svc,
		func(ctx context.Context, request interface{}) (interface{}, error) {
			return nil, nil
		},
	)

	return svc, endpoint
}

func TestAuthenticatedHostMiddleware(t *testing.T) {
	svc, endpoint := setupAuthenticatedHostTest(t)

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
			nodeKey:   svc.ds.nodeKey,
			shouldErr: false,
		},
	}

	for _, tt := range authenticatedHostTests {
		t.Run("", func(t *testing.T) {
			var r = struct{ NodeKey string }{NodeKey: tt.nodeKey}
			_, err := endpoint(context.Background(), r)
			if tt.shouldErr {
				assert.IsType(t, osqueryError{}, err)
			} else {
				assert.Nil(t, err)
			}
			assert.True(t, svc.ds.AuthenticateHostInvoked)
		})
	}
}
