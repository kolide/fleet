package service

import (
    "testing"

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
        HostName: "foo.local",
        PrimaryIP: "192.168.1.10",
    })
    require.Nil(t, err)

    _, err = ds.NewHost(&kolide.Host{
        HostName: "bar.local",
        PrimaryIP: "192.168.1.11",
    })
    require.Nil(t, err)

    q1, err := ds.NewQuery(&kolide.Query{
        Name: "query foo",
        Query: "select * from osquery_info;",
    })
    require.Nil(t, err)
    require.NotZero(t, q1.ID)

    l1, err := ds.NewLabel(&kolide.Label{
        Name: "label foo",
        QueryID: q1.ID,
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

}

func TestSearchWithOmit(t *testing.T) {

}