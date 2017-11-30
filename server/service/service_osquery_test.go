package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/kolide/fleet/server/config"
	hostctx "github.com/kolide/fleet/server/contexts/host"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/datastore/inmem"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/mock"
	"github.com/kolide/fleet/server/pubsub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrphanedQueryCampaign(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)

	_, err = ds.NewAppConfig(&kolide.AppConfig{EnrollSecret: ""})
	require.Nil(t, err)

	rs := pubsub.NewInmemQueryResults()

	svc, err := newTestService(ds, rs)
	require.Nil(t, err)

	ctx := context.Background()

	nodeKey, err := svc.EnrollAgent(ctx, "", "host123")
	require.Nil(t, err)

	host, err := ds.AuthenticateHost(nodeKey)
	require.Nil(t, err)

	ctx = viewer.NewContext(context.Background(), viewer.Viewer{
		User: &kolide.User{
			ID: 0,
		},
	})
	q := "select year, month, day, hour, minutes, seconds from time"
	campaign, err := svc.NewDistributedQueryCampaign(ctx, q, []uint{}, []uint{})
	require.Nil(t, err)

	campaign.Status = kolide.QueryRunning
	err = ds.SaveDistributedQueryCampaign(campaign)
	require.Nil(t, err)

	queryKey := fmt.Sprintf("%s%d", hostDistributedQueryPrefix, campaign.ID)

	expectedRows := []map[string]string{
		{
			"year":    "2016",
			"month":   "11",
			"day":     "11",
			"hour":    "6",
			"minutes": "12",
			"seconds": "10",
		},
	}
	results := map[string][]map[string]string{
		queryKey: expectedRows,
	}

	// Submit results
	ctx = hostctx.NewContext(context.Background(), *host)
	err = svc.SubmitDistributedQueryResults(ctx, results, map[string]string{})
	require.Nil(t, err)

	// The campaign should be set to completed because it is orphaned
	campaign, err = ds.DistributedQueryCampaign(campaign.ID)
	require.Nil(t, err)
	assert.Equal(t, kolide.QueryComplete, campaign.Status)
}

