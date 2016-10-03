package datastore

import (
	"testing"

	"github.com/kolide/kolide-ose/server/kolide"
	"github.com/stretchr/testify/assert"
)

func TestDeletePack(t *testing.T) {
	ds := setup(t)
	defer teardown(t, ds)
	testDeletePack(t, ds)
}

func testDeletePack(t *testing.T, ds kolide.Datastore) {
	pack := &kolide.Pack{
		Name: "foo",
	}
	err := ds.NewPack(pack)
	assert.Nil(t, err)
	assert.NotEqual(t, pack.ID, 0)

	pack, err = ds.Pack(pack.ID)
	assert.Nil(t, err)

	err = ds.DeletePack(pack.ID)
	assert.Nil(t, err)

	assert.NotEqual(t, pack.ID, 0)
	pack, err = ds.Pack(pack.ID)
	assert.NotNil(t, err)
}

func TestAddAndRemoveQueryFromPack(t *testing.T) {
	ds := setup(t)
	defer teardown(t, ds)
	testAddAndRemoveQueryFromPack(t, ds)
}

func testAddAndRemoveQueryFromPack(t *testing.T, ds kolide.Datastore) {
	pack := &kolide.Pack{
		Name: "foo",
	}
	err := ds.NewPack(pack)
	assert.Nil(t, err)

	q1 := &kolide.Query{
		Name:  "bar",
		Query: "bar",
	}
	_, err = ds.NewQuery(q1)
	assert.Nil(t, err)
	err = ds.AddQueryToPack(q1.ID, pack.ID)
	assert.Nil(t, err)

	q2 := &kolide.Query{
		Name:  "baz",
		Query: "baz",
	}
	_, err = ds.NewQuery(q2)
	assert.Nil(t, err)
	err = ds.AddQueryToPack(q2.ID, pack.ID)
	assert.Nil(t, err)

	queries, err := ds.GetQueriesInPack(pack)
	assert.Nil(t, err)
	assert.Len(t, queries, 2)

	err = ds.RemoveQueryFromPack(q1, pack)
	assert.Nil(t, err)

	queries, err = ds.GetQueriesInPack(pack)
	assert.Nil(t, err)
	assert.Len(t, queries, 1)
}
