package service

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/datastore/inmem"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func TestGetScheduledQueriesInPack(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)
	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)
	ctx := context.Background()

	u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID, true)
	q2 := test.NewQuery(t, ds, "bar", "select * from time;", u1.ID, true)
	p1 := test.NewPack(t, ds, "baz")
	sq1 := test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, false)

	queries, err := svc.GetScheduledQueriesInPack(ctx, p1.ID, kolide.ListOptions{})
	require.Nil(t, err)
	require.Len(t, queries, 1)
	assert.Equal(t, sq1.ID, queries[0].ID)

	test.NewScheduledQuery(t, ds, p1.ID, q2.ID, 60, false, false)
	test.NewScheduledQuery(t, ds, p1.ID, q2.ID, 60, true, false)

	queries, err = svc.GetScheduledQueriesInPack(ctx, p1.ID, kolide.ListOptions{})
	require.Nil(t, err)
	require.Len(t, queries, 3)
}

func TestGetScheduledQuery(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)
	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)
	ctx := context.Background()

	u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID, true)
	p1 := test.NewPack(t, ds, "baz")
	sq1 := test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, false)

	query, err := svc.GetScheduledQuery(ctx, sq1.ID)
	require.Nil(t, err)
	assert.Equal(t, uint(60), query.Interval)
}

func TestModifyScheduledQuery(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)
	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)
	ctx := context.Background()

	u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID, true)
	p1 := test.NewPack(t, ds, "baz")
	sq1 := test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, false)

	query, err := svc.GetScheduledQuery(ctx, sq1.ID)
	require.Nil(t, err)
	assert.Equal(t, uint(60), query.Interval)

	interval := uint(120)
	queryPayload := kolide.ScheduledQueryPayload{
		Interval: &interval,
	}
	query, err = svc.ModifyScheduledQuery(ctx, sq1.ID, queryPayload)
	assert.Equal(t, uint(120), query.Interval)

	queryVerify, err := svc.GetScheduledQuery(ctx, sq1.ID)
	require.Nil(t, err)
	assert.Equal(t, uint(120), queryVerify.Interval)
}

func TestDeleteScheduledQuery(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)
	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)
	ctx := context.Background()

	u1 := test.NewUser(t, ds, "Admin", "admin", "admin@kolide.co", true)
	q1 := test.NewQuery(t, ds, "foo", "select * from time;", u1.ID, true)
	p1 := test.NewPack(t, ds, "baz")
	sq1 := test.NewScheduledQuery(t, ds, p1.ID, q1.ID, 60, false, false)

	query, err := svc.GetScheduledQuery(ctx, sq1.ID)
	require.Nil(t, err)
	assert.Equal(t, uint(60), query.Interval)

	err = svc.DeleteScheduledQuery(ctx, sq1.ID)
	require.Nil(t, err)

	_, err = svc.GetScheduledQuery(ctx, sq1.ID)
	require.NotNil(t, err)
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func TestDecodeScheduleQueryRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/schedule", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeScheduleQueryRequest(context.Background(), request)
		require.Nil(t, err)

		params := r.(scheduleQueryRequest)
		assert.Equal(t, uint(5), params.PackID)
		assert.Equal(t, uint(1), params.QueryID)
		assert.Equal(t, uint(60), params.Interval)
		assert.Equal(t, true, *params.Snapshot)
	}).Methods("POST")

	var body bytes.Buffer
	body.Write([]byte(`{
		"pack_id": 5,
		"query_id": 1,
		"interval": 60,
		"snapshot": true
	}`))

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("POST", "/api/v1/kolide/schedule", &body),
	)
}

func TestDecodeModifyScheduledQueryRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/scheduled/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeModifyScheduledQueryRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(modifyScheduledQueryRequest)
		assert.Equal(t, uint(1), params.ID)
		assert.Equal(t, uint(5), *params.payload.PackID)
		assert.Equal(t, uint(6), *params.payload.QueryID)
		assert.Equal(t, true, *params.payload.Removed)
		assert.Equal(t, uint(60), *params.payload.Interval)
		assert.Equal(t, uint(1), *params.payload.Shard)
	}).Methods("PATCH")

	var body bytes.Buffer
	body.Write([]byte(`{
        "pack_id": 5,
		"query_id": 6,
		"removed": true,
		"interval": 60,
		"shard": 1
    }`))

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("PATCH", "/api/v1/kolide/scheduled/1", &body),
	)
}

func TestDecodeDeleteScheduledQueryRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/scheduled/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeDeleteScheduledQueryRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(deleteScheduledQueryRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("DELETE")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("DELETE", "/api/v1/kolide/scheduled/1", nil),
	)
}

func TestDecodeGetScheduledQueryRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/scheduled/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeGetScheduledQueryRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(getScheduledQueryRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("GET")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/api/v1/kolide/scheduled/1", nil),
	)
}

func TestDecodeGetScheduledQueriesInPackRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/packs/{id}/scheduled", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeGetScheduledQueriesInPackRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(getScheduledQueriesInPackRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("GET")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/api/v1/kolide/packs/1/scheduled", nil),
	)
}