func TestUpdateHostIntervals(t *testing.T) {
	ds := new(mock.Store)

	svc, err := newTestService(ds, nil)
	require.Nil(t, err)

	ds.ListDecoratorsFunc = func(opt ...kolide.OptionalArg) ([]*kolide.Decorator, error) {
		return []*kolide.Decorator{}, nil
	}
	ds.ListPacksFunc = func(opt kolide.ListOptions) ([]*kolide.Pack, error) {
		return []*kolide.Pack{}, nil
	}
	ds.ListLabelsForHostFunc = func(hid uint) ([]kolide.Label, error) {
		return []kolide.Label{}, nil
	}
	ds.AppConfigFunc = func() (*kolide.AppConfig, error) {
		return &kolide.AppConfig{FIMInterval: 400}, nil
	}
	ds.FIMSectionsFunc = func() (kolide.FIMSections, error) {
		sections := kolide.FIMSections{
			"etc": []string{
				"/etc/%%",
			},
		}
		return sections, nil
	}

	var testCases = []struct {
		initHost       kolide.Host
		finalHost      kolide.Host
		configOptions  map[string]interface{}
		saveHostCalled bool
	}{
		// Both updated
		{
			kolide.Host{
				ConfigTLSRefresh: 60,
			},
			kolide.Host{
				DistributedInterval: 11,
				LoggerTLSPeriod:     33,
				ConfigTLSRefresh:    60,
			},
			map[string]interface{}{
				"distributed_interval": 11,
				"logger_tls_period":    33,
				"logger_plugin":        "tls",
			},
			true,
		},
		// Only logger_tls_period updated
		{
			kolide.Host{
				DistributedInterval: 11,
				ConfigTLSRefresh:    60,
			},
			kolide.Host{
				DistributedInterval: 11,
				LoggerTLSPeriod:     33,
				ConfigTLSRefresh:    60,
			},
			map[string]interface{}{
				"distributed_interval": 11,
				"logger_tls_period":    33,
			},
			true,
		},
		// Only distributed_interval updated
		{
			kolide.Host{
				ConfigTLSRefresh: 60,
				LoggerTLSPeriod:  33,
			},
			kolide.Host{
				DistributedInterval: 11,
				LoggerTLSPeriod:     33,
				ConfigTLSRefresh:    60,
			},
			map[string]interface{}{
				"distributed_interval": 11,
				"logger_tls_period":    33,
			},
			true,
		},
		// Kolide not managing distributed_interval
		{
			kolide.Host{
				ConfigTLSRefresh:    60,
				DistributedInterval: 11,
			},
			kolide.Host{
				DistributedInterval: 11,
				LoggerTLSPeriod:     33,
				ConfigTLSRefresh:    60,
			},
			map[string]interface{}{
				"logger_tls_period": 33,
			},
			true,
		},
		// SaveHost should not be called with no changes
		{
			kolide.Host{
				DistributedInterval: 11,
				LoggerTLSPeriod:     33,
				ConfigTLSRefresh:    60,
			},
			kolide.Host{
				DistributedInterval: 11,
				LoggerTLSPeriod:     33,
				ConfigTLSRefresh:    60,
			},
			map[string]interface{}{
				"distributed_interval": 11,
				"logger_tls_period":    33,
			},
			false,
		},
	}

	for _, tt := range testCases {
		ds.FIMSectionsFuncInvoked = false

		t.Run("", func(t *testing.T) {
			ctx := hostctx.NewContext(context.Background(), tt.initHost)

			ds.GetOsqueryConfigOptionsFunc = func() (map[string]interface{}, error) {
				return tt.configOptions, nil
			}

			saveHostCalled := false
			ds.SaveHostFunc = func(host *kolide.Host) error {
				saveHostCalled = true
				assert.Equal(t, tt.finalHost, *host)
				return nil
			}

			cfg, err := svc.GetClientConfig(ctx)
			require.Nil(t, err)
			assert.Equal(t, tt.saveHostCalled, saveHostCalled)
			require.True(t, ds.FIMSectionsFuncInvoked)
			require.Condition(t, func() bool {
				_, ok := cfg.Schedule["file_events"]
				return ok
			})
			assert.Equal(t, 400, int(cfg.Schedule["file_events"].Interval))
			assert.Equal(t, "SELECT * FROM file_events;", cfg.Schedule["file_events"].Query)
			require.NotNil(t, cfg.FilePaths)
			require.Condition(t, func() bool {
				_, ok := cfg.FilePaths["etc"]
				return ok
			})
			assert.Len(t, cfg.FilePaths["etc"], 1)
		})
	}
}

func setupOsqueryTests(t *testing.T) (kolide.Datastore, kolide.Service, *clock.MockClock) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)

	_, err = ds.NewAppConfig(&kolide.AppConfig{EnrollSecret: ""})
	require.Nil(t, err)

	mockClock := clock.NewMockClock()
	svc, err := newTestServiceWithClock(ds, nil, mockClock)
	require.Nil(t, err)

	return ds, svc, mockClock
}

type notFoundError struct{}

func (e notFoundError) Error() string {
	return "not found"
}

func (e notFoundError) IsNotFound() bool {
	return true
}

func TestAuthenticationErrors(t *testing.T) {
	ms := new(mock.Store)
	ms.MarkHostSeenFunc = func(*kolide.Host, time.Time) error {
		return nil
	}
	ms.AuthenticateHostFunc = func(nodeKey string) (*kolide.Host, error) {
		return nil, nil
	}

	svc, err := newTestService(ms, nil)
	require.Nil(t, err)
	ctx := context.Background()

	_, err = svc.AuthenticateHost(ctx, "")
	require.NotNil(t, err)
	require.True(t, err.(osqueryError).NodeInvalid())

	_, err = svc.AuthenticateHost(ctx, "foo")
	require.Nil(t, err)

	// return not found error
	ms.AuthenticateHostFunc = func(nodeKey string) (*kolide.Host, error) {
		return nil, notFoundError{}
	}

	_, err = svc.AuthenticateHost(ctx, "foo")
	require.NotNil(t, err)
	require.True(t, err.(osqueryError).NodeInvalid())

	// return other error
	ms.AuthenticateHostFunc = func(nodeKey string) (*kolide.Host, error) {
		return nil, errors.New("foo")
	}

	_, err = svc.AuthenticateHost(ctx, "foo")
	require.NotNil(t, err)
	require.False(t, err.(osqueryError).NodeInvalid())
}
