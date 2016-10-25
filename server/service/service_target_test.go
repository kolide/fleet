package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/kolide/kolide-ose/server/datastore"
	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestSearchTargets(t *testing.T) {
	ds, err := datastore.New("inmem", "")
	require.Nil(t, err)

	svc, err := newTestService(ds)
	require.Nil(t, err)

	ctx := context.Background()

	h1, err := ds.NewHost(&kolide.Host{
		HostName:  "foo.local",
		PrimaryIP: "192.168.1.10",
	})
	require.Nil(t, err)

	l1, err := ds.NewLabel(&kolide.Label{
		Name:    "label foo",
		QueryID: 1,
	})

	results, count, err := svc.SearchTargets(ctx, "foo", nil)
	require.Nil(t, err)

	require.Len(t, results.Hosts, 1)
	assert.Equal(t, h1.HostName, results.Hosts[0].HostName)

	require.Len(t, results.Labels, 1)
	assert.Equal(t, l1.Name, results.Labels[0].Name)

	assert.Equal(t, uint(1), count)
}

func TestCountHostsInTargets(t *testing.T) {
	ds, err := datastore.New("inmem", "")
	require.Nil(t, err)

	svc, err := newTestService(ds)
	require.Nil(t, err)

	ctx := context.Background()

	h1, err := ds.NewHost(&kolide.Host{
		HostName:  "foo.local",
		PrimaryIP: "192.168.1.10",
	})
	require.Nil(t, err)

	h2, err := ds.NewHost(&kolide.Host{
		HostName:  "bar.local",
		PrimaryIP: "192.168.1.11",
	})
	require.Nil(t, err)

	h3, err := ds.NewHost(&kolide.Host{
		HostName:  "baz.local",
		PrimaryIP: "192.168.1.12",
	})
	require.Nil(t, err)

	l1, err := ds.NewLabel(&kolide.Label{
		Name:    "label foo",
		QueryID: 1,
	})
	require.Nil(t, err)
	require.NotZero(t, l1.ID)
	l1ID := fmt.Sprintf("%d", l1.ID)

	for _, h := range []*kolide.Host{h1, h2, h3} {
		err = ds.RecordLabelQueryExecutions(h, map[string]bool{l1ID: true}, time.Now())
		assert.Nil(t, err)
	}

	count, err := svc.CountHostsInTargets(ctx, nil, []kolide.Label{*l1})
	assert.Nil(t, err)
	assert.Equal(t, uint(3), count)
}

func TestSearchWithOmit(t *testing.T) {

}
