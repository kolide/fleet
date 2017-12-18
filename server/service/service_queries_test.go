package service

import (
	"context"
	"testing"

	"github.com/kolide/fleet/server/config"
	"github.com/kolide/fleet/server/contexts/viewer"
	"github.com/kolide/fleet/server/datastore/inmem"
	"github.com/kolide/fleet/server/kolide"
	"github.com/kolide/fleet/server/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyQueryYaml(t *testing.T) {
	var gotAuthorID uint
	var gotQueries []*kolide.Query
	ds := &mock.Store{
		QueryStore: mock.QueryStore{
			ApplyQueriesFunc: func(authorID uint, queries []*kolide.Query) error {
				gotAuthorID = authorID
				gotQueries = queries
				return nil
			},
		},
	}
	svc := service{
		ds: ds,
	}

	ctx := context.Background()

	// Error due to missing user
	err := svc.ApplyQueryYaml(ctx, "foo: bar")
	require.NotNil(t, err)

	// Error due to invalid yaml
	ctx = viewer.NewContext(ctx, viewer.Viewer{User: &kolide.User{ID: 1}})
	err = svc.ApplyQueryYaml(ctx, "bad_yaml")
	require.NotNil(t, err)

	// Success
	err = svc.ApplyQueryYaml(ctx, `
---
apiVersion: k8s.kolide.com/v1alpha1
kind: OsqueryQuery
spec:
  name: osquery_version
  description: The version of the Launcher and Osquery process
  query: select launcher.version, osquery.version from kolide_launcher_info launcher, osquery_info osquery;
  support:
    launcher: 0.3.0
    osquery: 2.9.0
---
apiVersion: k8s.kolide.com/v1alpha1
kind: OsqueryQuery
spec:
  name: osquery_schedule
  description: Report performance stats
  query: select name, interval, executions, output_size
---
`,
	)
	require.Nil(t, err)
	assert.Equal(t, uint(1), gotAuthorID)
	assert.Equal(t, []*kolide.Query{
		&kolide.Query{
			Name:        "osquery_version",
			Description: "The version of the Launcher and Osquery process",
			Query:       "select launcher.version, osquery.version from kolide_launcher_info launcher, osquery_info osquery;",
		},
		&kolide.Query{
			Name:        "osquery_schedule",
			Description: "Report performance stats",
			Query:       "select name, interval, executions, output_size",
		},
	}, gotQueries)

}

func TestRoundtripQueriesYaml(t *testing.T) {
	var expectedQueries []*kolide.Query
	var gotQueries []*kolide.Query
	var gotAuthorID uint
	ds := &mock.Store{
		QueryStore: mock.QueryStore{
			ListQueriesFunc: func(opt kolide.ListOptions) ([]*kolide.Query, error) {
				return expectedQueries, nil
			},
			ApplyQueriesFunc: func(authorID uint, queries []*kolide.Query) error {
				gotAuthorID = authorID
				gotQueries = queries
				return nil
			},
		},
	}
	svc := service{
		ds: ds,
	}

	ctx := context.Background()
	ctx = viewer.NewContext(ctx, viewer.Viewer{User: &kolide.User{ID: 1}})

	var testCases = []struct {
		queries []*kolide.Query
	}{
		{[]*kolide.Query{}},
		{[]*kolide.Query{
			&kolide.Query{
				Name:        "frold",
				Description: "bringle",
				Query:       "dahmp",
			},
		}},
		{[]*kolide.Query{
			&kolide.Query{
				Name:        "frold",
				Description: "",
				Query:       "dahmp",
			},
			&kolide.Query{
				Name:        "shmoot",
				Description: "mingus",
				Query:       "kramplit",
			},
		}},
	}

	for _, tt := range testCases {
		t.Run("", func(t *testing.T) {
			expectedQueries = tt.queries
			gotQueries = nil

			yml, err := svc.GetQueryYaml(ctx)
			require.Nil(t, err)

			err = svc.ApplyQueryYaml(ctx, yml)
			require.Nil(t, err)

			assert.Equal(t, expectedQueries, gotQueries)
		})
	}
}

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
