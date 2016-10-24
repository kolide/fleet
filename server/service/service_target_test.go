package service

import (
    "testing"

    "github.com/kolide/kolide-ose/server/datastore"
    "github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
)

func TestSearchTargets(t *testing.T) {
    ds, err := datastore.New("inmem", "")
	assert.Nil(t, err)

    h1, err := ds.NewHost(&kolide.Host{
        HostName: "foo.local",
        PrimaryIP: "192.168.1.10",
    })
    assert.Nil(t, err)
    _ = h1

    h2, err := ds.NewHost(&kolide.Host{
        HostName: "bar.local",
        PrimaryIP: "192.168.1.11",
    })
    assert.Nil(t, err)
    _ = h2

    q1, err := ds.NewQuery(&kolide.Query{
        Name: "query 1",
        Query: "select * from osquery_info;",
    })
    assert.Nil(t, err)
    assert.NotZero(t, q1.ID)

    l1, err := ds.NewLabel(&kolide.Label{
        Name: "label 1",
        QueryID: q1.ID,
    })
    _ = l1
}