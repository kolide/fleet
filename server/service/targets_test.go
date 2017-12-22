package service

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/datastore/inmem"
	"github.com/kolide/fleet/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func TestSearchTargets(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)

	svc, err := newTestService(ds, nil)
	require.Nil(t, err)

	ctx := context.Background()

	h1, err := ds.NewHost(&kolide.Host{
		HostName: "foo.local",
	})
	require.Nil(t, err)

	l1, err := ds.NewLabel(&kolide.Label{
		Name:  "label foo",
		Query: "query foo",
	})
	require.Nil(t, err)

	results, err := svc.SearchTargets(ctx, "foo", nil, nil)
	require.Nil(t, err)

	require.Len(t, results.Hosts, 1)
	assert.Equal(t, h1.HostName, results.Hosts[0].HostName)

	require.Len(t, results.Labels, 1)
	assert.Equal(t, l1.Name, results.Labels[0].Name)
}

func TestSearchWithOmit(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)

	svc, err := newTestService(ds, nil)
	require.Nil(t, err)

	ctx := context.Background()

	h1, err := ds.NewHost(&kolide.Host{
		HostName: "foo.local",
		NodeKey:  "1",
		UUID:     "1",
	})
	require.Nil(t, err)

	h2, err := ds.NewHost(&kolide.Host{
		HostName: "foobar.local",
		NodeKey:  "2",
		UUID:     "2",
	})
	require.Nil(t, err)

	l1, err := ds.NewLabel(&kolide.Label{
		Name:  "label foo",
		Query: "query foo",
	})

	{
		results, err := svc.SearchTargets(ctx, "foo", nil, nil)
		require.Nil(t, err)

		require.Len(t, results.Hosts, 2)

		require.Len(t, results.Labels, 1)
		assert.Equal(t, l1.Name, results.Labels[0].Name)
	}

	{
		results, err := svc.SearchTargets(ctx, "foo", []uint{h2.ID}, nil)
		require.Nil(t, err)

		require.Len(t, results.Hosts, 1)
		assert.Equal(t, h1.HostName, results.Hosts[0].HostName)

		require.Len(t, results.Labels, 1)
		assert.Equal(t, l1.Name, results.Labels[0].Name)
	}
}

func TestSearchHostsInLabels(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)

	svc, err := newTestService(ds, nil)
	require.Nil(t, err)

	ctx := context.Background()

	h1, err := ds.NewHost(&kolide.Host{
		HostName: "foo.local",
		NodeKey:  "1",
		UUID:     "1",
	})
	require.Nil(t, err)

	h2, err := ds.NewHost(&kolide.Host{
		HostName: "bar.local",
		NodeKey:  "2",
		UUID:     "2",
	})
	require.Nil(t, err)

	h3, err := ds.NewHost(&kolide.Host{
		HostName: "baz.local",
		NodeKey:  "3",
		UUID:     "3",
	})
	require.Nil(t, err)

	l1, err := ds.NewLabel(&kolide.Label{
		Name:  "label foo",
		Query: "query foo",
	})
	require.Nil(t, err)
	require.NotZero(t, l1.ID)

	for _, h := range []*kolide.Host{h1, h2, h3} {
		err = ds.RecordLabelQueryExecutions(h, map[uint]bool{l1.ID: true}, time.Now())
		assert.Nil(t, err)
	}

	results, err := svc.SearchTargets(ctx, "baz", nil, nil)
	require.Nil(t, err)

	require.Len(t, results.Hosts, 1)
	assert.Equal(t, h3.HostName, results.Hosts[0].HostName)
}

func TestSearchResultsLimit(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	require.Nil(t, err)

	svc, err := newTestService(ds, nil)
	require.Nil(t, err)

	ctx := context.Background()

	for i := 0; i < 15; i++ {
		_, err := ds.NewHost(&kolide.Host{
			HostName: fmt.Sprintf("foo.%d.local", i),
			NodeKey:  fmt.Sprintf("%d", i+1),
			UUID:     fmt.Sprintf("%d", i+1),
		})
		require.Nil(t, err)
	}
	targets, err := svc.SearchTargets(ctx, "foo", nil, nil)
	require.Nil(t, err)
	assert.Len(t, targets.Hosts, 10)
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func TestDecodeSearchTargetsRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/targets", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeSearchTargetsRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(searchTargetsRequest)
		assert.Equal(t, "bar", params.Query)
		assert.Len(t, params.Selected.Hosts, 3)
		assert.Len(t, params.Selected.Labels, 2)
	}).Methods("POST")
	var body bytes.Buffer

	body.Write([]byte(`{
        "query": "bar",
		"selected": {
			"hosts": [
				1,
				2,
				3
			],
			"labels": [
				1,
				2
			]
		}
    }`))

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("POST", "/api/v1/kolide/targets", &body),
	)
}
