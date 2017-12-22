package service

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/datastore/inmem"
	"github.com/kolide/fleet/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

////////////////////////////////////////////////////////////////////////////////
// Service
////////////////////////////////////////////////////////////////////////////////

func TestListQueries(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	queries, err := svc.ListQueries(ctx, kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, queries, 0)

	name := "foo"
	query := "select * from time"
	_, err = svc.NewQuery(ctx, kolide.QueryPayload{
		Name:  &name,
		Query: &query,
	})
	assert.Nil(t, err)

	queries, err = svc.ListQueries(ctx, kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, queries, 1)
}

func TestGetQuery(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	query := &kolide.Query{
		Name:  "foo",
		Query: "select * from time;",
	}
	query, err = ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotZero(t, query.ID)

	queryVerify, err := svc.GetQuery(ctx, query.ID)
	assert.Nil(t, err)
	assert.Equal(t, query.ID, queryVerify.ID)
}

func TestNewQuery(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	createTestUsers(t, ds)
	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	user, err := ds.User("admin1")
	require.Nil(t, err)

	ctx := context.Background()
	ctx = viewer.NewContext(ctx, viewer.Viewer{User: user})

	name := "foo"
	query := "select * from time;"
	q, err := svc.NewQuery(ctx, kolide.QueryPayload{
		Name:  &name,
		Query: &query,
	})
	assert.Nil(t, err)
	assert.Equal(t, "Test Name admin1", q.AuthorName)
	assert.Equal(t, []kolide.Pack{}, q.Packs)

	queries, err := ds.ListQueries(kolide.ListOptions{})
	assert.Nil(t, err)
	if assert.Len(t, queries, 1) {
		assert.Equal(t, "Test Name admin1", queries[0].AuthorName)
	}
}

func TestModifyQuery(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	query := &kolide.Query{
		Name:  "foo",
		Query: "select * from time;",
	}
	query, err = ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotZero(t, query.ID)

	newName := "bar"
	queryVerify, err := svc.ModifyQuery(ctx, query.ID, kolide.QueryPayload{
		Name: &newName,
	})
	assert.Nil(t, err)

	assert.Equal(t, query.ID, queryVerify.ID)
	assert.Equal(t, "bar", queryVerify.Name)
}

func TestDeleteQuery(t *testing.T) {
	ds, err := inmem.New(config.TestConfig())
	assert.Nil(t, err)

	svc, err := newTestService(ds, nil)
	assert.Nil(t, err)

	ctx := context.Background()

	query := &kolide.Query{
		Name:  "foo",
		Query: "select * from time;",
	}
	query, err = ds.NewQuery(query)
	assert.Nil(t, err)
	assert.NotZero(t, query.ID)

	err = svc.DeleteQuery(ctx, query.ID)
	assert.Nil(t, err)

	queries, err := ds.ListQueries(kolide.ListOptions{})
	assert.Nil(t, err)
	assert.Len(t, queries, 0)
}

////////////////////////////////////////////////////////////////////////////////
// Transport
////////////////////////////////////////////////////////////////////////////////

func TestDecodeCreateQueryRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/queries", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeCreateQueryRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(createQueryRequest)
		assert.Equal(t, "foo", *params.payload.Name)
		assert.Equal(t, "select * from time;", *params.payload.Query)
	}).Methods("POST")

	var body bytes.Buffer
	body.Write([]byte(`{
        "name": "foo",
        "query": "select * from time;"
    }`))

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("POST", "/api/v1/kolide/queries", &body),
	)
}

func TestDecodeModifyQueryRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/queries/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeModifyQueryRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(modifyQueryRequest)
		assert.Equal(t, "foo", *params.payload.Name)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("PATCH")

	var body bytes.Buffer
	body.Write([]byte(`{
        "name": "foo"
    }`))

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("PATCH", "/api/v1/kolide/queries/1", &body),
	)
}

func TestDecodeDeleteQueryRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/queries/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeDeleteQueryRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(deleteQueryRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("DELETE")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("DELETE", "/api/v1/kolide/queries/1", nil),
	)
}

func TestDecodeGetQueryRequest(t *testing.T) {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kolide/queries/{id}", func(writer http.ResponseWriter, request *http.Request) {
		r, err := decodeGetQueryRequest(context.Background(), request)
		assert.Nil(t, err)

		params := r.(getQueryRequest)
		assert.Equal(t, uint(1), params.ID)
	}).Methods("GET")

	router.ServeHTTP(
		httptest.NewRecorder(),
		httptest.NewRequest("GET", "/api/v1/kolide/queries/1", nil),
	)
}
