package datastore

import (
	"testing"

	"github.com/kolide/fleet/server/kolide"
)

func testListScheduledQueriesInPack(t *testing.T, ds kolide.Datastore) {
	// 	u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	// 	q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID, true)
	// 	q2 := test.NewQuery(t, ds, "bar", "select * from time;", u1.ID, true)
	// 	p1 := test.NewPack(t, ds, "baz")

	// 	test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, false)

	// 	queries, err := ds.ListScheduledQueriesInPack(p1.ID, kolide.ListOptions{})
	// 	require.Nil(t, err)
	// 	require.Len(t, queries, 1)
	// 	assert.Equal(t, uint(60), queries[0].Interval)

	// 	test.NewScheduledQuery(t, ds, p1.ID, q2.ID, 60, false, false)
	// 	test.NewScheduledQuery(t, ds, p1.ID, q2.ID, 60, true, false)

	// 	queries, err = ds.ListScheduledQueriesInPack(p1.ID, kolide.ListOptions{})
	// 	require.Nil(t, err)
	// 	require.Len(t, queries, 3)
}
