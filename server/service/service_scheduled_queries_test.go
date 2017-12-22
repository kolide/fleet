package service

import "testing"

func TestGetScheduledQueriesInPack(t *testing.T) {
	// ds, err := inmem.New(config.TestConfig())
	// assert.Nil(t, err)
	// svc, err := newTestService(ds, nil)
	// assert.Nil(t, err)
	// ctx := context.Background()

	// u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	// q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID, true)
	// q2 := test.NewQuery(t, ds, "bar", "select * from time;", u1.ID, true)
	// p1 := test.NewPack(t, ds, "baz")
	// sq1 := test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, false)

	// queries, err := svc.GetScheduledQueriesInPack(ctx, p1.ID, kolide.ListOptions{})
	// require.Nil(t, err)
	// require.Len(t, queries, 1)
	// assert.Equal(t, sq1.ID, queries[0].ID)

	// test.NewScheduledQuery(t, ds, p1.ID, q2.ID, 60, false, false)
	// test.NewScheduledQuery(t, ds, p1.ID, q2.ID, 60, true, false)

	// queries, err = svc.GetScheduledQueriesInPack(ctx, p1.ID, kolide.ListOptions{})
	// require.Nil(t, err)
	// require.Len(t, queries, 3)
}
